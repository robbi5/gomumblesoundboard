# GoMumbleSoundboard
This is a fork of [gomumblesoundboard](https://github.com/robbi5/gomumblesoundboard), which is a soundboard for [mumble](https://www.mumble.info/).

## Requirements
* go >= 1.14 (dunno, but probably a current one)
* mumble server
* folder with sounds
* ffmpeg (`brew install ffmpeg` / `sudo apt-get install ffmpeg`)
* opus-header (`brew install opus` / `sudo apt-get install libopus-dev`)

Tested on debian bullseye

## Install
```
go get github.com/feuerrot/gomumblesoundboard
```

## Usage
```
$GOPATH/bin/gomumblesoundboard --server yourmumbleserver.com:64738 ~/soundboard
```

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
* `--maxvol 50` maximum volume in %

## Development
```
git clone https://github.com/feuerrot/gomumblesoundboard.git
cd gomumblesoundboard
go build .
```

## License
MIT

## Thanks to
* @bontibon / @layeh for [gumble](https://github.com/layeh/gumble)
* various people for [gin](https://github.com/gin-gonic/gin)
* @robbi5 for the [original version](https://github.com/robbi5/gomumblesoundboard) of this software
