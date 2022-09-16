package objects

type IPAddress string

type Port string

type NodeID struct {
	IP IPAddress

	Port Port 
}

func NewNodeID(ip IPAddress, port Port) NodeID {
	// TODO: implement
	return NodeID{}
}

func (id NodeID) Serialize() string {
	// TODO: implement
	return ""
}

func DeserializeNodeID(id string) NodeID {
	// TODO: implement
	return NodeID{}
}