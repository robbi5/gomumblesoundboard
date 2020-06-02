package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"layeh.com/gumble/gumbleutil"
)

func main() {
	targetChannel := flag.String("channel", "Root", "channel the bot will join")
	files := make(map[string]string)
	var volume float32 = 0.5

	gumbleutil.Main(
		gumbleutil.AutoBitrate,
		gumbleutil.Listener{
			Connect: func(e *gumble.ConnectEvent) {
				var stream *gumbleffmpeg.Stream

				for _, file := range flag.Args() {
					key := filepath.Base(file)
					files[key] = file
				}

				fmt.Printf("GoMumbleSoundboard loaded (%d files)\n", len(files))
				fmt.Printf("Connected to %s\n", e.Client.Conn.RemoteAddr())
				fmt.Printf("Current Channel: %s\n", e.Client.Self.Channel.Name)

				if *targetChannel != "" && e.Client.Self.Channel.Name != *targetChannel {
					channelPath := strings.Split(*targetChannel, "/")
					target := e.Client.Self.Channel.Find(channelPath...)
					if target == nil {
						fmt.Printf("Cannot find channel named %s\n", *targetChannel)
						os.Exit(1)
					}
					e.Client.Self.Move(target)
					fmt.Printf("Moved to: %s\n", target.Name)
				}

				// Start webserver
				m := martini.Classic()
				// martini.Static() is used, so public/index.html gets automagically served
				m.Get("/files.json", func() string {
					keys := make([]string, 0, len(files))
					for k := range files {
						keys = append(keys, k)
					}
					// Sort keys into alphabetical order. Sick of things moving around
					ss := sort.StringSlice(keys)
					ss.Sort()

					js, _ := json.Marshal(ss)
					return string(js)
				})
				m.Get("/play/:file", func(params martini.Params) (int, string) {
					file, ok := files[params["file"]]
					if !ok {
						return 404, "not found"
					}
					stream := gumbleffmpeg.New(e.Client, gumbleffmpeg.SourceFile(file))
					stream.Volume = volume
					err := stream.Play()
					if err != nil {
						return 400, fmt.Sprintf("%s\n", err)
					}
					return 200, fmt.Sprintf("Playing %s\n", file)
				})
				m.Get("/volume/:volume", func(params martini.Params) (int, string) {
					strVol := params["volume"]
					vol, err := strconv.Atoi(strVol)
					if err != nil {
						return 400, "NaN"
					}

					if vol < 0 && vol > 100 {
						return 400, "Number too small or too large"
					}

					volume = float32(vol) / 100
					return 200, fmt.Sprintf("Volume set to %d\n", vol)
				})
				m.Get("/stop", func() string {
					stream.Stop()
					return "ok"
				})
				m.Run()
			},
			Disconnect: func(e *gumble.DisconnectEvent) {
				os.Exit(1)
			},
		})
}
