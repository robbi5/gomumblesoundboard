# GoMumbleSoundboard
A Soundboard for the [Mumble](http://mumble.info) voice chat software written in [Go](http://golang.org).
![gomumblesoundboard](https://cloud.githubusercontent.com/assets/172415/5604064/6524d658-93a5-11e4-8009-beb179c03b81.png)
## Requirements

* mumble server
* folder with sounds
* ffmpeg (`brew install ffmpeg`)
* opus-header (`brew install opus`)

## Install

    $ go get -u github.com/layeh/gumble/gumble
    $ go get -u github.com/go-martini/martini
    $ go build

## Usage

    ./gomumblesoundboard --server yourmumbleserver.com:64738 --insecure --channel ChannelName ~/SoundboardFiles/*.mp3

Then open [http://localhost:3000](http://localhost:3000) and press all the buttons!

### Supported command line arguments

* `--server localhost:64738` Mumble server address
* `--username gumble-bot` client username
* `--password hunter2` client password
* `--insecure` skip server certificate verification
* `--certificate` user certificate file (PEM)
* `--key` user certificate key file (PEM)
* `--channel ChannelName` Mumble channel to join

## License

MIT

## Thanks
Thanks to @bontibon / @layeh for [gumble](https://github.com/layeh/gumble) and @codegangsta for [martini](https://github.com/go-martini/martini).
