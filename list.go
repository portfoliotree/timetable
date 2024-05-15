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

	firstIndex, lastIndex := 0, 0

	if last := list[len(list)-1]; t0.After(last.time) {
		firstIndex = len(list)
	} else if i, ok := slices.BinarySearchFunc(list, Cell[Value]{time: t0}, Cell[Value].compareTimes); ok {
		firstIndex = i
	} else if i >= 0 && i < len(list) {
		firstIndex = i
	}

	if last := list[len(list)-1]; t1.After(last.time) {
		lastIndex = len(list)
	} else if i, ok := slices.BinarySearchFunc(list, Cell[Value]{time: t1}, Cell[Value].compareTimes); ok {
		lastIndex = i + 1
	} else if i >= 0 && i < len(list) {
		lastIndex = i
	}

	return list[firstIndex:lastIndex:lastIndex]
}
