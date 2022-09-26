package objects

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestSerializeNodeID(t *testing.T){
	ip := "122.116.233.149"
	port := "8080"
	nodeID := NewNodeID(IPAddress(ip), Port(port))

	expectedNodeIDStr := "122.116.233.149:8080"
	nodeIDStr := nodeID.Serialize()

	require.Equal(t, expectedNodeIDStr, nodeIDStr)	
}

func TestDeserializeNodeIDReturnsNodeID(t *testing.T) {
	nodeIDStr := "122.116.233.149:8080"
	
	expectedIP := "122.116.233.149"
	expectedPort := "8080"

	nodeID, err := DeserializeNodeID(nodeIDStr)

	require.NoError(t, err)
	require.Equal(t, expectedPort, string(nodeID.Port))
	require.Equal(t, expectedIP, string(nodeID.IP))
}

func TestDeserializeNodeIDReturnsInvalidPortError(t *testing.T) {
	invalidPortNodeIDStr:= "122.116.233.149:1234212314"
	

	_, err := DeserializeNodeID(invalidPortNodeIDStr)

	require.Error(t, err, InvalidPortNumber)
}

func TestDeserializeNodeIDReturnsInvalidIPError(t *testing.T) {
	invalidIPNodeIDStr:= "42.42.42:8080"
	

	_, err := DeserializeNodeID(invalidIPNodeIDStr)

	require.Error(t, err, InvalidIPAddress)
}