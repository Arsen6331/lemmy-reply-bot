package main

import (
	"net/url"
	"os"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pelletier/go-toml/v2"
	"go.arsenm.dev/logger/log"
	"go.arsenm.dev/pcre"
)

type Config struct {
	Lemmy struct {
		InstanceURL string `toml:"instanceURL"`
		Account     struct {
			UserOrEmail string `toml:"userOrEmail"`
			Password    string `toml:"password"`
		} `toml:"account"`
	} `toml:"lemmy"`
	Replies []Reply `toml:"reply"`
}

type Reply struct {
	Regex string `toml:"regex"`
	Msg   string `toml:"msg"`
}

var (
	cfg             = Config{}
	compiledRegexes = map[string]*pcre.Regexp{}
	compiledTmpls   = map[string]*template.Template{}
)

func loadConfig(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fi.Mode().Perm() != 0o600 {
		log.Fatal("Your config file's permissions are insecure. Please use chmod to set them to 600. Refusing to start.").Send()
	}

	fl, err := os.Open(path)
	if err != nil {
		return err
	}

	err = toml.NewDecoder(fl).Decode(&cfg)
	if err != nil {
		return err
	}

	err = compileReplies(cfg.Replies)
	if err != nil {
		return err
	}

	validateConfig()
	return nil
}

func compileReplies(replies []Reply) error {
	for i, reply := range replies {
		if _, ok := compiledRegexes[reply.Regex]; ok {
			continue
		}

		re, err := pcre.Compile(reply.Regex)
		if err != nil {
			return err
		}
		compiledRegexes[reply.Regex] = re

		tmpl, err := template.
			New(strconv.Itoa(i)).
			Funcs(tmplFuncs).
			Funcs(sprig.TxtFuncMap()).
			Parse(reply.Msg)
		if err != nil {
			return err
		}
		compiledTmpls[reply.Regex] = tmpl
	}
	return nil
}

func validateConfig() {
	_, err := url.Parse(cfg.Lemmy.InstanceURL)
	if err != nil {
		log.Fatal("Lemmy instance URL is not valid").Err(err).Send()
	}

	for i, reply := range cfg.Replies {
		re := compiledRegexes[reply.Regex]

		if re.MatchString(reply.Msg) {
			log.Fatal("Regular expression matches message. This may create an infinite loop. Refusing to start.").Int("reply-index", i).Send()
		}
	}
}
