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

// Healthy Gossip Node implements a node that shares its own database to peers and pulls other peers' database, merging it into its
// own to implement database consistency via a pull gossip method.
// 
// Invariants:
// - The max number of [nodeID]'s in [database], with the same ip address (different port number) should be three
// - Once something is added to [blacklist], it can't be removed.
type BadGossipNode struct {
	nodeID objects.NodeID
	
	database *objects.Database

	peers map[objects.NodeID]struct{}

	blacklist map[objects.NodeID]struct{}
	
	mutex sync.Mutex
}

func NewAdverserialGossipNode(ip string, port string) *BadGossipNode {
	nodeID := objects.NewNodeID(ip, port)
	db := objects.InitializeDatabase()

	return &BadGossipNode {
		nodeID: nodeID, 
		database: db,
		peers: make(map[objects.NodeID]struct{}),
		blacklist: make(map[objects.NodeID]struct{}),
	}
}

func (n *BadGossipNode) BoostrapNode(){
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
func (n *BadGossipNode) gossip() {
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

func (n *BadGossipNode) listen() {
	// setup listener
	ln, err := net.Listen("tcp", n.nodeID.Serialize())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer ln.Close()
	
	for {
		_, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// ADVERSERIAL MODES (UNCOMMENT WHICH ADVERSERIAL FEATURE TO USE)

		// 1. ACCEPT CONNECTION AND DON'T SEND BACK ANYTHING

		// 2. SEND BACK RANDOM CHARACTERS
		// // once connection is received, send back serialized GossipValue of database
		// if _, err = conn.Write([]byte("3214oi2klc ;kr d,.mnfqcew")); err != nil {
		// 	fmt.Println(err.Error())
		// }

		// // close the connection
		// conn.Close()

		// 3. CORRECT DATABASE FORMAT BUT DANGEROUS GOSSIP INFO (INCORRECT IP AND FUTURE TIMESTAMP)
		// // once connection is received, send back serialized GossipValue of database
		// if _, err = conn.Write([]byte("211.66.250.91:8080,1964282751,89\n")); err != nil {
		// 	fmt.Println(err.Error())
		// }

		// // close the connection
		// conn.Close()
	}
}

func (n *BadGossipNode) AddPeer(peer objects.NodeID) error {
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

func (n *BadGossipNode) UpdateValue(v int64) {
	gossipValue := objects.NewGossipValue(time.Now(), v)
	n.database.SetGossipValue(n.nodeID, gossipValue)
}

func (n *BadGossipNode) GetDatabase() *objects.Database {
	return n.database
}

// Invariant: 
// 	- [peers] cannot equal 0
func (n *BadGossipNode) getRandomPeerNodeID() (objects.NodeID, bool) {
	if len(n.peers) == 0 {
		return objects.NodeID{}, false
	}
	var nodeID objects.NodeID
	for id, _ := range n.peers {
		nodeID = id
	}
	return nodeID, true
}