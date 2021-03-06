package turn

import (
	"fmt"

	"net"

	"github.com/pions/pkg/stun"
	"github.com/pions/turn/internal/client"
	"github.com/pions/turn/internal/server"
)

// Server is the interface for Pion TURN server callbacks
type Server interface {
	AuthenticateRequest(username string, srcAddr *stun.TransportAddr) (password string, ok bool)
}

// Client is the interface for Pion STUN requests
type Client interface {
	SendSTUNRequest(serverIP net.IP, serverPort int) (interface{}, error)
}

// StartArguments are the arguments for the Pion TURN server
type StartArguments struct {
	Server    Server
	Realm     string
	UDPPort   int
	RelayIP   net.IP
	RelayPort func() int
}

// ClientArguments are the arguments for the Pion client
type ClientArguments struct {
	BindingAddress string
	ServerIP       net.IP
	ServerPort     int
}

// Start the Pion TURN server
func Start(args StartArguments) {
	fmt.Println(server.NewServer(
		args.Realm,
		args.Server.AuthenticateRequest,
		args.RelayIP,
		args.RelayPort,
	).Listen("0.0.0.0", args.UDPPort))
}

// StartClient starts a Pion client
func StartClient(args ClientArguments) (interface{}, error) {
	c, err := client.NewClient(args.BindingAddress)
	if err != nil {
		return nil, err
	}
	return c.SendSTUNRequest(args.ServerIP, args.ServerPort)
}
