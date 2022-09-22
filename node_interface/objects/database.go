package objects

// Invariants:
// - Cannot be more than 3 entries with the same IPAddress 
//- aka cannot exist connections to more than 3 ports at the same IP
type Database struct {
	db map[NodeID]*GossipValue
}

func InitializeDatabase() *Database {
	// TODO: implement
	// db := map[NodeID]*GossipValue{}
	return &Database{
		db: map[NodeID]*GossipValue{},
	}
}

func CreateDatabaseFromString(dbString string) *Database {
	// TODO: implement
	return nil
}

func (db *Database) GetGossipValue(id NodeID) GossipValue {
	// TODO: implement
	return GossipValue{}
}

func (db *Database) SetGossipValue(id NodeID, v GossipValue) {
	// TODO: implement
}

func (db *Database) Serialize() string {
	// TODO: implement
	return ""
}

func DeserializeDatabase(string) (*Database) {
	// TODO: implement
	return nil
}

func (db *Database) Upsert(dbToUpsert *Database) {
	// TODO: implement
}

func (db *Database) PrintDatabase() {
	// TODO: implement
}

// Takes in a string representing a database and checks that it follows
// the line protocol format
func _validateDatabaseFormat(db string) bool {
	// TODO: implement
	return false
}