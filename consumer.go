package line

import (
	"fmt"
	"github.com/lorestudios/line/protocol"
	"github.com/nats-io/nats.go"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
)

// Consumer is a Reader that reads packets from a NATS subject.
type Consumer struct {
	line     *Line
	subject  string
	handlers *protocol.Handlers

	sub     *nats.Subscription
	channel chan *nats.Msg
}

// NewConsumer returns a new Consumer ready to use with the required data.
func NewConsumer(l *Line, subject string, handlers *protocol.Handlers) *Consumer {
	return &Consumer{
		line:     l,
		subject:  subject,
		handlers: handlers,
		channel:  make(chan *nats.Msg, 64),
	}
}

// ReadPacket reads a packet from the NATS consumer and returns the message object, the packet read and
// any errors that have occurred.
func (r *Consumer) ReadPacket() (*nats.Msg, packet.Packet, error) {
	msg := <-r.channel
	if len(msg.Data) == 0 {
		return msg, nil, fmt.Errorf("message empty")
	}
	pk, err := r.line.ReadPacket(msg.Data)
	return msg, pk, err
}

// Close deals with any cleanup that needs to be done after the Start loop is escaped.
func (r *Consumer) Close() {
	if r.sub.IsValid() {
		_ = r.sub.Unsubscribe()
	}
	close(r.channel)
	if !r.line.Conn().IsClosed() {
		_ = r.line.Conn().Drain()
	}
}

// Start subscribes to the subject and starts reading and handling packets from it.
func (r *Consumer) Start() error {
	sub, err := r.line.Conn().ChanSubscribe(r.subject, r.channel)
	if err != nil {
		logrus.Errorf("could not subscribe to subject %s: %v", r.subject, err)
		return err
	}
	r.sub = sub

	defer r.Close()
	for {
		msg, pk, err := r.ReadPacket()
		if err != nil {
			logrus.Debugf("error reading packet %T: %v\n", pk, err)
			return err
		}
		go func() {
			h, ok := r.handlers.FindHandler(pk.ID())
			if !ok {
				logrus.Errorf("unhandled packet %T", pk)
				return
			}
			if err := h.Handle(msg, pk); err != nil {
				logrus.Errorf("client unable to handle packet: %v", err)
			}
		}()
	}
}
