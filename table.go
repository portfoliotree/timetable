package timetable

import (
	"encoding/json"
	"slices"
	"sort"
	"time"
)

type Compact[Value any] struct {
	times  []time.Time
	values [][]Value
}

func New[Value any](columns ...List[Value]) *Compact[Value] {
	table := new(Compact[Value])
	for _, column := range columns {
		table = table.AddColumnFillMissingWithZero(column)
	}
	return table
}

type encodedCompact[Value any] struct {
	Times  []time.Time `json:"times"  bson:"times"  yaml:"times"`
	Values [][]Value   `json:"values" bson:"values" yaml:"values"`
}

func (table *Compact[Value]) UnmarshalJSON(buf []byte) error {
	var enc encodedCompact[Value]
	err := json.Unmarshal(buf, &enc)
	table.times = enc.Times
	table.values = enc.Values
	table.sort()
	return err
}

type sorter struct {
	less func(a, b int) bool
	swap func(a, b int)
	len  int
}

func (s sorter) Len() int           { return s.len }
func (s sorter) Less(i, j int) bool { return s.less(i, j) }
func (s sorter) Swap(i, j int)      { s.swap(i, j) }

func (table *Compact[Value]) sort() {
	sort.Sort(sorter{
		less: func(a, b int) bool { return table.times[a].Before(table.times[b]) },
		swap: func(a, b int) {
			table.times[a], table.times[b] = table.times[b], table.times[a]
			for i := range table.values {
				table.values[i][a], table.values[i][b] = table.values[i][b], table.values[i][a]
			}
		},
		len: len(table.times),
	})
}

func (table *Compact[Value]) MarshalJSON() ([]byte, error) {
	table.sort()
	return json.Marshal(encodedCompact[Value]{Times: table.times, Values: table.values})
}

func (table *Compact[Value]) Times() []time.Time { return slices.Clone(table.times) }
func (table *Compact[Value]) Values() [][]Value {
	result := make([][]Value, len(table.values))
	for i, values := range table.values {
		result[i] = slices.Clone(values)
	}
	return result
}

func (table *Compact[Value]) UnderlyingValues() [][]Value  { return table.values }
func (table *Compact[Value]) UnderlyingTimes() []time.Time { return table.times }
func (table *Compact[Value]) NumberOfColumns() int         { return len(table.values) }
func (table *Compact[Value]) NumberOfRows() int            { return len(table.times) }
func (table *Compact[Value]) NumberOfCells() int           { return len(table.values) * len(table.times) }
func (table *Compact[Value]) FirstTime() time.Time         { return table.times[0] }
func (table *Compact[Value]) LastTime() time.Time          { return table.times[len(table.times)-1] }
func (table *Compact[Value]) isEmpty() bool                { return table == nil || table.times == nil }

func (table *Compact[Value]) Column(column int) (List[Value], bool) {
	if column < 0 || column >= len(table.values) {
		var zero List[Value]
		return zero[:], false
	}
	list := make(List[Value], len(table.times))
	for row := range table.times {
		list[row].time = table.times[row]
		list[row].value = table.values[column][row]
	}
	return list, true
}

func (table *Compact[Value]) Row(t time.Time) ([]Value, bool) {
	index, found := slices.BinarySearchFunc(table.times, t, time.Time.Compare)
	if !found {
		var zero [0]Value
		return zero[:], false
	}
	list := make([]Value, len(table.values))
	for column := range table.values {
		list[column] = table.values[column][index]
	}
	return list, true
}

func (table *Compact[Value]) AddColumnFillMissingWithZero(list List[Value]) *Compact[Value] {
	return table.AddColumn(list, zeroValue[Value])
}

func zeroValue[Value any](time.Time) Value {
	var zero Value
	return zero
}

func zeroTable[Value any](n int) *Compact[Value] {
	values := make([][]Value, n)
	var zeroValues [0]Value
	for i := range values {
		values[i] = zeroValues[:0:0]
	}
	var zeroTimes [0]time.Time
	return &Compact[Value]{
		times:  zeroTimes[:0:0],
		values: values,
	}
}

func (table *Compact[Value]) AddColumn(list List[Value], missing func(time.Time) Value) *Compact[Value] {
	if table.isEmpty() {
		return addInitialColumn(list)
	}
	if table.NumberOfRows() == 0 {
		return zeroTable[Value](len(table.values) + 1)
	}
	list = slices.Clone(list)
	slices.SortFunc(list, Cell[Value].compareTimes)
	t0, t1 := table.FirstTime(), table.LastTime()
	list = list.Between(t0, t1)
	if len(list) == 0 {
		return zeroTable[Value](len(table.values) + 1)
	}
	t0, t1 = list.FirstTime(), list.LastTime()
	updated := table.Between(t0, t1)
	if updated.NumberOfRows() == 0 {
		return zeroTable[Value](len(table.values) + 1)
	}
	return updated.addAdditionalColumn(list, missing)
}

func addInitialColumn[Value any](list List[Value]) *Compact[Value] {
	table := new(Compact[Value])
	slices.SortFunc(list, Cell[Value].compareTimes)
	newValues := make([]Value, 0, max(len(list), len(table.times)))
	table.times = make([]time.Time, 0, len(list))
	for _, element := range list {
		table.times = append(table.times, element.time)
		newValues = append(newValues, element.value)
	}
	table.values = append(table.values, newValues)
	return table
}

func (table *Compact[Value]) addAdditionalColumn(list List[Value], missing func(time.Time) Value) *Compact[Value] {
	var missingTimes []time.Time
	for _, cell := range list {
		_, found := slices.BinarySearchFunc(table.times, cell.time, time.Time.Compare)
		if !found {
			missingTimes = append(missingTimes, cell.time)
		}
	}
	times := slices.Grow(missingTimes, len(table.times)+len(missingTimes))
	for _, t := range table.times {
		times = append(times, t)
	}
	slices.SortFunc(times, time.Time.Compare)
	times = slices.Compact(times)
	values := make([][]Value, len(table.values)+1)

	for column := range table.values {
		values[column] = slices.Grow(table.values[column], len(times))
		values[column] = values[column][:0]
		for _, t := range times {
			index, found := slices.BinarySearchFunc(table.times, t, time.Time.Compare)
			var value Value
			if found {
				value = table.values[column][index]
			} else {
				value = missing(t)
			}
			values[column] = append(values[column], value)
		}
	}
	for _, t := range times {
		index, found := slices.BinarySearchFunc(list, Cell[Value]{time: t}, Cell[Value].compareTimes)
		var value Value
		if found {
			value = list[index].value
		} else {
			value = missing(t)
		}
		values[len(values)-1] = append(values[len(values)-1], value)
	}
	times = slices.Clip(times)
	for i := range values {
		values[i] = slices.Clip(values[i])
	}
	return &Compact[Value]{times: times, values: values}
}

func (table *Compact[Value]) Between(t0, t1 time.Time) *Compact[Value] {
	if table.isEmpty() || len(table.times) == 0 {
		return &Compact[Value]{
			times:  nil,
			values: make([][]Value, len(table.values)),
		}
	}
	if t1.Before(t0) {
		t0, t1 = t1, t0
	}

	var firstIndex, lastIndex int
	list := table.times
	if last := list[len(list)-1]; t0.After(last) {
		firstIndex = len(list)
	} else if first := list[0]; t0.Before(first) {
		lastIndex = 0
	} else if i, ok := slices.BinarySearchFunc(list, t0, time.Time.Compare); ok {
		firstIndex = i
	} else if i >= 0 && i < len(list) {
		firstIndex = i
	}

	if last := list[len(list)-1]; t1.After(last) {
		lastIndex = len(list)
	} else if first := list[0]; t1.Before(first) {
		lastIndex = 0
	} else if i, ok := slices.BinarySearchFunc(list, t1, time.Time.Compare); ok {
		lastIndex = i + 1
	} else if i >= 0 && i < len(list) {
		lastIndex = i
	}

	values := make([][]Value, len(table.values))
	for i := range table.values {
		values[i] = table.values[i][firstIndex:lastIndex:lastIndex]
	}
	return &Compact[Value]{
		times:  table.times[firstIndex:lastIndex:lastIndex],
		values: values,
	}
}
