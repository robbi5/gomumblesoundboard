# GoMumbleSoundboard
A Soundboard for the [Mumble](http://mumble.info) voice chat software written in [Go](http://golang.org).
![gomumblesoundboard](https://cloud.githubusercontent.com/assets/172415/5604064/6524d658-93a5-11e4-8009-beb179c03b81.png)
## Requirements

* mumble server
* folder with sounds
* godep (`go get github.com/tools/godep`)
* ffmpeg (`brew install ffmpeg` / `sudo apt-get install libav-tools`*)
* opus-header (`brew install opus` / `sudo apt-get install libopus-dev`)

\* On ubuntu you may need to symlink ffmpeg to avconv: `sudo ln -s /usr/bin/avconv /usr/bin/ffmpeg`

## Install
Using [godep](https://github.com/tools/godep) to handle dependency version locking:

    $ git clone https://github.com/robbi5/gomumblesoundboard
    $ cd gomumblesoundboard
    $ godep restore
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
