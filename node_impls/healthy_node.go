package node_impls

import (
	"github.com/tedim52/gossip_two/node_interface/objects"

	"sync"
	"time"
	"net"
	"errors"
)

const (
	defaultInitValue = 0
)

type GossipMessage struct {
	serializedDB []byte
}

// Healthy Gossip Node implements a node that shares its own database to peers and pulls other peers' database, merging it into its
// own to implement database consistency via a pull gossip method.
// 
// Invariants:
// - Value associated with [nodeID] in [database] must always be equivalent to [currVal]
// - The max number of [nodeID]'s in [database], with the same ip address (different port number) should be three
// - Once something is added to [blacklist], it can't be removed.
type GossipNode struct {
	nodeID objects.NodeID

	currVal objects.GossipValue
	
	database *objects.Database

	peers map[objects.NodeID]struct{}

	blacklist map[objects.NodeID]struct{}
	
	mutex sync.Mutex
}

func NewHealthyGossipNode(ip string, port string) *GossipNode {
	nodeID := objects.NewNodeID(ip, port)
	initValue := objects.NewGossipValue(time.Now(), 0)
	db := objects.InitializeDatabase()
	db.SetGossipValue(nodeID, initValue)

	return &GossipNode {
		nodeID: nodeID, 
		currVal: initValue,
		database: db,
		peers: map[objects.NodeID]struct{}{},
		blacklist: map[objects.NodeID]struct{}{},
	}
}

func (n *GossipNode) BoostrapNode(){
	go n.listen()
	go n.gossip()
}

// gossip initiates the sending of gossip messages to
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

	// Check that this node is not in the blacklist
	if _, exists := n.blacklist[peer]; exists {
		return errors.New("Error adding peer. Peer was blacklisted.")
	}

	// dial node
	_, _ = net.Dial("tcp", peer.NodeID)

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
	n.mutex.Lock()
	defer n.mutex.Unlock()

	return n.currVal
}

func (n *GossipNode) UpdateValue(v int64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	gossipValue := objects.NewGossipValue(time.Now(), v)
	n.database.SetGossipValue(n.nodeID, gossipValue)
	n.currVal = gossipValue
}

func (n* GossipNode) GetDatabase() *objects.Database {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	return n.database
}