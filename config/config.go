package config

import (
	"github.com/VastleLLC/VastleX/log"
	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"io/ioutil"
	"os"
)

var Config = LoadConfig()

// Config is a parsed configuration file for the proxy.
type Structure struct {
	Listener struct {
		Host string
		Port int
	}
	Minecraft struct {
		Auth        bool
		Motd        string
		ShowVersion bool
		MaxPlayers  int
	}
	Lobby struct {
		Enabled bool
		Host    string
		Port    int
	}
	Logging struct {
		Debug bool
	}
}

// LoadConfig loads and parses the configuration for the proxy.
func LoadConfig() (config Structure) {
	config, err := readConfig()
	if err != nil {
		config = DefaultConfig()
		err = createDefaultConfig()
		log.FatalError("Failed to create default configuration", err)
	}
	return config
}

// readConfig loads and parses the configuration for the proxy.
func readConfig() (conf Structure, err error) {
	dat, err := ioutil.ReadFile("./config.toml")
	if err != nil {
		return conf, err
	}
	err = toml.Unmarshal(dat, &conf)
	return
}

// createDefaultConfig writes the default config to the configuration file.
func createDefaultConfig() error {
	dat, err := toml.Marshal(DefaultConfig())
	if err != nil {
		return err
	}
	return ioutil.WriteFile("./config.toml", dat, os.ModePerm)
}

// DefaultConfig returns the default configuration for the proxy.
func DefaultConfig() Structure {
	return Structure{
		Listener: struct {
			Host string
			Port int
		}{
			Host: "0.0.0.0",
			Port: 19132,
		},
		Minecraft: struct {
			Auth        bool
			Motd        string
			ShowVersion bool
			MaxPlayers  int
		}{
			Auth:        true,
			Motd:        text.Bold()(text.Red()("Vastle")) + text.Bold()(text.White()("X")),
			ShowVersion: false,
			MaxPlayers:  0,
		},
		Lobby: struct {
			Enabled bool
			Host    string
			Port    int
		}{
			Enabled: true,
			Host:    "127.0.0.1",
			Port:    19133,
		},
		Logging: struct {
			Debug bool
		}{},
	}
}
