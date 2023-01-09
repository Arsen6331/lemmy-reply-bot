module go.arsenm.dev/lemmy-reply-bot

go 1.19

replace go.arsenm.dev/go-lemmy => /home/arsen/Code/go-lemmy

require (
	github.com/pelletier/go-toml/v2 v2.0.6
	github.com/spf13/pflag v1.0.5
	github.com/vmihailenco/msgpack/v5 v5.3.5
	go.arsenm.dev/go-lemmy v0.16.8-0.20230109205406-c0aced05f0cd
	go.arsenm.dev/logger v0.0.0-20230104225304-d706171ea6df
	go.arsenm.dev/pcre v0.0.0-20220530205550-74594f6c8b0e
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gookit/color v1.5.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/sys v0.1.0 // indirect
	modernc.org/libc v1.16.8 // indirect
	modernc.org/mathutil v1.4.1 // indirect
	modernc.org/memory v1.1.1 // indirect
)
