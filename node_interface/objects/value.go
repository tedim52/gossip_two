package objects

type Timestamp string

type State int

type Value {
	Time Timestamp
	
	val State
}

func NewValue(t Timestamp, v Value) return (Value) {
	// TODO: implement
	return nil
}

func (v Value) serialize() return (string) {
	// TODO: implement
	return "unimplemented"
}
func  DeserializeValue(db string) return (Value) {
	// TODO: implement
	return nil
}