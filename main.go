package main

import (
	"github.com/VastleLLC/VastleX/vastlex/player"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func main() {
	listener := &minecraft.Listener{
		ServerName: "Vastlex",
	}
	address := "0.0.0.0:19132"
	err := listener.Listen("raknet", address)
	if err != nil {
		println(err.Error())
		return
	}
	println("Vastlex running on " + address)
	for {
		c, err := listener.Accept()
		if err != nil {
			println(err.Error())
			return
		}
		p := player.NewPlayer(c.(*minecraft.Conn))

		err = p.Send("127.0.0.1:1002")
		if err != nil {
			_ = p.Conn().WritePacket(&packet.Disconnect{Message: text.Red()("Unknown error while connecting you to a lobby!")})
			println("Error while connecting " + p.Conn().IdentityData().DisplayName + " (" + p.Conn().IdentityData().XUID + ") to a server: " + err.Error())
			continue
		} else {
			println(p.Conn().IdentityData().DisplayName + " (" + p.Conn().IdentityData().XUID + ") connected to the proxy and has been placed into a server")
		}
	}
}
