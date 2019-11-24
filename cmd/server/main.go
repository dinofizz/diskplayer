package main

import (
	"github.com/dinofizz/diskplayer"
	"log"
)

func main() {
	diskplayer.ReadConfig(diskplayer.DEFAULT_CONFIG_NAME)
	ds := diskplayer.NewDiskplayerServer(nil, nil)
	e := ds.RunRecordServer()
	if e != nil {
		log.Fatal(e)
	}
}
