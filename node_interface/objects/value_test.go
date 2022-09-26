package objects

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestSerializeGossipValue(t *testing.T){
	timeStr := "12341543143141234"
	time, _ := stringTimeToTime(timeStr)
	value := int64(1234)
	gossipVal := NewGossipValue(time, value)

	expectedGossipValStr := "12341543143141234,1234"
	gossipValStr := gossipVal.Serialize()

	require.Equal(t, expectedGossipValStr, gossipValStr)	
}

func TestDeserializeGossipValueReturnsCorrectGossipValue(t *testing.T) {
	gossipValStr := "12341543143141234,1234"
	
	expectedTimeStr := "12341543143141234"
	expectedValStr := "1234"

	gossipVal, err := DeserializeGossipValue(gossipValStr)

	require.NoError(t, err)
	require.Equal(t, expectedTimeStr, gossipVal.GetTimeString())
	require.Equal(t, expectedValStr, gossipVal.GetValueString())
}

func TestDeserializeNodeIDReturnsInvalidGossipValueFormatForIncorrecTime(t *testing.T) {
	invalidGossipValStr:= "hello,1234212314"

	_, err := DeserializeGossipValue(invalidGossipValStr)

	require.Error(t, err, InvalidGossipValueFormat)
}

func TestDeserializeNodeIDReturnsInvalidGossipValueFormatForIncorrectValue(t *testing.T) {
	invalidGossipValStr:= "12341543143141234,1212asfafda"

	_, err := DeserializeGossipValue(invalidGossipValStr)

	require.Error(t, err, InvalidGossipValueFormat)
}

func TestDeserializeNodeIDReturnsInvalidGossipValueFormatForNoDelimeter(t *testing.T) {
	invalidGossipValStr:= "12341543143141234hello"

	_, err := DeserializeNodeID(invalidGossipValStr)

	require.Error(t, err, InvalidGossipValueFormat)
}