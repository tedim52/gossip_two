package objects

type Timestamp string

type State int

type Value struct {
	Time Timestamp
	
	val State
}

func NewValue(t Timestamp, v Value) Value {
	// TODO: implement
	return Value{}
}

func (v Value) serialize() string {
	// TODO: implement
	return "unimplemented"
}

func  DeserializeValue(db string) Value {
	// TODO: implement
	return Value{}
}