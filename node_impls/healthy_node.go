
// Healthy Gossip Node implementation
type GossipNode struct {
	nodeID NodeID

	currVal Value
	
	database *Database

	peers map[NodeID]struct{}

	blacklist map[NodeID]struct{}
	
	mutex sync.Mutex
}

func (n *GossipNode) Gossip() {
	int clock = 0
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


func (n *GossipNode) AddPeer(<ip:tcp> peer) {
	mutex.Lock()
	defer mutex.Release()

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
}

func (n *GossipNode) UpdateValue(val) {
	mutex.Lock()
	defer mutex.Release()

	currentValue = val
}