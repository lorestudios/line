package encoding

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// Encoding is the interface that wraps the basic Encode and Decode methods.
type Encoding interface {
	// Encode encodes the given packet.
	Encode(pk packet.Packet) ([]byte, error)
	// Decode decodes the given data into a packet.
	Decode(b []byte) (packet.Packet, error)
}
