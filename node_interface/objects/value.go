package objects

import (
	"time"
)

type Timestamp time.Time

type GossipValue struct {
	time time.Time
	
	value int64
}

func NewGossipValue(t Timestamp, v int64) GossipValue {
	// TODO: implement
	return GossipValue{}
}

func (v GossipValue) serialize() string {
	// TODO: implement
	return "unimplemented"
}

func  DeserializeGossipValue(db string) GossipValue {
	// TODO: implement
	return GossipValue{}
}

func (v GossipValue) GetStringTime() string {
	return ""
}

func (v GossipValue) GetIntTime() int64 {
	return 0
}

func (v GossipValue) GetValue() int64 {
	return v.value
}