package objects

type IPAddress string

type Port string

type NodeID {
	IP IPAddress

	Port Port 
}

func NewNodeID(ip IPAddress, port Port) NodeID {
	// TODO: implement
	return nil
}

func (id NodeID) Serialize() return (string) {
	// TODO: implement
	return ""
}

func DeserializeNodeID(id string) return (NodeID) {
	// TODO: implement
	return nil
}