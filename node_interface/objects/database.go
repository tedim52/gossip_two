package objects

import (
	"fmt"
)

var(
	InvalidDatabaseFormat = errors.New("Invalid database format.")
)

// Database representes the database used by a gossip node mapping [NodeID]'s of its peers to
// [GossipValue] representing the nodes most up to date information on the contents of its peer.
// The database is serialized and deserialized based on the following format:
//
// Format:
//	NodeID1,GossipValue1
//  NodeID2,GossipValue2
//	NodeID3,GossipValue3
//
// where NodeID format is '<ip-address>:<port>' and GossipValue format is '<timestamp>,<value'
// An example serialized database looks like:
// 122.116.233.149:8080,1234154131241,123\n
// 121.104.230.38:3000,122134423,81\n
// 121.104.230.38:3001,1221344233,85\n
//
// Invariants:
// - Cannot be more than 3 [NodeID] entries with the same [IPAddress]
//	- aka cannot exist connections to more than 3 ports at the same IP
// - There should be no [GossipValue]'s in [db] that have a [time] later than this node's current time
type Database struct {
	db map[NodeID]GossipValue
}

func InitializeDatabase() *Database {
	return &Database{
		db: map[NodeID]GossipValue{},
	}
}

// GetGossipValue returns the GossipValue associated with [id] if a value exists along with true.
// If no entry is found associated with [id], an empty GossipValue and false is returned.
func (db *Database) GetGossipValue(id NodeID) (GossipValue, bool) {
	gossipVal, found := db.db[id]
	if !found {
		return Gossip{}, false
	}
	return gossipVal, true
}

func (db *Database) SetGossipValue(id NodeID, v GossipValue) {
	db.db[id] = v
}

func (db *Database) Serialize() string {
	// TODO: implement
	return ""
}

// DeserializeDatabase takes a string representing a database and returns a Database struct. 
// Returns error if the database string is an invalid format.
func DeserializeDatabase(string) (*Database, error) {
	// TODO: implement
	return nil
}

// Upsert takes in a [dbToUpsert] and merges the mappings in [db] with the mappings in [db] according to the following rules:
// If [dbToUpsert] contains a NodeID entry not in [db], add this entry to [db]
// If [dbToUpsert] contains a NodeID entry in [db], ONLY add this entry to [db] if the timestamp associated with the GossipValue is later than 
// the timestamp associated with GossipValue aready in [db]
func (db *Database) Upsert(dbToUpsert *Database) {
	// TODO: implement
}

func (db *Database) PrintDatabase() {
	// TODO: implement
}

// Serializes a single database entry into the following format: 'NodeID,GossipValue'
// ex. 122.116.233.149:8080,1234154131241,123
// Invariant: 
// 
func (db *Database) serializeDatabaseEntry(id NodeID) (string)  {
	return fmt.Sprintf("%s,%s")
}

// Takes in a string representing a database and checks that it follows
// the line protocol format
func validateDatabaseFormat(db string) bool {
	// TODO: implement
	return false
}