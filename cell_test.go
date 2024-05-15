package timetable_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/timetable"
)

type (
	Value = int
	Cell  = timetable.Cell[Value]
)

const (
	day0      = "2022-10-20" // Thursday
	day1      = "2022-10-21" // Friday
	day2      = "2022-10-24" // Monday
	day3      = "2022-10-25" // Tuesday
	firstDay  = day0
	lastDay   = day3
	dayBefore = "2022-10-19"
	dayAfter  = "2022-10-26"
)

func date(value string) time.Time {
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func elT(t string) Cell {
	return timetable.NewCell[Value](date(t), 0)
}

func elV(t string, v Value) Cell {
	return timetable.NewCell[Value](date(t), v)
}

func TestCell_Time(t *testing.T) {
	cell := elT(day0)
	assert.Equal(t, date(day0), cell.Time())
}

func TestCell_Value(t *testing.T) {
	cell := elV(day0, 50)
	assert.Equal(t, 50, cell.Value())
}

func TestCell_GoString(t *testing.T) {
	cell := elV(day0, 1)
	const s = `timetable.Cell[int]{time: "2022-10-20 00:00:00 +0000 UTC", value: 1}`
	assert.Equal(t, s, cell.GoString())
}
