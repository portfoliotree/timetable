package timetable

import (
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
