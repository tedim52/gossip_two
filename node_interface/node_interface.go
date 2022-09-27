package node_interface

import (
	"github.com/tedim52/gossip_two/node_interface/objects"
)

type GossipNode interface {
	// BoostrapNode starts the gossip node by:
	// - initiating gossip to other peers
	// - initiating listening for messages from other peers
	// Gossip occurs in a pull, anti-entropy fashion meaning
	//	1. pull -> nodes will "prompt" or "query" other nodes to get their updates
	// 	2. anti-entropy -> when a node gets "prompted" or "queried", it will send the entirety of its database back
	BoostrapNode()

	// AddPeer attempts to add a peer with [id] to the nodes peer list
	// so that it will be considered for future gossip exchanges
	AddPeer(id objects.NodeID) (error)

	// UpdateValue updates the nodes current value to [val]
	UpdateValue(val int64)

	GetDatabase() (*objects.Database)
}