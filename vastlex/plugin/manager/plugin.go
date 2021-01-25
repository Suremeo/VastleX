package manager

import (
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"io"
	"os"
	"os/exec"
	"sync"
)

var Plugins []*Plugin
var Indexes map[string]int
var Mutex sync.Mutex

type Plugin struct {
	Name           string
	Version        int32
	encoder        *actions.Encoder
	decoder        *actions.Decoder
	cmd            *exec.Cmd
	close          *sync.Once
	packetsHandled []int16
}

// loadPlugin loads the specified plugin from the plugins folder.
func loadPlugin(name string) *Plugin {
	log.DefaultLogger.Debug("Loading plugin file: " + name)
	cmd := exec.Command("./plugins/" + name)
	p, err := New(name, cmd)
	if err != nil {
		log.DefaultLogger.Warn("Plugin file '" + name + "' failed to load: " + err.Error())
		return nil
	}
	err = cmd.Start()
	if err != nil {
		log.DefaultLogger.Warn("Plugin file '" + name + "' failed to load: " + err.Error())
		return nil
	}
	action, err := p.decoder.Read()
	if err != nil {
		log.DefaultLogger.Warn("Plugin file '" + name + "' failed to load: " + err.Error())
		p.Close()
		return nil
	}
	if action.ID() != actions.IDInit {
		log.DefaultLogger.Warn("Plugin file '" + name + "' did not send the init action first.")
		p.Close()
		return nil
	} else {
		if handles[action.ID()] != nil {
			handles[action.ID()](p, action)
		}
		go p.readActions()
	}
	return p
}

func New(name string, cmd *exec.Cmd) (*Plugin, error) {
	Mutex.Lock()
	defer Mutex.Unlock()
	out, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	in, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	l, _ := os.Create(name + "-err.txt")
	cmd.Stderr = l
	l, _ = os.Create(name + "-out.txt")
	p := &Plugin{
		Name:    name,
		Version: 0,
		encoder: actions.NewEncoder(out),
		decoder: actions.NewDecoder(&LoggingReadCloser{
			from: in,
			to:   l,
		}),
		cmd:   cmd,
		close: &sync.Once{},
	}
	Indexes[name] = len(Plugins)
	Plugins = append(Plugins, p)
	return p, nil
}

func WriteAll(action actions.Action) {
	for _, plugin := range Plugins {
		log.DefaultLogger.Error(plugin.WriteAction(action))
	}
}

func (plugin *Plugin) readActions() {
	<-Ready
	for {
		action, err := plugin.decoder.Read()
		if err != nil {
			log.DefaultLogger.Warn("Failed to read action from plugin, exiting the plugin: " + err.Error())
			plugin.Close()
			return
		}
		if handles[action.ID()] != nil {
			handles[action.ID()](plugin, action)
		}
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

type LoggingReadCloser struct {
	from io.ReadCloser
	to   io.Writer
}

func (l *LoggingReadCloser) Read(b []byte) (i int, err error) {
	i, err = l.from.Read(b)
	_, _ = l.to.Write(b)
	return
}

func (l *LoggingReadCloser) Close() error {
	return l.from.Close()
}
