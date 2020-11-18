package manager

import (
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"os"
	"os/exec"
	"path/filepath"
)

// Init initializes the plugin system as the proxy and loads all the plugins.
func Init() {
	if f, err := os.Stat("./plugins"); !os.IsNotExist(err) {
		if f.IsDir() {
			if err := filepath.Walk("./plugins", func(name string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					loadPlugin(name)
				}
				return nil
			}); err == nil {
				return
			}
		}
	}
	log.DefaultLogger.Fatal(os.Mkdir("plugins", os.ModePerm))
	log.DefaultLogger.Debug("Plugin directory was not found and has been created.")
}

// loadPlugin loads the specified plugin from the plugins folder.
func loadPlugin(name string) {
	Mutex.Lock()
	defer Mutex.Unlock()
	log.DefaultLogger.Debug("Loading plugin '" + name + "'...")
	cmd := exec.Command("./plugins/" + name)
	err := cmd.Start()
	if err != nil {
		log.DefaultLogger.Warn("Plugin '" + name + "' failed to load: " + err.Error())
		return
	}
}