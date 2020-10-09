package main

import (
	"github.com/VastleLLC/VastleX/log"
	"github.com/VastleLLC/VastleX/vastlex"
)

func main() {
	err := vastlex.Start()
	if err != nil {
		log.FatalError("VastleX crashed!", err)
	}
}
