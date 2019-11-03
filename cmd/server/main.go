package main

import (
	"github.com/dinofizz/diskplayer/internal/config"
	"github.com/dinofizz/diskplayer/internal/server"
)


func main() {
	config.ReadConfig()
	server.RunRecordServer()
}

