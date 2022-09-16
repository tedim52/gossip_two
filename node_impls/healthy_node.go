package node_impls

import (
	"github.com/tedim52/gossip_two/node_interface/objects"
	"sync"
)

// Healthy Gossip Node implementation
type GossipNode struct {
	nodeID objects.NodeID

	currVal objects.Value
	
	database *objects.Database

	peers map[objects.NodeID]struct{}

	blacklist map[objects.NodeID]struct{}
	
	mutex sync.Mutex
}

func NewHealthyGossipNode() *GossipNode {
	return nil
}

func (n *GossipNode) Gossip() {
	clock := 0
	for {
		clock++
		if clock % 3 == 0 {
			// select random node from peers
			// dial node

			// err check
				// do necessary error handling
				// if dial doesn't work, add node id to blacklist

			// if successful
				// wait for response from node
				// read response into buffer
				// validate response from node
				// if response is not valid
					// do necessary error handling
				// if response is valid
					// upsert database
				// WHAT DO WE PRINT OUT HERE? ONLY THE DIFFERENT VALUES???
			// close the connection
		}
	}
}

func (n *GossipNode) Listen() {
	// setup listener
	for {
		// accept connections

		// once connection is received
		// send back serialized value of database
		// close the connection
	}
}


func (n *GossipNode) AddPeer(peer objects.NodeID) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// check that this node is not in the blacklist

	// dial node

	// err check
		// do necessary error handling
		// if dial doesn't work, add node id to blacklist

	// if successful
		// wait for response from node
		// read response into buffer
		// validate response from node
		// if response is not valid
			// do necessary error handling
		// if response is valid
			// upsert database
			// add node to peer list
	
	// close the connection
	return nil
}

func (n *GossipNode) UpdateValue(v objects.Value) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// TODO: implement
}