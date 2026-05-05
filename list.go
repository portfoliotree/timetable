package timetable

import (
	"slices"
	"time"
)

type List[Value any] []Cell[Value]

func (list List[Value]) LastTime() time.Time  { return list[len(list)-1].time }
func (list List[Value]) FirstTime() time.Time { return list[0].time }

// Between returns a slice of a list. You may want pass result into slices.Clone()
// before using any mutating functions (such as Insert).
func (list List[Value]) Between(t0, t1 time.Time) List[Value] {
	if len(list) == 0 {
		return nil
	}
	if t1.Before(t0) {
		t0, t1 = t1, t0
	}
	slices.SortFunc(list, Cell[Value].compareTimes)
	last := list[len(list)-1].time

	var firstIndex, lastIndex int
	if t0.After(last) {
		firstIndex = len(list)
	} else {
		firstIndex, _ = slices.BinarySearchFunc(list, Cell[Value]{time: t0}, Cell[Value].compareTimes)
	}

	switch i, ok := slices.BinarySearchFunc(list, Cell[Value]{time: t1}, Cell[Value].compareTimes); {
	case t1.After(last):
		lastIndex = len(list)
	case ok:
		lastIndex = i + 1
	default:
		lastIndex = i
	}

	return list[firstIndex:lastIndex:lastIndex]
}
