package protocol

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// registeredPackets holds the registered packets.
var registeredPackets = map[uint32]func() packet.Packet{}

// Register registers a function that returns a packet for a specific ID. Packets with this ID coming in from
// connections will resolve to the packet returned by the function passed.
func Register(id uint32, pk func() packet.Packet) {
	registeredPackets[id] = pk
}

// Pool holds packets indexed by id.
type Pool map[uint32]packet.Packet

// NewPool returns a new pool with all supported packets.
func NewPool() Pool {
	p := Pool{}
	for id, pk := range registeredPackets {
		p[id] = pk()
	}
	return p
}
