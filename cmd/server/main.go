package main

import (
	"github.com/dinofizz/diskplayer"
)

func main() {
	diskplayer.ReadConfig(diskplayer.DEFAULT_CONFIG_NAME)
	diskplayer.RunRecordServer()
}
