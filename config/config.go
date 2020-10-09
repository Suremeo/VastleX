package config

import (
	"crypto/rand"
	"encoding/hex"
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
	Fallback struct {
		Enabled bool
		Host    string
		Port    int
	}
	Lobby struct {
		Enabled bool
		Host    string
		Port    int
	}
	Debug struct {
		Logging bool
		Profiling bool
	}
	Proxy struct {
		Secret string
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
	log.Debugging = config.Debug.Logging
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

// DefaultConfig returns the default configuration for the proxy & generates a unique secret.
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
		Fallback: struct {
			Enabled bool
			Host    string
			Port    int
		}{
			Enabled: true,
			Host:    "127.0.0.1",
			Port:    19133,
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
		Debug: struct {
			Logging bool
			Profiling bool
		}{},
		Proxy: struct {
			Secret string
		}{
			Secret: func() string { // Generate a cryptographically secure secret for the configuration.
				b := make([]byte, 32)
				_, err := rand.Read(b)
				log.FatalError("Failed to generate default secret for the configuration", err)
				return hex.EncodeToString(b)
			}(),
		},
	}
}
