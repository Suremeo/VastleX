package actions

import (
	"encoding/json"
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

// Info is sent by the proxy to inform the client of any changes in the proxy information.
type Info struct {
	*protobuf.InfoAction
}

func (i *Info) ID() int16 {
	return IDInfo
}

func (i *Info) DecodeConfig() (config config.Structure) {
	_ = json.Unmarshal(i.Config, &config)
	return
}

func (i *Info) EncodeConfig(config config.Structure) *Info {
	i.Config, _ = json.Marshal(config)
	return i
}
