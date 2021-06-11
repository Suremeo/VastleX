package minecraft

import (
	"bytes"
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/custompackets"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// The majority of the code in this file was from https://github.com/Sandertv/gophertunnel/blob/master/minecraft/packet.go and has been edited to fit our uses.

func init() {
	packet.Register(custompackets.IDVastleXTransfer, func() packet.Packet { return &custompackets.VastleXTransfer{} })
	packet.Register(packet.IDCraftingData, func() packet.Packet { return &custompackets.CraftingData{} })
	//packet.Register(packet.IDSetScore, func() packet.Packet { return &custompackets.NoParse{Id: packet.IDSetScore} })
	//packet.Register(packet.IDAddItemActor, func() packet.Packet { return &custompackets.AddItemActor{} })
}

// packetData holds the data of a Minecraft packet.
type packetData struct {
	h       *packet.Header
	full    []byte
	payload *bytes.Buffer
	c       *Connection
}

// parseData parses the packet data slice passed into a packetData struct.
func parseData(data []byte, c *Connection) (*packetData, error) {
	buf := bytes.NewBuffer(data)
	header := &packet.Header{}
	if err := header.Read(buf); err != nil {
		// We don't return this as an error as it's not in the hand of the user to control this. Instead,
		// we return to reading a new packet.
		return nil, fmt.Errorf("error reading packet header: %v", err)
	}
	return &packetData{c: c, h: header, full: data, payload: buf}, nil
}

// decode decodes the packet payload held in the packetData and returns the packet.Packet decoded.
func (p *packetData) decode() (pk packet.Packet, err error) {
	// Attempt to fetch the packet with the right packet ID from the Pool
	pkf, ok := p.c.Pool[p.h.PacketID]
	if !ok {
		// No packet with the ID. This may be a custom packet of some sorts.
		pk = &packet.Unknown{PacketID: p.h.PacketID}
	}
	pk = pkf()


	r := protocol.NewReader(p.payload, 0)
	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("%T: %w", pk, recoveredErr.(error))
		}
	}()
	pk.Unmarshal(r)
	if p.payload.Len() != 0 {
		return pk, fmt.Errorf("%T: %v unread bytes left: 0x%x", pk, p.payload.Len(), p.payload.Bytes())
	}
	return pk, nil
}
