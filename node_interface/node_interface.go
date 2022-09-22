package node_interface

import (
	"github.com/tedim52/gossip_two/node_interface/objects"
)

type GossipNode interface {
	BoostrapNode()

	AddPeer(id objects.NodeID) (error)

	GetValue() (objects.GossipValue)

	UpdateValue(int64)

	GetDatabase() (*objects.Database)
}