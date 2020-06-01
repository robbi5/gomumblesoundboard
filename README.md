# GoMumbleSoundboard
This is a fork of [gomumblesoundboard](https://github.com/robbi5/gomumblesoundboard), which tries to integrate some minor fixes.

A Soundboard for the [Mumble](http://mumble.info) voice chat software written in [Go](http://golang.org).
![gomumblesoundboard](https://cloud.githubusercontent.com/assets/172415/19899199/7921df8e-a05f-11e6-8545-13731eaacf10.png)


## Requirements
* go >= 1.14 (dunno, but probably a current one)
* mumble server
* folder with sounds
* ffmpeg (`brew install ffmpeg` / `sudo apt-get install libav-tools`\*)
* opus-header (`brew install opus` / `sudo apt-get install libopus-dev`)

Tested on debian bullseye

## Install
```
go get github.com/feuerrot/gomumblesoundboard
```

## Usage
```
$GOPATH/bin/gomumblesoundboard --server yourmumbleserver.com:64738 ~/soundboard/*
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

### Supported environment variables
* `HOST`
* `PORT` (default: 3000)

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
* @codegangsta for [martini](https://github.com/go-martini/martini)
* @robbi5 for the [original version](https://github.com/robbi5/gomumblesoundboard) of this software
