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
)

const (
	numLinesToRead = 256
	timeBetweenGossips = 3
)

// Healthy Gossip Node implements a node that shares its own database to peers and pulls other peers' database, merging it into its
// own to implement database consistency via a pull gossip method.
// 
// Invariants:
// - The max number of [nodeID]'s in [database], with the same ip address (different port number) should be three
// - Once something is added to [blacklist], it can't be removed.
type GossipNode struct {
	nodeID objects.NodeID
	
	database *objects.Database

	peers map[objects.NodeID]struct{}

	blacklist map[objects.NodeID]struct{}
	
	mutex sync.Mutex
}

func NewHealthyGossipNode(ip string, port string) *GossipNode {
	nodeID := objects.NewNodeID(ip, port)
	db := objects.InitializeDatabase()

	return &GossipNode {
		nodeID: nodeID, 
		database: db,
		peers: make(map[objects.NodeID]struct{}),
		blacklist: make(map[objects.NodeID]struct{}),
	}
}

func (n *GossipNode) BoostrapNode(){
	// start listening on this node
	go n.listen()

	// start gossiping every [timeBetweenGossips] seconds
	go func(){
		for range time.Tick(timeBetweenGossips * time.Second) {
			n.gossip()
		}
	}()
}

// gossip initiates the sending of gossip messages to
func (n *GossipNode) gossip() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// tracks errors throughout gossip
	var err error

	peer, found := n.getRandomPeerNodeID()
	if !found {
		return
	}
	// Check that this node is not in the blacklist
	if _, found := n.blacklist[peer]; found {
		return
	}

	// Dial node
	conn, err := net.Dial("tcp", peer.Serialize())
	// err check
	if err != nil {
		fmt.Println(err.Error())
		// if dial doesn't work, add node id to blacklist
		n.blacklist[peer] = struct{}{}
		return
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// read response into buffer
	reader := bufio.NewReader(conn)
	var messageBuffer []byte
	lineCounter := 0
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
				fmt.Println(err.Error())
				return
			}
		}
		messageBuffer = append(messageBuffer, bytes...)
		lineCounter++
	}

	// validate response from node
	peerDBStr := string(messageBuffer)
	peerDB, err := objects.DeserializeDatabase(peerDBStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// upsert database
	n.database.Upsert(peerDB)

	// close the connection
	conn.Close()
}

func (n *GossipNode) listen() {
	// setup listener
	ln, err := net.Listen("tcp", n.nodeID.Serialize())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer ln.Close()
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// once connection is received, send back serialized GossipValue of database
		if _, err = conn.Write([]byte(n.database.Serialize())); err != nil {
			fmt.Println(err.Error())
		}

		// close the connection
		conn.Close()
	}
}

func (n *GossipNode) AddPeer(peer objects.NodeID) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// Check that this node is not in the blacklist
	if _, found := n.blacklist[peer]; found {
		return errors.New("Error adding peer. Peer was blacklisted.")
	}

	// Dial node
	conn, err := net.Dial("tcp", peer.Serialize())
	if err != nil {
		n.blacklist[peer] = struct{}{}
		return err
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		n.blacklist[peer] = struct{}{}
		return err
	}

	// add node to peer set
	n.peers[peer] = struct{}{}

	// read response into buffer
	reader := bufio.NewReader(conn)
	var messageBuffer []byte
	lineCounter := 0
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

	// validate response from node
	peerDBStr := string(messageBuffer)
	peerDB, err := objects.DeserializeDatabase(peerDBStr)
	if err != nil {
		return err
	}

	// upsert database
	n.database.Upsert(peerDB)

	// close the connection
	conn.Close()
	
	return nil
}

func (n *GossipNode) UpdateValue(v int64) {
	gossipValue := objects.NewGossipValue(time.Now(), v)
	n.database.SetGossipValue(n.nodeID, gossipValue)
}

func (n *GossipNode) GetDatabase() *objects.Database {
	return n.database
}

// Invariant: 
// 	- [peers] cannot equal 0
func (n *GossipNode) getRandomPeerNodeID() (objects.NodeID, bool) {
	if len(n.peers) == 0 {
		return objects.NodeID{}, false
	}
	var nodeID objects.NodeID
	for id, _ := range n.peers {
		nodeID = id
	}
	return nodeID, true
}