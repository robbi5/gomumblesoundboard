package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"layeh.com/gumble/gumbleutil"
	_ "layeh.com/gumble/opus"
)

//go:generate go-assets-builder public -s "/public" -o assets.go

type staticAssetsFS struct {
	fs *assets.FileSystem
}

func (f staticAssetsFS) Exists(prefix string, path string) bool {
	if prefix != "/" {
		panic("We don't support prefixes except for the empty one")
	}

	_, ok := f.fs.Files[path]
	return ok
}

func (f staticAssetsFS) Open(name string) (http.File, error) {
	file, err := f.fs.Open(name)
	return file, err
}

func main() {
	targetChannel := flag.String("channel", "Root", "channel the bot will join")
	soundfiles := make(map[string]string)
	var volume float32 = 0.5

	gumbleutil.Main(
		gumbleutil.AutoBitrate,
		gumbleutil.Listener{
			Connect: func(e *gumble.ConnectEvent) {
				var stream *gumbleffmpeg.Stream

				for _, dir := range flag.Args() {
					fmt.Printf("Dir: %s\n", dir)
					files, err := ioutil.ReadDir(dir)
					if err != nil {
						continue
					}

					for _, file := range files {
						if file.IsDir() == false {
							fmt.Printf("File: %s\t%s\n", file.Name(), filepath.Join(dir, file.Name()))
							soundfiles[file.Name()] = filepath.Join(dir, file.Name())
						}
					}
				}

				fmt.Printf("GoMumbleSoundboard loaded (%d files)\n", len(soundfiles))
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

				r := gin.Default()
				r.Use(static.Serve("/", staticAssetsFS{fs: Assets}))

				r.GET("/files.json", func(c *gin.Context) {
					keys := make([]string, 0, len(soundfiles))
					for k := range soundfiles {
						keys = append(keys, k)
					}
					// Sort keys into alphabetical order. Sick of things moving around
					ss := sort.StringSlice(keys)
					ss.Sort()
					c.JSON(200, ss)
				})
				r.GET("/play/:file", func(c *gin.Context) {
					file, ok := soundfiles[c.Param("file")]
					if !ok {
						c.AbortWithError(404, fmt.Errorf("%s: file not found", c.Param("file")))
					}
					stream := gumbleffmpeg.New(e.Client, gumbleffmpeg.SourceFile(file))
					stream.Volume = volume
					err := stream.Play()
					if err != nil {
						c.AbortWithError(400, err)
					}
					c.String(200, fmt.Sprintf("Playing %s\n", file))
				})
				r.GET("/volume/:volume", func(c *gin.Context) {
					strVol := c.Param("volume")
					vol, err := strconv.Atoi(strVol)
					if err != nil {
						c.AbortWithError(400, fmt.Errorf("couldn't convert %s to integer: %v", strVol, err))
					}

					if vol < 0 && vol > 100 {
						c.AbortWithError(400, fmt.Errorf("number too small or too large: %s", strVol))
					}

					volume = float32(vol) / 100
					c.String(200, fmt.Sprintf("volume set to %d", vol))
				})
				r.GET("/stop", func(c *gin.Context) {
					stream.Stop()
					c.String(200, "ok")
				})
				r.Run(":3000")
			},
			Disconnect: func(e *gumble.DisconnectEvent) {
				os.Exit(1)
			},
		})
}
