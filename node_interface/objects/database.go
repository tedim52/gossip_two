package objects

import (
	"fmt"
	"errors"
	"strings"
	"time"
	"sync"
)

const (
	entryDelimeter = ","
	newEntryDelimeter = "\n"
	maxNumPortsPerIP = 3
)

var(
	InvalidDatabaseFormat = errors.New("Invalid database format.")
)

// Database representes a thread safe database implementation used by a gossip node mapping [NodeID]'s of its peers to
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
// - Cannot be more than [maxNumPortsPerIP] [NodeID] entries with the same [IPAddress]
//	- aka cannot exist connections to more than 3 ports at the same IP
// - There should be no [GossipValue]'s in [db] that have a [time] later than this node's current time
// - All IPAddress' in [ipToNumPorts] are associated with at least one IPAddress in a NodeID in [db]
type Database struct {
	db map[NodeID]GossipValue

	ipToNumPorts map[IPAddress]int

	mutex sync.RWMutex
}

func InitializeDatabase() *Database {
	return &Database{
		db: make(map[NodeID]GossipValue),
		ipToNumPorts: make(map[IPAddress]int),
	}
}

// GetGossipValue returns the GossipValue associated with [id] if a value exists along with true.
// If no entry is found associated with [id], an empty GossipValue and false is returned.
func (db *Database) GetGossipValue(id NodeID) (GossipValue, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	gossipVal, found := db.db[id]
	if !found {
		return GossipValue{}, false
	}
	return gossipVal, true
}

// TODO: this function violates single responsibility principle by both updating and returning whether or not smth was to be updated or not, thus there's a side effect. 
//	find better way to do this

// SetGossipValue sets [id] to [v] in the database unless the following is the case in order to abide by invariants:
// - There are already [maxNumPortsPerIP] ports associated with the same IPAddress in [db]
// - The time associated with [v] is past this nodes local time ("in the future")
// - If there is already a value in the [db] associated with [id], and the timestamp of the value is after the timestamp of [v]
// Returns true if something was updated, and false otherwise. 
func (db *Database) SetGossipValue(id NodeID, v GossipValue) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if v.GetTime().After(time.Now()) {
		return
	}
	currGossipVal, found := db.db[id]
	if found && currGossipVal.GetTime().After(v.GetTime()) {
		return
	}
	db.db[id] = v
	// THIS IS BAD THIS IS A SIDE EFFECT BUT IT GETS THE JOB DONE, PRINT ONLY EXACTLY WHEN AN UPDATE OCCURS
	// TODO: is there a more clean way to do this? maybe return updated node ids and print at the gossip or main level
	// only print if already in db and new value
	if found && currGossipVal.GetValue() != v.GetValue(){
		fmt.Println(fmt.Sprintf("%s --> %s", id.Serialize(), v.GetValueString()))
	}
}

func (db *Database) Serialize() string {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	dbStr := ""
	for id, _ := range db.db { 
		dbStr = dbStr + db.serializeDatabaseEntry(id) + "\n"
	}
	return dbStr
}

// DeserializeDatabase takes a [dbStr] representing a database and returns a Database struct. 
// Returns error if the database string is an invalid format.
func DeserializeDatabase(dbStr string) (*Database, error) {
	db := InitializeDatabase()
	dbEntryStrList := strings.Split(dbStr, newEntryDelimeter)
	dbEntryStrList = dbEntryStrList[:len(dbEntryStrList) - 1]
	for _, entryStr  := range(dbEntryStrList) {
		entryValueStrList := strings.SplitAfterN(entryStr, entryDelimeter, 2)
		if len(entryValueStrList) != 2 {
			return nil, InvalidDatabaseFormat
		}
		nodeIDStr := entryValueStrList[0]
		nodeIDStr = nodeIDStr[:len(nodeIDStr)-1]
		nodeID, err := DeserializeNodeID(nodeIDStr)
		if err != nil {
			// could turn this into a continue to be more liberal
			return nil, err
		}
		gossipValStr := entryValueStrList[1]
		gossipVal, err := DeserializeGossipValue(gossipValStr)
		if err != nil {
			// could turn this into a continue to be more liberal
			return nil, err
		}
		db.SetGossipValue(nodeID, gossipVal)
	}
	return db, nil
}

// Upsert takes in a [dbToUpsert] and merges the mappings in [db] with the mappings in [db] according to the following rules:
// If [dbToUpsert] contains a NodeID entry not in [db], add this entry to [db]
// If [dbToUpsert] contains a NodeID entry in [db], ONLY add this entry to [db] if the timestamp associated with the GossipValue is later than 
// the timestamp associated with GossipValue aready in [db]
func (db *Database) Upsert(dbToUpsert *Database) {
	for _, nodeID := range dbToUpsert.GetNodeIDs() {
		gossipVal, _ := dbToUpsert.GetGossipValue(nodeID)
		db.SetGossipValue(nodeID, gossipVal)
	}
}

func (db *Database) Size() int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return len(db.db)
}

// Retrieves all NodeIDs mapped in [db]
func (db *Database) GetNodeIDs() []NodeID {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	nodeIDs := make([]NodeID, len(db.db))
	i := 0
	for nodeID, _ := range db.db {
		nodeIDs[i] = nodeID
		i++
	}
	return nodeIDs
}

// Serializes a single database entry into the following format: 'NodeID,GossipValue'
// ex. 122.116.233.149:8080,1234154131241,123
// Invariant: 
// 	[id] must exist as an entry in [db]
func (db *Database) serializeDatabaseEntry(id NodeID) (string)  {
	value, _ := db.db[id]
	return fmt.Sprintf("%v,%v", id.Serialize(), value.Serialize())
}