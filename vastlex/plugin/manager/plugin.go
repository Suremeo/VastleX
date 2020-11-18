package manager

import (
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"os/exec"
	"sync"
)

var Plugins []*Plugin
var Indexes map[string]int
var Mutex sync.Mutex

type Plugin struct {
	Name           string
	Version        int
	encoder        *actions.Encoder
	decoder        *actions.Decoder
	cmd            *exec.Cmd
	close          *sync.Once
	packetsHandled []int16
}

func New(cmd *exec.Cmd) (*Plugin, error) {
	Mutex.Lock()
	defer Mutex.Lock()
	out, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	in, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	return &Plugin{
		Name:    "Unknown",
		Version: 0,
		encoder: actions.NewEncoder(out),
		decoder: actions.NewDecoder(in),
		cmd:     cmd,
		close:   &sync.Once{},
	}, nil
}

func WriteAll(action actions.Action) {
	for _, plugin := range Plugins {
		log.DefaultLogger.Error(plugin.WriteAction(action))
	}
}

func (plugin *Plugin) readActions() {
	for {
		action, err := plugin.decoder.Read()
		if err != nil {
			log.DefaultLogger.Warn("Failed to read action from plugin, exiting the plugin")
			return
		}
		action.ID()
	}
}

func (plugin *Plugin) WriteAction(action actions.Action) error {
	return plugin.encoder.Encode(action)
}

func (plugin *Plugin) Close() {
	plugin.close.Do(func() {
		log.DefaultLogger.Error(plugin.cmd.Process.Kill())
	})
}
