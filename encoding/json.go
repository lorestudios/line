package encoding

import (
	"encoding/json"
	"fmt"
	"github.com/lorestudios/line/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// JsonEncoding is a simple JSON encoder/decoder for packets. It is used with communication with the web server as
// it is less performant but more flexible and easier to use.
type JsonEncoding struct {
	pool protocol.Pool
}

// NewJsonEncoding returns a new JSON encoder/decoder ready to use.
func NewJsonEncoding(pool protocol.Pool) *JsonEncoding {
	return &JsonEncoding{pool: pool}
}

// Encode ...
func (j *JsonEncoding) Encode(p packet.Packet) ([]byte, error) {
	return json.Marshal(p)
}

// Decode ...
func (j *JsonEncoding) Decode(b []byte) (packet.Packet, error) {
	var jp jsonPacket
	if err := json.Unmarshal(b, &jp); err != nil {
		return nil, err
	}
	pk, ok := j.pool[jp.ID]
	if !ok {
		return nil, fmt.Errorf("unknown packet %v", jp.ID)
	}
	if err := json.Unmarshal(b, &pk); err != nil {
		return nil, err
	}
	return pk, nil
}

// jsonPacket is a simple packet that only contains an ID. It is used to decode packets from JSON.
type jsonPacket struct {
	ID uint32
}
