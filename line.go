package line

import (
	"github.com/lorestudios/line/encoding"
	"github.com/nats-io/nats.go"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
	"time"
)

// Line contains a NATS connection and an encoding to be used for connections.
type Line struct {
	name     string
	conn     *nats.Conn
	encoding encoding.Encoding
}

// NewLine returns a new Line with an initialized connection and an encoding ready to use.
func NewLine(conf Config, enc encoding.Encoding) *Line {
	name := conf.Name
	for {
		nc, err := nats.Connect(conf.Address,
			nats.Token(conf.Token),
			nats.ReconnectWait(time.Second*5),
			nats.ClosedHandler(func(conn *nats.Conn) {
				logrus.Error("Could not reconnect to NATS server.")
			}),
			nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
				logrus.Errorf("Disconnected from NATS server: %v", err)
			}),
			nats.ReconnectHandler(func(conn *nats.Conn) {
				logrus.Info("Reconnected to NATS server!")
			}),
			nats.Name(name),
		)
		if err != nil {
			logrus.Errorf("Failed to connect to NATS, retrying in 5 seconds: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		return &Line{
			name:     name,
			conn:     nc,
			encoding: enc,
		}
	}
}

// Name returns the name of the NATS connection.
func (n *Line) Name() string {
	return n.name
}

// Conn returns the NATS connection for direct use.
func (n *Line) Conn() *nats.Conn {
	return n.conn
}

// ReadPacket reads a packet from the provided bytes and returns it.
func (n *Line) ReadPacket(b []byte) (packet.Packet, error) {
	return n.encoding.Decode(b)
}

// WritePacket writes a packet in to a byte slice which can be sent over NATS.
func (n *Line) WritePacket(pk packet.Packet) ([]byte, error) {
	return n.encoding.Encode(pk)
}
