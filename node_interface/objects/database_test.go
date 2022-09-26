package objects

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestSerializeDatabaseWithOneEntry(t *testing.T){
	db := InitializeDatabase()

	expectedDBStr := "12341543143141234,1234"
	dbStr := gossipVal.Serialize()

	require.Equal(t, expectedDBStr, dbStr)	
}

func TestSerializeDatabaseWithOneMultipleEntries(t *testing.T){
	db := InitializeDatabase()

	expectedDBStr := "12341543143141234,1234"
	dbStr := gossipVal.Serialize()

	require.Equal(t, expectedDBStr, dbStr)	
}


func TestDeserializeDatabaseReturnsCorrectDatabase(t *testing.T) {
	gossipValStr := "12341543143141234,1234"
	
	expectedTimeStr := "12341543143141234"
	expectedValStr := "1234"

	gossipVal, err := DeserializeDatabase(gossipValStr)

	require.NoError(t, err)
	require.Equal(t, expectedTimeStr, gossipVal.GetTimeString())
	require.Equal(t, expectedValStr, gossipVal.GetValueString())
}

func TestDeserializeNodeIDReturnsInvalidDatabaseFormat(t *testing.T) {
	invalidGossipValStr:= "hello,1234212314"

	_, err := DeserializeDatabase(invalidGossipValStr)

	require.Error(t, err, InvalidDatabaseFormat)
}

func TestDeserializeNodeIDReturnsInvalidDatabaseFormat(t *testing.T) {
	invalidGossipValStr:= "12341543143141234,1212asfafda"

	_, err := DeserializeDatabase(invalidGossipValStr)

	require.Error(t, err, InvalidDatabaseFormat)
}

func TestDeserializeNodeIDReturnsInvalidDatabaseFormat(t *testing.T) {
	invalidGossipValStr:= "12341543143141234hello"

	_, err := DeserializeNodeID(invalidGossipValStr)

	require.Error(t, err, InvalidDatabaseFormat)
}