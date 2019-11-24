package main

import (
	"github.com/dinofizz/diskplayer"
	"log"
)

func main() {
	diskplayer.ReadConfig(diskplayer.DEFAULT_CONFIG_NAME)
	s := &diskplayer.RealDiskplayerServer{}
	e := s.RunRecordServer()
	if e != nil {
		log.Fatal(e)
	}
}
