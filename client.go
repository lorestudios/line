package line

import (
	"fmt"
	"github.com/lorestudios/line/protocol"
	"github.com/nats-io/nats.go"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"time"
)

// Client represents a client that can handle packets coming to it and write packets in response.
type Client struct {
	*Consumer
}

// NewClient returns a new Client ready to use with the required data.
func NewClient(l *Line, s string, h *protocol.Handlers) *Client {
	return &Client{
		Consumer: NewConsumer(l, s, h),
	}
}

// Send sends a packet to the given subject without waiting for a response.
func (c *Client) Send(subject string, pk packet.Packet) error {
	data, err := c.line.WritePacket(pk)
	if err != nil {
		return err
	}
	return c.line.Conn().Publish(subject, data)
}

// Request sends a request packet to the given subject and awaits a response. The response and any errors are returned.
func (c *Client) Request(subject string, pk packet.Packet) (packet.Packet, error) {
	data, err := c.line.WritePacket(pk)
	if err != nil {
		return nil, err
	}
	m, err := c.line.Conn().Request(subject, data, time.Second*5)
	if err != nil {
		return nil, fmt.Errorf("no response received from poet: %v", err)
	}
	p, err := c.line.ReadPacket(m.Data)
	return p, err
}

// Reply sends a response to a request packet. Any errors are returned.
func (c *Client) Reply(msg *nats.Msg, pk packet.Packet) error {
	data, err := c.line.WritePacket(pk)
	if err != nil {
		return err
	}
	return msg.Respond(data)
}
