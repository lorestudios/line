# Line

NATS API wrapper for inter-service communication.

# Example Usage

Below is a bare bones usage example of Line. For most use cases, you'll want to use the [client and consumer](#client-and-consumer).
```go
// Initialize pool
protocol.Register(IDExamplePacket, func () packet.Packet { return &ExamplePacket{} })
pool := protocol.NewPool()

// Create new line
c := line.DefaultConfig()
l := line.New(c, encoding.NewNBTEncoder(pool))

// Publish a packet
data, _ := nats.WritePacket(&ExamplePacket{})
_ = nats.Conn().Publish("subject", data)
```

# Packets

Communication is over packets using gophertunnel's packet implementation. Line comes with a packet handler registry
too, though the responsibility to handle packet read/writes is on the user.

### Packet Pool

```go
protocol.register(IDExamplePacket, &ExamplePacket{})
```

### Handlers

```go
// create a new handler registry and register a packet.
h := handler.NewHandlers()
h.RegisterHandler(IDExamplePacket, &ExamplePacketHandler{})

// find a handler and handle a packet.
exampleHandler, ok := h.FindHandler(IDExamplePacket)
if ok {
    exampleHandler.Handle(pk) // pk is a packet.Packet
}
```

# Client and Consumer
For ease of use, line comes with a client and consumer implementation.

### Client
```go
// create a new client
client := client.NewClient(line, "wild-1", handlers)
defer client.Close()
err := client.Start()

// Send data to another service.
err = client.Send("poet", &ExamplePacket{})

// Request data from another service.
// resp is a packet.Packet as well
resp, err := client.Request("poet", &ExampleRequestPacket{})
fmt.Println(resp.(*ExampleResponsePacket).Data)

// Respond to data from another service.
// this is usually inside a handler, msg is a *nats.Msg
err = client.Respond(msg, &ExampleResponsePacket{})
```
