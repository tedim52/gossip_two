package node_interface

type GossipNode interface {
	Gossip()
	
	Listen()

	AddPeer(id NodeID)

	UpdateValue(v Value)
}