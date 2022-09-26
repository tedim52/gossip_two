package objects

import (
	"time"
	"strconv"
	"errors"
	"fmt"
	"strings"
)

const (
	gossipValDelimeter = ","
)
var(
	InvalidGossipValueFormat = errors.New("Invalid Gossip Value format.")
)

type GossipValue struct {
	time time.Time
	
	value int64
}

func NewGossipValue(t time.Time, v int64) GossipValue {
	return GossipValue{
		time: t,
		value: v,
	}
}

// Serialized gossip value into the following format 'time,value' where
// 	value is a string representing the int64 of value
// 	time is a string representing number of seconds after 1970
func (v GossipValue) Serialize() string {
	return fmt.Sprintf("%s%s%s", v.GetTimeString(), gossipValDelimeter, v.GetValueString())
}

// Deserializes a GossipValue string in the following format 'time,value'
// Returns error if format is incorrect
func  DeserializeGossipValue(valueStr string) (GossipValue, error) {
	gossipValStrList := strings.Split(valueStr, gossipValDelimeter)
	if len(gossipValStrList) == 1 || len(gossipValStrList) > 2  {
		return GossipValue{}, InvalidGossipValueFormat
	}
	timeStr := gossipValStrList[0]
	valStr:= gossipValStrList[1]
	timeInt, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return GossipValue{}, InvalidGossipValueFormat
	}
    time := time.Unix(timeInt, 0)
	valInt, err := strconv.ParseInt(valStr, 10, 32)
	if err != nil {
		return GossipValue{}, InvalidGossipValueFormat
	}
	return NewGossipValue(time, valInt), nil
}

func (v GossipValue) GetTime() time.Time {
	return v.time
}

func (v GossipValue) GetTimeString() string {
	return strconv.FormatInt(v.time.Unix(), 10)
}

func (v GossipValue) GetValue() int64 {
	return v.value
}

func (v GossipValue) GetValueString() string {
	return fmt.Sprint(v.value)
}

// CompareTime returns true if GossipValue [a] precedes GossipValue [b]
// this is determined based on which value was gossiped first
func CompareTime(a GossipValue, b GossipValue) bool{
	timeA := a.GetTime()
	timeB := b.GetTime()
	return timeA.Before(timeB)
}

func stringTimeToTime(timeStr string) (time.Time, error){
	timeInt, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Now(), err // change time now
	}
    time := time.Unix(timeInt, 0)
	return time, nil
}