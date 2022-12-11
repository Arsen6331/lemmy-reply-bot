package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/pflag"
	"github.com/vmihailenco/msgpack/v5"
	"go.arsenm.dev/go-lemmy"
	"go.arsenm.dev/go-lemmy/types"
	"go.arsenm.dev/logger/log"
)

func main() {
	configPath := pflag.StringP("config", "c", "./lemmy-reply-bot.toml", "Path to the config file")
	dryRun := pflag.BoolP("dry-run", "D", false, "Don't actually send comments, just check for matches")
	pflag.Parse()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := loadConfig(*configPath)
	if err != nil {
		log.Fatal("Error loading config file").Err(err).Send()
	}

	rhc := retryablehttp.NewClient()
	rhc.Logger = retryableLogger{}

	c, err := lemmy.NewWithClient(cfg.Lemmy.InstanceURL, rhc.StandardClient())
	if err != nil {
		log.Fatal("Error creating new Lemmy API client").Err(err).Send()
	}

	err = c.Login(ctx, types.Login{
		UsernameOrEmail: cfg.Lemmy.Account.UserOrEmail,
		Password:        cfg.Lemmy.Account.Password,
	})
	if err != nil {
		log.Fatal("Error logging in to Lemmy instance").Err(err).Send()
	}

	log.Info("Successfully logged in to Lemmy instance").Send()

	replyCh := make(chan replyJob, 200)

	if !*dryRun {
		go commentReplyWorker(ctx, c, replyCh)
	}

	commentWorker(ctx, c, replyCh)
}

func commentWorker(ctx context.Context, c *lemmy.Client, replyCh chan<- replyJob) {
	repliedIDs := map[int]struct{}{}

	repliedStore, err := os.Open("replied.bin")
	if err == nil {
		err = msgpack.NewDecoder(repliedStore).Decode(&repliedIDs)
		if err != nil {
			log.Warn("Error decoding reply store").Err(err).Send()
		}
		repliedStore.Close()
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cr, err := c.Comments(ctx, types.GetComments{
				Sort:  types.NewOptional(types.CommentSortNew),
				Limit: types.NewOptional(200),
			})
			if err != nil {
				log.Warn("Error while trying to get comments").Err(err).Send()
				continue
			}

			for _, commentView := range cr.Comments {
				if _, ok := repliedIDs[commentView.Comment.ID]; ok {
					continue
				}

				for i, reply := range cfg.Replies {
					re := compiledRegexes[reply.Regex]
					if !re.MatchString(commentView.Comment.Content) {
						continue
					}

					log.Info("Matched comment body").
						Int("reply-index", i).
						Int("comment-id", commentView.Comment.ID).
						Send()

					job := replyJob{
						CommentID: commentView.Comment.ID,
						PostID:    commentView.Comment.PostID,
					}

					matches := re.FindStringSubmatch(commentView.Comment.Content)
					job.Content = expandStr(reply.Msg, func(s string) string {
						i, err := strconv.Atoi(s)
						if err != nil {
							log.Debug("Message variable is not an integer, returning empty string").Str("var", s).Send()
							return ""
						}

						if i+1 > len(matches) {
							log.Debug("Message variable exceeds match length").Int("length", len(matches)).Int("var", i).Send()
							return ""
						}

						log.Debug("Message variable found, returning").Int("var", i).Str("found", matches[i]).Send()
						return matches[i]
					})

					replyCh <- job

					repliedIDs[commentView.Comment.ID] = struct{}{}
				}
			}
		case <-ctx.Done():
			repliedStore, err := os.Create("replied.bin")
			if err != nil {
				log.Warn("Error creating reply store file").Err(err).Send()
				return
			}

			err = msgpack.NewEncoder(repliedStore).Encode(repliedIDs)
			if err != nil {
				log.Warn("Error encoding replies to reply store").Err(err).Send()
			}

			repliedStore.Close()
			return
		}
	}
}

type replyJob struct {
	Content   string
	CommentID int
	PostID    int
}

func commentReplyWorker(ctx context.Context, c *lemmy.Client, ch <-chan replyJob) {
	for {
		select {
		case reply := <-ch:
			cr, err := c.CreateComment(ctx, types.CreateComment{
				PostID:   reply.PostID,
				ParentID: types.NewOptional(reply.CommentID),
				Content:  reply.Content,
			})
			if err != nil {
				log.Warn("Error while trying to create new comment").Err(err).Send()
			}

			log.Info("Created new comment").
				Int("post-id", reply.PostID).
				Int("parent-id", reply.CommentID).
				Int("comment-id", cr.CommentView.Comment.ID).
				Send()

			// Make sure requests don't happen too quickly
			time.Sleep(1 * time.Second)
		case <-ctx.Done():
			return
		}
	}
}

func expandStr(s string, mapping func(string) string) string {
	strings.ReplaceAll(s, "$$", "${_escaped_dollar_symbol}")
	return os.Expand(s, func(s string) string {
		if s == "_escaped_dollar_symbol" {
			return "$"
		}
		return mapping(s)
	})
}
