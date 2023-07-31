package encoding

import (
	"bytes"
	"fmt"
	protocol2 "github.com/lorestudios/line/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NBTEncoding is an NBT encoder/decoder for packets. It is the main encoder/decoder used over NATS, except for
// web services.
type NBTEncoding struct {
	buf  *bytes.Buffer
	pool protocol2.Pool
}

// NewNBTEncoding returns a new NBT encoder/decoder ready to use.
func NewNBTEncoding(pool protocol2.Pool) *NBTEncoding {
	return &NBTEncoding{
		buf:  bytes.NewBuffer(make([]byte, 0, 256)),
		pool: pool,
	}
}

// Encode ...
func (n *NBTEncoding) Encode(pk packet.Packet) ([]byte, error) {
	id := pk.ID()
	if err := protocol.WriteVaruint32(n.buf, id); err != nil {
		return nil, err
	}
	pk.Marshal(protocol.NewWriter(n.buf, 0))

	defer n.buf.Reset()
	return n.buf.Bytes(), nil
}

// Decode ...
func (n *NBTEncoding) Decode(b []byte) (packet.Packet, error) {
	buf := bytes.NewBuffer(b)

	var id uint32
	if err := protocol.Varuint32(buf, &id); err != nil {
		return nil, err
	}

	pk, ok := n.pool[id]
	if !ok {
		return nil, fmt.Errorf("unknown packet %v", id)
	}

	pk.Marshal(protocol.NewReader(buf, 0, false))
	if buf.Len() > 0 {
		return nil, fmt.Errorf("still have %v bytes unread", buf.Len())
	}

	return pk, nil
}
