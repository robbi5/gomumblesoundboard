# GoMumbleSoundboard
A Soundboard for the [Mumble](http://mumble.info) voice chat software written in [Go](http://golang.org).
![gomumblesoundboard](https://cloud.githubusercontent.com/assets/172415/5604064/6524d658-93a5-11e4-8009-beb179c03b81.png)
## Requirements

* go >= 1.6
* mumble server
* folder with sounds
* ffmpeg (`brew install ffmpeg` / `sudo apt-get install libav-tools`\*)
* opus-header (`brew install opus` / `sudo apt-get install libopus-dev`)

\* On ubuntu you may need to symlink ffmpeg to avconv: `sudo ln -s /usr/bin/avconv /usr/bin/ffmpeg`

Tested on OS X 10.11 and Ubuntu 14.04 LTS.

## Install

    $ go get github.com/robbi5/gomumblesoundboard

## Usage

    cd $GOPATH/src/github.com/robbi5/gomumblesoundboard
    $GOPATH/bin/gomumblesoundboard --server yourmumbleserver.com:64738 --insecure --channel ChannelName ~/SoundboardFiles/*.mp3

Then open [http://localhost:3000](http://localhost:3000) and press all the buttons!

### Supported command line arguments

* `--server localhost:64738` Mumble server address
* `--username gumble-bot` client username
* `--password hunter2` client password
* `--insecure` skip server certificate verification
* `--certificate` user certificate file (PEM)
* `--key` user certificate key file (PEM)
* `--channel ChannelName` Mumble channel to join.  
  If the channel is a sub channel, you need to enter the full path like `Parent/ChannelName`

### Supported environment variables

* `HOST`
* `PORT` (default: 3000)

## Development

    $ git clone https://github.com/robbi5/gomumblesoundboard.git
    $ cd gomumblesoundboard
    $ go build

For updating/editing dependencies, use [godep](https://github.com/tools/godep).

## License

MIT

## Thanks
Thanks to @bontibon / @layeh for [gumble](https://github.com/layeh/gumble) and @codegangsta for [martini](https://github.com/go-martini/martini).
