package objects

import (
	"fmt"
	"errors"
	"regexp"
	"strings"
)

const (
	nodeIDDelimeter = ":"
	portRegexStr = "^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([0-5]{0,5})|([0-9]{1,4}))$"
	ipAddressRegexStr = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
)
var (
	ipAddressRegexPat = regexp.MustCompile(ipAddressRegexStr)
	portRegexPat = regexp.MustCompile(portRegexStr)

	InvalidIPAddress = errors.New("Invalid IP adderss format. Please provide a valid IPv4 or IPv6 address.")
	InvalidPortNumber = errors.New("Invalid TCP port format. Please provide a valid TCP port.")
	InvalidNodeID = errors.New("Invalid node id format. Please provide a node id in the following format '<ip-address>:<port>'")
)

type IPAddress string

type Port string

type NodeID struct {
	IP IPAddress

	Port Port 
	
	NodeID string
}

// Creates a NodeID representation
// Invariants:
// 	- [ip] and [port] must have correct corresponding formats, according to regexes
func NewNodeID(ip string, port string) NodeID {
	return NodeID{
		IP: IPAddress(ip),
		Port: Port(port),
		NodeID: fmt.Sprintf("%s:%s", string(ip), string(port)),
	}
}

func (id NodeID) Serialize() string {
	return id.NodeID
}

// Deserializes a NodeID string in the following format '<ip-address>:<port>'
// Returns error if format is incorrect
func DeserializeNodeID(idStr string) (NodeID, error) {
	nodeIDStrList := strings.Split(idStr, nodeIDDelimeter)
	if len(nodeIDStrList) == 1 || len(nodeIDStrList) > 2  {
		return NodeID{}, InvalidNodeID
	}
	ipAddresss := nodeIDStrList[0]
	port := nodeIDStrList[1]
	if !ipAddressRegexPat.Match([]byte(ipAddresss)) {
		return NodeID{}, InvalidIPAddress
	}
	if !portRegexPat.Match([]byte(port)) {
		return NodeID{}, InvalidPortNumber
	}
	return NewNodeID(ipAddresss, port), nil
}
