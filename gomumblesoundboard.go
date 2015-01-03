package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/go-martini/martini"
)

func main() {
	files := make(map[string]string)
	var stream *gumble_ffmpeg.Stream
	targetChannel := flag.String("channel", "Root", "channel the bot will join")

	gumbleutil.Main(func(_ *gumble.Config, client *gumble.Client) {
		var err error
		stream, err = gumble_ffmpeg.New(client)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		for _, file := range flag.Args() {
			key := filepath.Base(file)
			files[key] = file
		}
	}, gumbleutil.Listener{
		// Connect event
		Connect: func(e *gumble.ConnectEvent) {
			fmt.Printf("GoMumbleSoundboard loaded (%d files)\n", len(files))
			fmt.Printf("Connected to %s\n", e.Client.Conn().RemoteAddr())
			if e.WelcomeMessage != "" {
				fmt.Printf("Welcome message: %s\n", e.WelcomeMessage)
			}
			fmt.Printf("Channel: %s\n", e.Client.Self().Channel().Name())

			if e.Client.Self().Channel().Name() != *targetChannel {
				target := e.Client.Self().Channel().Find(*targetChannel)
				e.Client.Self().Move(target)
				fmt.Printf("Moved to: %s\n", target.Name())
			}

			// Start webserver
			m := martini.Classic()
			// martini.Static() is used, so public/index.html gets automagically served
			m.Get("/files.json", func() string {
				keys := make([]string, 0, len(files))
			    for k := range files {
			        keys = append(keys, k)
			    }
				js, _ := json.Marshal(keys)
				return string(js)
			})
			m.Get("/play/:file", func(params martini.Params) (int, string) {
				file, ok := files[params["file"]]
				if !ok {
					return 404, "not found"
				}
				stream.Stop()
				if err := stream.Play(file); err != nil {
					return 400, fmt.Sprintf("%s\n", err)
				} else {
					return 200, fmt.Sprintf("Playing %s\n", file)
				}
			})
			m.Get("/stop", func() string {
				stream.Stop()
				return "ok"
			})
			m.Run()
		},
	})
}