package node_impls

import (
	"github.com/tedim52/gossip_two/node_interface/objects"

	"fmt"
	"bufio"
	"sync"
	"time"
	"net"
	"errors"
	"io"
	"math/rand"
)

const (
	defaultInitValue = 0
	numLinesToRead = 256
)

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
		peers: make(map[objects.NodeID]struct{}),
		blacklist: make(map[objects.NodeID]struct{}),
	}
}

func (n *GossipNode) BoostrapNode(){
	go n.listen()
	go n.gossip()
}

// gossip initiates the sending of gossip messages to
func (n *GossipNode) gossip() {
	fmt.Println("starting to gossip...")
	clock := 0
	for {
		clock++
		if clock % 3 == 0 {
			n.mutex.Lock()
			// tracks errors throughout gossip
			var err error
			// select random node from peers
			peer := n.getRandomPeerNodeID()
	
			fmt.Println("attempting to gossip with a random peer...")
			// Dial node
			conn, err := net.Dial("tcp", peer.Serialize())
			// err check
			if err != nil {
				// do necessary error handling
				// if dial doesn't work, add node id to blacklist
				fmt.Println("error dialing peer, adding peer to blacklist and removing from peers...")
				fmt.Println(err.Error())
				n.blacklist[peer] = struct{}{}
				delete(n.peers, peer)
				continue
			}
		
			fmt.Println("reading response from peer...")
			// read response into buffer
			reader := bufio.NewReader(conn)
			var messageBuffer []byte
			lineCounter := 0
			fmt.Println("reading message from peer...")
			for {
				if lineCounter == numLinesToRead {
					break
				}
				bytes, err := reader.ReadBytes(byte('\n'))
				if err != nil {
					if err == io.EOF {
						messageBuffer = append(messageBuffer, bytes...)
						break
					} else {
						break
					}
				}
				messageBuffer = append(messageBuffer, bytes...)
				lineCounter++
			}
			if err != nil && err != io.EOF {
				// this means that smth went wrong and we actually don't want to do any gossip so keep moving clock
				fmt.Println(err.Error())
				continue
			}
		
			fmt.Println("validating message from peer...")
			// validate response from node
			peerDBStr := string(messageBuffer)
			// if response is not valid
				// do necessary error handling
			peerDB, err := objects.DeserializeDatabase(peerDBStr)
			if err != nil {
				fmt.Println("errors in message from from peer...")
				continue
			}
		
			fmt.Println("successfully retrieved message from peer, updating this nodes database now...")
			// upsert database
			n.database.Upsert(peerDB)
		
			fmt.Println("closing connection to peer...")
			// close the connection
			conn.Close()
			n.mutex.Unlock()
		}
	}
}

func (n *GossipNode) listen() {
	fmt.Println("starting to listen...")

	// setup listener
	fmt.Println("setting up a listener on this nodes port...")
	ln, err := net.Listen("tcp", n.nodeID.Serialize())
	if err != nil {
		fmt.Println("error occurred setting up listener for node")
	}
	defer ln.Close()
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error occurred accepting a connection on port for node.")
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("accepted a connection on port, now gossiping database...")

		// once connection is received
		// send back serialized GossipValue of database
		if _, err = conn.Write([]byte(n.database.Serialize())); err != nil {
			fmt.Println("error occurred while gossiping database...")
			fmt.Println(err.Error())
		}
		fmt.Println("successfully sent database to node that connected with me!")

		// close the connection
		conn.Close()
	}
}

func (n *GossipNode) AddPeer(peer objects.NodeID) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Println("checking if peeris in blacklist...")
	// Check that this node is not in the blacklist
	if _, exists := n.blacklist[peer]; exists {
		return errors.New("Error adding peer. Peer was blacklisted.")
	}

	fmt.Println("dialing peers...")
	// Dial node
	conn, err := net.Dial("tcp", peer.Serialize())
	// err check
		// do necessary error handling
		// if dial doesn't work, add node id to blacklist
	if err != nil {
		fmt.Println("error dialing peer, adding peer to blacklist and removing from peers...")
		n.blacklist[peer] = struct{}{}
		delete(n.peers, peer)
		return err
	}

	fmt.Println("reading response from peer...")
	// read response into buffer
	reader := bufio.NewReader(conn)
	var messageBuffer []byte
	lineCounter := 0
	fmt.Println("reading message from peer...")
	for {
		if lineCounter == numLinesToRead {
			break
		}
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF {
				messageBuffer = append(messageBuffer, bytes...)
				break
			} else {
				return err
			}
		}
		messageBuffer = append(messageBuffer, bytes...)
		lineCounter++
	}

	fmt.Println("validating message from peer...")
	// validate response from node
	peerDBStr := string(messageBuffer)
	// if response is not valid
		// do necessary error handling
	peerDB, err := objects.DeserializeDatabase(peerDBStr)
	if err != nil {
		fmt.Println("errors in message from from peer...")
		return err
	}

	fmt.Println("successfully retrieved message from peer, updating this nodes database now...")
	// upsert database
	n.database.Upsert(peerDB)

	fmt.Println("adding this node to peer list...")
	// add node to peer list
	n.peers[peer] = struct{}{}
	
	fmt.Println("closing connection to peer...")
	// close the connection
	conn.Close()
	
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
	// this value should always be set given database invariants, but we check in case
	if n.database.SetGossipValue(n.nodeID, gossipValue) {
		n.currVal = gossipValue
	}
}

func (n *GossipNode) GetDatabase() *objects.Database {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	return n.database
}

func (n *GossipNode) getRandomPeerNodeID()  objects.NodeID {
	k := rand.Intn(len(n.peers))
	peerNodeIDs := make([]objects.NodeID, len(n.peers))
	for nodeID, _ := range(n.peers) {
		peerNodeIDs = append(peerNodeIDs, nodeID)
	}
	return peerNodeIDs[k]
}