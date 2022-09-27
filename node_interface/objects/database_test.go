package objects

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestSerializeDatabaseWithOneEntry(t *testing.T){
	db := InitializeDatabase()
	nodeID := NewNodeID("127.0.0.1", "8080")
	time, _ := stringTimeToTime("1664228446")
	gossipVal := NewGossipValue(time, 4)
	db.SetGossipValue(nodeID, gossipVal)

	expectedDBStr := "127.0.0.1:8080,1664228446,4\n"
	dbStr := db.Serialize()

	require.Equal(t, expectedDBStr, dbStr)	
}

// weak test -> no order enforced on the iterating of db entries
// so nondeterministic
func TestSerializeDatabaseWithOneMultipleEntries(t *testing.T){
	db := InitializeDatabase()

	nodeIDOne := NewNodeID("127.0.0.1", "8080")
	timeOne, _ := stringTimeToTime("1664228446")
	gossipValOne := NewGossipValue(timeOne, 4)

	nodeIDTwo := NewNodeID("121.104.230.38", "3000")
	timeTwo, _ := stringTimeToTime("1663218247")
	gossipValTwo := NewGossipValue(timeTwo, 7)

	nodeIDThree := NewNodeID("60.60.164.141", "4001")
	timeThree, _ := stringTimeToTime("1664228459")
	gossipValThree := NewGossipValue(timeThree, 1234)

	db.SetGossipValue(nodeIDOne, gossipValOne)
	db.SetGossipValue(nodeIDTwo, gossipValTwo)
	db.SetGossipValue(nodeIDThree, gossipValThree)

	expectedDBStr := "127.0.0.1:8080,1664228446,4\n121.104.230.38:3000,1663218247,7\n60.60.164.141:4001,1664228459,1234\n"

	dbStr := db.Serialize()

	require.Equal(t, expectedDBStr, dbStr)	
}

// weak test only verifies no error and size of database
func TestDeserializeDatabaseReturnsCorrectDatabase(t *testing.T) {
	dbStr := "127.0.0.1:8080,1664228446,4\n121.104.230.38:3000,1663218247,7\n60.60.164.141:4001,1664228459,1234\n"
	
	db, err := DeserializeDatabase(dbStr)

	require.NoError(t, err)
	require.Equal(t, db.Size(), 3)
}

func TestDeserializeNodeIDReturnsInvalidDatabaseFormatForNoComma(t *testing.T) {
	dbStr := "127.0.0.1:8080,16642284464\n121.104.230.38:3000,1663218247,7\n60.60.164.141:4001,1664228459,1234\n"
	
	_, err := DeserializeDatabase(dbStr)

	require.Error(t, err)
}

func TestDeserializeNodeIDReturnsInvalidDatabaseFormatForNoNewLine(t *testing.T) {
	dbStr := "127.0.0.1:8080,16642284464121.104.230.38:3000,1663218247,7\n60.60.164.141:4001,1664228459,1234\n"
	
	_, err := DeserializeDatabase(dbStr)

	require.Error(t, err)
}

// weak test only verifies size
func TestDoesNotSetGossipValueWithFutureTime(t *testing.T){
	db := InitializeDatabase()
	nodeID := NewNodeID("127.0.0.1", "8080")

	// time that is after current time (obviously a flaky test, bc the time may be hit), need to subt out time.Now()
	time, _ := stringTimeToTime("18111164237052") 
	gossipVal := NewGossipValue(time, 4)

	// should not set this value
	db.SetGossipValue(nodeID, gossipVal)

	require.Equal(t, db.Size(), 0)	
}

// this test should be based on [maxNumPortsPerIP]
// also weak, only checks number of entries
func TestDoesNotSetGossipValueMoreThanMaxPortsPerIP(t *testing.T){
	db := InitializeDatabase()

	nodeIDOne := NewNodeID("127.0.0.1", "8080")
	timeOne, _ := stringTimeToTime("1664228446")
	gossipValOne := NewGossipValue(timeOne, 4)

	nodeIDTwo := NewNodeID("127.0.0.1", "3000")
	timeTwo, _ := stringTimeToTime("1663218247")
	gossipValTwo := NewGossipValue(timeTwo, 7)

	nodeIDThree := NewNodeID("127.0.0.1", "4008")
	timeThree, _ := stringTimeToTime("1664228459")
	gossipValThree := NewGossipValue(timeThree, 1234)

	nodeIDFour := NewNodeID("127.0.0.1", "4001")
	timeFour,  _ := stringTimeToTime("1664228459")
	gossipValFour := NewGossipValue(timeFour, 1234)

	db.SetGossipValue(nodeIDOne, gossipValOne)
	db.SetGossipValue(nodeIDTwo, gossipValTwo)
	db.SetGossipValue(nodeIDThree, gossipValThree)

	// should not set this Gossip Value
	db.SetGossipValue(nodeIDFour, gossipValFour)

	require.Equal(t, 3, db.Size())	
}

// Weak test
func TestUpsert(t *testing.T) {
	db := InitializeDatabase()
	dbTwo := InitializeDatabase()

	nodeIDOne := NewNodeID("127.0.0.1", "8080")
	timeOne, _ := stringTimeToTime("1664228446")
	gossipValOne := NewGossipValue(timeOne, 4)

	nodeIDTwo := NewNodeID("127.0.0.1", "3000")
	timeTwo, _ := stringTimeToTime("1663218247")
	gossipValTwo := NewGossipValue(timeTwo, 7)

	db.SetGossipValue(nodeIDOne, gossipValOne)
	dbTwo.SetGossipValue(nodeIDTwo, gossipValTwo)

	db.Upsert(dbTwo)

	require.Equal(t, db.Size(), 2)	
}