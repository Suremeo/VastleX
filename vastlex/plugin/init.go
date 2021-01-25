package plugin

import (
	"github.com/VastleLLC/VastleX/vastlex"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
	"os"
	"sync"
)

var Encoder = actions.NewEncoder(os.Stdout)
var Decoder = actions.NewDecoder(os.Stdin)

var pluginVastleXinit = &pluginVastleX{
	players: &sync.Map{},
}

func Init(name string, version int32, handleEvents bool) {
	Name = name
	Version = version
	pluginVastleXinit.handleEvents = handleEvents
	vastlex.VastleX = pluginVastleXinit
	log.DefaultLogger = &logger{}

	err := WriteAction(&actions.Init{
		InitAction: &protobuf.InitAction{
			Name:    name,
			Version: version,
		},
	})
	if err != nil {
		log.DefaultLogger.Warn("Failed to send init action, exiting the plugin")
		os.Exit(0)
		return
	}
	action, err := Decoder.Read()
	if err != nil {
		log.DefaultLogger.Warn("Failed to read action from server, exiting the plugin")
		os.Exit(0)
		return
	}
	if action.ID() != actions.IDInfo {
		log.DefaultLogger.Warn("Packet sent before the info packet, exiting plugin.")
		os.Exit(0)
		return
	}
	handleInfo(action)
	go handleReading()
}

func WriteAction(action actions.Action) error {
	return Encoder.Encode(action)
}

func handleReading() {
	for {
		action, err := Decoder.Read()
		if err != nil {
			log.DefaultLogger.Warn("Failed to read action from server, exiting the plugin")
			os.Exit(0)
			return
		}
		if handles[action.ID()] != nil {
			handles[action.ID()](action)
		}
	}
}
