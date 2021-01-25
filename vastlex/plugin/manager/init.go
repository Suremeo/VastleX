package manager

import (
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Ready = make(chan struct{})
var vastlex interfaces.VastleX

// Init initializes the plugin system as the proxy and loads all the plugins.
func Init(v interfaces.VastleX) {
	vastlex = v
	Indexes = map[string]int{}
	if f, err := os.Stat("./plugins"); !os.IsNotExist(err) {
		if f.IsDir() {
			if err := filepath.Walk("./plugins", func(name string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					loadPlugin(strings.TrimPrefix(strings.TrimPrefix(name, "plugins/"), "plugins\\"))
				}
				return nil
			}); err == nil {
				log.DefaultLogger.Info("Loaded " + strconv.Itoa(len(Plugins)) + " plugin(s).")
				return
			}
		}
	}
	log.DefaultLogger.Fatal(os.Mkdir("plugins", os.ModePerm))
	log.DefaultLogger.Debug("Plugin directory was not found and has been created.")
}
