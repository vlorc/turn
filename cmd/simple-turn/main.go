package main

import (
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pions/pkg/stun"
	"github.com/pions/turn"
)

type myTurnServer struct {
	usersMap map[string]string
}

func (m *myTurnServer) AuthenticateRequest(username string, srcAddr *stun.TransportAddr) (password string, ok bool) {
	if password, ok := m.usersMap[username]; ok {
		return password, true
	}
	return "", false
}

func makePort(port string) func() int {
	arr := strings.Split(port, "-")
	min, _ := strconv.Atoi(arr[0])
	max, _ := strconv.Atoi(arr[1])
	count := min
	return func() int {
		if count++; count > max {
			count = min
		}
		return count
	}
}

func main() {
	m := &myTurnServer{usersMap: make(map[string]string)}

	users := os.Getenv("USERS")
	if users == "" {
		log.Panic("USERS is a required environment variable")
	}
	for _, kv := range regexp.MustCompile(`(\w+)=(\w+)`).FindAllStringSubmatch(users, -1) {
		m.usersMap[kv[1]] = kv[2]
	}

	realm := os.Getenv("REALM")
	if realm == "" {
		log.Panic("REALM is a required environment variable")
	}

	udpPortStr := os.Getenv("UDP_PORT")
	if udpPortStr == "" {
		log.Panic("UDP_PORT is a required environment variable")
	}
	udpPort, err := strconv.Atoi(udpPortStr)
	if err != nil {
		log.Panic(err)
	}
	replyIp := os.Getenv("REPLY_IP")
	replyPort := os.Getenv("REPLY_PORT")

	turn.Start(turn.StartArguments{
		Server:    m,
		Realm:     realm,
		UDPPort:   udpPort,
		RelayIP:   net.ParseIP(replyIp),
		RelayPort: makePort(replyPort),
	})
}
