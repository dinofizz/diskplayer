# diskplayer

diskplayer is essentially a Spotify client, written in Go, built around [zmb3's existing Spotify Web API wrapper
](https://github.com/zmb3/spotify).

This software is part of the larger "diskplayer" project of mine: a 3.5" floppy disk music player, running on a
 Raspberry Pi.

There are two components of diskplayer. The "player" binary and the "recorder" binary. The player binary can be used to obtain a new long-lived Spotify client authentication token, as well as accepting Spotify URI or a path to a file containing a Spotify URI which it will attempt to play.
  
The recorder binary runs as an HTTP server, and exposes a simple UI which can be used to record a Spotify URI to a chosen location.

## Build

### Requirements

displayer was developed using Go 1.13, and uses go modules to install its required dependencies.



## Configuration

In this repository there exists a `diskplayer.yaml` configuration file which must be updated with the relevant config values.

diskplayer will search for the `diskplayer.yaml` configuration file in one of the following locations:

* `/etc/diskplayer/`
* `$HOME/.config/diskplayer/`
* or the current directory from which the `player` or `recorder` binary is being run.

For the Spotify-related config values see [the documentation for the zmb3's Spotify wrapper](https://github.com/zmb3/spotify#authentication). The callback URL must match that as configured when you set up your Spotify API application.

The `recorder.folder_path` configuration value represents to the folder to which the disk device will be mounted during the recording process. You will need to ensure that this folder exists.

## Player Usage

## Recorder Usage

tests


