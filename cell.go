package timetable

import (
	"encoding/json"
	"fmt"
	"time"
)

type Cell[Value any] struct {
	time  time.Time
	value Value
}

func NewCell[Value any](time time.Time, value Value) Cell[Value] {
	return Cell[Value]{time: time, value: value}
}

func (c Cell[Value]) compareTimes(o Cell[Value]) int { return c.time.Compare(o.time) }

func (c Cell[Value]) Time() time.Time { return c.time }

func (c Cell[Value]) Value() Value { return c.value }

func (c Cell[Value]) GoString() string {
	return fmt.Sprintf("%T{time: %q, value: %#v}", c, c.time, c.value)
}

type encodedCell[Value any] struct {
	Time  time.Time `json:"time"  bson:"time"  yaml:"time"`
	Value Value     `json:"value" bson:"value" yaml:"value"`
}

func (c *Cell[Value]) UnmarshalJSON(buf []byte) error {
	var enc encodedCell[Value]
	err := json.Unmarshal(buf, &enc)
	c.time = enc.Time
	c.value = enc.Value
	return err
}

//goland:noinspection GoMixedReceiverTypes
func (c Cell[Value]) MarshalJSON() ([]byte, error) {
	return json.Marshal(encodedCell[Value]{Time: c.time, Value: c.value})
}
