package node_interface

import (
	"github.com/tedim52/gossip_two/node_interface/objects"
)

type GossipNode interface {

	// BoostrapNode starts the gossip node by:
	// - initiating gossip to other peers
	// - initiating listening for messages from other peers
	BoostrapNode()

	// AddPeer attempts to add a peer with [id] to the nodes peer list
	// so that it will be considered for future gossip exchanges
	AddPeer(id objects.NodeID) (error)

	GetValue() (objects.GossipValue)

	// UpdateValue updates the nodes current value to [va]
	UpdateValue(val int64)

	GetDatabase() (*objects.Database)
}