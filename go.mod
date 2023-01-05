module go.arsenm.dev/lemmy-reply-bot

go 1.19

//replace go.arsenm.dev/go-lemmy => /home/arsen/Code/go-lemmy

require (
	github.com/pelletier/go-toml/v2 v2.0.6
	github.com/spf13/pflag v1.0.5
	github.com/vmihailenco/msgpack/v5 v5.3.5
	go.arsenm.dev/go-lemmy v0.0.0-20230105214607-754cfa602c10
	go.arsenm.dev/logger v0.0.0-20221007032343-cbffce4f4334
	go.arsenm.dev/pcre v0.0.0-20220530205550-74594f6c8b0e
)

require (
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gookit/color v1.5.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	modernc.org/libc v1.16.8 // indirect
	modernc.org/mathutil v1.4.1 // indirect
	modernc.org/memory v1.1.1 // indirect
)
