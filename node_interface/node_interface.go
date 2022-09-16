package node_interface

import (
	"github.com/tedim52/gossip_two/node_interface/objects"
)

type GossipNode interface {
	Gossip()
	
	Listen()

	AddPeer(id NodeID) (error)

	UpdateValue(v Value)
}