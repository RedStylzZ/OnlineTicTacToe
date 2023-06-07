package main

import (
	"time"

	"github.com/RedStylzZ/OnlineTicTacToe/logging"
	"github.com/RedStylzZ/OnlineTicTacToe/p2p"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&logging.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func main() {
	setting := p2p.Settings{
		Name:       "Server",
		Version:    "1.0",
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(setting)
	log.Infof("Listening on %s", setting.ListenAddr)
	go server.Start()
	time.Sleep(time.Second)

	cSetting := p2p.Settings{
		Name:       "Client",
		Version:    "1.0",
		ListenAddr: ":4000",
	}
	log.Infof("Listening on %s", cSetting.ListenAddr)
	cServer := p2p.NewServer(cSetting)
	go cServer.Start()
	err := cServer.Connect(setting.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
