package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"

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

	c, err := lemmy.NewWebSocket(cfg.Lemmy.InstanceURL)
	if err != nil {
		log.Fatal("Error creating new Lemmy API client").Err(err).Send()
	}

	err = c.ClientLogin(ctx, types.Login{
		UsernameOrEmail: cfg.Lemmy.Account.UserOrEmail,
		Password:        cfg.Lemmy.Account.Password,
	})
	if err != nil {
		log.Fatal("Error logging in to Lemmy instance").Err(err).Send()
	}

	log.Info("Successfully logged in to Lemmy instance").Send()

	joinAll(c)

	c.OnReconnect(func(c *lemmy.WSClient) {
		joinAll(c)
		log.Info("Successfully reconnected to WebSocket").Send()
	})

	replyCh := make(chan replyJob, 200)

	if !*dryRun {
		go commentReplyWorker(ctx, c, replyCh)
	}

	commentWorker(ctx, c, replyCh)
}

func commentWorker(ctx context.Context, c *lemmy.WSClient, replyCh chan<- replyJob) {
	repliedIDs := map[int]struct{}{}

	repliedStore, err := os.Open("replied.bin")
	if err == nil {
		err = msgpack.NewDecoder(repliedStore).Decode(&repliedIDs)
		if err != nil {
			log.Warn("Error decoding reply store").Err(err).Send()
		}
		repliedStore.Close()
	}

	for {
		select {
		case res := <-c.Responses():
			if res.IsOneOf(types.UserOperationCRUDCreateComment, types.UserOperationCRUDEditComment) {
				var cr types.CommentResponse
				err = lemmy.DecodeResponse(res.Data, &cr)
				if err != nil {
					log.Warn("Error while trying to decode comment").Err(err).Send()
					continue
				}

				if _, ok := repliedIDs[cr.CommentView.Comment.ID]; ok {
					continue
				}

				for i, reply := range cfg.Replies {
					re := compiledRegexes[reply.Regex]
					if !re.MatchString(cr.CommentView.Comment.Content) {
						continue
					}

					log.Info("Matched comment body").
						Int("reply-index", i).
						Int("comment-id", cr.CommentView.Comment.ID).
						Send()

					job := replyJob{
						CommentID: cr.CommentView.Comment.ID,
						PostID:    cr.CommentView.Comment.PostID,
					}

					matches := re.FindAllStringSubmatch(cr.CommentView.Comment.Content, -1)
					job.Content, err = executeTmpl(compiledTmpls[reply.Regex], matches)
					if err != nil {
						log.Warn("Error while executing template").Err(err).Send()
						continue
					}

					replyCh <- job

					repliedIDs[cr.CommentView.Comment.ID] = struct{}{}
				}
			}
		case err := <-c.Errors():
			log.Warn("Lemmy client error").Err(err).Send()
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

func commentReplyWorker(ctx context.Context, c *lemmy.WSClient, ch <-chan replyJob) {
	for {
		select {
		case reply := <-ch:
			err := c.Request(types.UserOperationCRUDCreateComment, types.CreateComment{
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
				Send()
		case <-ctx.Done():
			return
		}
	}
}

func executeTmpl(tmpl *template.Template, matches [][]string) (string, error) {
	sb := &strings.Builder{}
	err := tmpl.Execute(sb, matches)
	return sb.String(), err
}

func joinAll(c *lemmy.WSClient) {
	err := c.Request(types.UserOperationUserJoin, nil)
	if err != nil {
		log.Fatal("Error joining WebSocket user context").Err(err).Send()
	}

	err = c.Request(types.UserOperationCommunityJoin, types.CommunityJoin{
		CommunityID: 0,
	})
	if err != nil {
		log.Fatal("Error joining WebSocket community context").Err(err).Send()
	}
}
