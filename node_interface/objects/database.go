package objects

// Invariants:
// - Cannot be more than 3 entries with the same IPAddress 
//- aka cannot exist connections to more than 3 ports at the same IP
type Database map[NodeID]*Value


func InitializeDatabase(string) *Database {
	// TODO: implement
	return nil
}

func CreateDatabaseFromString(string) *Database {
	// TODO: implement
	return nil
}

func (d *Database) GetValue(id NodeID) return (Value) {
	// TODO: implement
	return nil
}

func (d* Database) SetValue(id NodeID, v Value) {
	// TODO: implement
}

func (d* Database) Serialize() return (string) {
	// TODO: implement
	return ""
}

func DeserializeDatabase(string) (*Database) {
	// TODO: implement
	return nil
}

func (d *Database) Upsert(d *Database) {
	// TODO: implement
	return nil
}

func (d* Database) PrintDatabase() {
	// TODO: implement
}

// Takes in a string representing a database and checks that it follows
// the line protocol format
func _validateDatabaseFormat(db string) return (bool) {
	// TODO: implement
	return false
}