package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"layeh.com/gumble/gumbleutil"
	_ "layeh.com/gumble/opus"
)

//go:embed public
var Assets embed.FS

type File struct {
	Name     string `json:"name"`
	Folder   string `json:"folder"`
	FullPath string
}

func (f File) String() string {
	return f.Folder + "/" + f.Name
}

var soundfiles map[string]File

func scanDirsFunc(l string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	validSuffix := []string{
		".mp3",
		".m4a",
		".ogg",
		".flac",
		".opus",
		".wav",
		".MPG",
	}
	validSuffixCheck := false
	for _, s := range validSuffix {
		if strings.HasSuffix(info.Name(), s) {
			validSuffixCheck = true
		}
	}
	if !validSuffixCheck {
		return nil
	}

	if info.IsDir() == false {
		fmt.Printf("File: %s\t%s\n", info.Name(), l)
		dir, file := path.Split(l)
		split := strings.Split(dir, "/")
		f := File{
			FullPath: l,
			Name:     file,
			Folder:   split[len(split)-2],
		}

		soundfiles[f.String()] = f
	}

	return nil
}

func scanDirs(directories []string) {
	soundfiles = make(map[string]File)
	for _, dir := range directories {
		err := filepath.Walk(dir, scanDirsFunc)
		if err != nil {
			fmt.Printf("Error at %s: %v", dir, err)
		}
	}
}

func main() {
	targetChannel := flag.String("channel", "Root", "channel the bot will join")
	maxVolume := flag.String("maxvol", "100", "Set the maximum Volume in %, the volume set in the UI is multiplied with it")
	var volume float32 = 1

	gumbleutil.Main(
		gumbleutil.AutoBitrate,
		gumbleutil.Listener{
			Connect: func(e *gumble.ConnectEvent) {
				stream := gumbleffmpeg.New(e.Client, nil)
				stream.Volume = volume
				scanDirs(flag.Args())

				e.Client.Self.SetSelfDeafened(true)

				maxVolumeF, err := strconv.Atoi(*maxVolume)
				if err != nil {
					fmt.Printf("Invalid MaxVolume %d", maxVolumeF)
					os.Exit(1)
				}
				maxvol := float32(maxVolumeF) / 100
				fmt.Printf("maximum Volume: %.1f%%\n", maxvol*100)

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

				r := gin.New()
				s := http.FileServer(http.FS(FSPrefixer("public", Assets)))
				r.NoRoute(func(context *gin.Context) {
					s.ServeHTTP(context.Writer, context.Request)
				})

				r.GET("/files.json", func(c *gin.Context) {
					files := make([]File, 0)
					for _, f := range soundfiles {
						files = append(files, File{
							Name:   f.Name,
							Folder: f.Folder,
						})
					}

					// Sort keys into alphabetical order. Sick of things moving around
					c.JSON(200, files)
				})
				r.GET("/play/:file", func(c *gin.Context) {
					unescape, err := url.PathUnescape(c.Param("file"))
					if err != nil {
						c.AbortWithError(404, err)
						return
					}
					f, ok := soundfiles[unescape]
					if !ok {
						c.AbortWithError(404, fmt.Errorf("%s: file not found", unescape))
						return
					}
					if stream.State() == gumbleffmpeg.StatePlaying {
						c.AbortWithError(400, fmt.Errorf("already playing a sound, gtfo"))
						return
					}
					e.Client.Self.SetSelfMuted(false)
					stream = gumbleffmpeg.New(e.Client, gumbleffmpeg.SourceFile(f.FullPath))
					stream.Volume = volume
					if err := stream.Play(); err != nil {
						c.AbortWithError(400, err)
						return
					}
					go func() {
						stream.Wait()
						e.Client.Self.SetSelfDeafened(true)
					}()
					c.String(200, fmt.Sprintf("Playing %s\n", f.FullPath))
				})
				r.GET("/volume/:volume", func(c *gin.Context) {
					strVol := c.Param("volume")
					vol, err := strconv.Atoi(strVol)
					if err != nil {
						c.AbortWithError(400, fmt.Errorf("couldn't convert %s to integer: %v", strVol, err))
						return
					}

					if vol < 0 && vol > 100 {
						c.AbortWithError(400, fmt.Errorf("number too small or too large: %s", strVol))
						return
					}

					volume = float32(vol) / 100 * maxvol
					c.String(200, fmt.Sprintf("volume set to %.1f%%", volume*100))
				})
				r.GET("/stop", func(c *gin.Context) {
					stream.Stop()
					c.String(200, "ok")
				})
				r.GET("/rescan", func(c *gin.Context) {
					scanDirs(flag.Args())
					c.Redirect(http.StatusTemporaryRedirect, "/")
				})
				r.GET("/restart", func(c *gin.Context) {
					os.Exit(255)
				})
				r.Run(":3000")
			},
			Disconnect: func(e *gumble.DisconnectEvent) {
				os.Exit(1)
			},
		})
}

type FuncFS func(name string) (fs.File, error)

func (f FuncFS) Open(name string) (fs.File, error) {
	return f(name)
}

func FSPrefixer(prefix string, f fs.FS) fs.FS {
	return FuncFS(func(name string) (fs.File, error) {
		if name == "." {
			name = ""
		}

		if name == "" {
			name = prefix
		} else {
			name = prefix + "/" + name
		}

		return f.Open(name)
	})
}
