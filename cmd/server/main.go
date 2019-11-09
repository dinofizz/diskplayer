package main

import (
	"github.com/dinofizz/diskplayer"
)

func main() {
	diskplayer.ReadConfig()
	diskplayer.RunRecordServer()
}
