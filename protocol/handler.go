package protocol

import (
	"github.com/nats-io/nats.go"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Handler is an interface implemented by packet handlers.
type Handler interface {
	// Handle is responsible for handling a packet. Both the message and the packet are passed so that the message
	// can be used to send a response.
	Handle(msg *nats.Msg, pk packet.Packet) error
}

// Handlers is a struct dealing with the handling of packets with a registry of handlers.
type Handlers struct {
	// handlers is a map of all registered handlers.
	handlers map[uint32]Handler
}

// NewHandlers returns a new Handlers ready to use.
func NewHandlers() *Handlers {
	return &Handlers{make(map[uint32]Handler)}
}

// FindHandler returns a Handler from the provided packet ID and a boolean indicating whether
// a handler registered for the packet ID exists or not.
func (h *Handlers) FindHandler(id uint32) (Handler, bool) {
	v, ok := h.handlers[id]
	return v, ok
}

// RegisterHandler registers a Handler for the provided packet ID.
func (h *Handlers) RegisterHandler(id uint32, v Handler) {
	h.handlers[id] = v
}
