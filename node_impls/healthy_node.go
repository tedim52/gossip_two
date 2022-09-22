package node_impls

import (
	"github.com/tedim52/gossip_two/node_interface/objects"

	// "fmt"
	"sync"
	"time"
)

const (
	defaultInitValue = 0
)

// Healthy Gossip Node implementation
type GossipNode struct {
	nodeID objects.NodeID

	currVal objects.GossipValue
	
	database *objects.Database

	peers map[objects.NodeID]struct{}

	blacklist map[objects.NodeID]struct{}
	
	mutex sync.Mutex
}

func NewHealthyGossipNode(ip objects.IPAddress, port objects.Port) *GossipNode {

	return &GossipNode {
		nodeID: objects.NewNodeID(ip, port),
		currVal: objects.NewGossipValue(objects.Timestamp(time.Now()), 0),
		database: objects.InitializeDatabase(),
		peers: map[objects.NodeID]struct{}{},
		blacklist: map[objects.NodeID]struct{}{},
	}
}

func (n *GossipNode) BoostrapNode(){
	go n.listen()
	go n.gossip()
}

func (n *GossipNode) gossip() {
	// fmt.Println("starting to gossip...")
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
				// WHAT DO WE PRINT OUT HERE? ONLY THE DIFFERENT GossipValueS???
			// close the connection
		}
	}
}

func (n *GossipNode) listen() {
	// fmt.Println("starting to listen...")
	// setup listener
	for {
		// accept connections

		// once connection is received
		// send back serialized GossipValue of database
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

func (n* GossipNode) GetValue() objects.GossipValue {
	return n.currVal
}

func (n *GossipNode) UpdateValue(v int64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	
}

func (n* GossipNode) GetDatabase() *objects.Database {
	return n.database
}