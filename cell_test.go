package timetable_test

import (
	"encoding/json"
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

func TestCell_JSON(t *testing.T) {
	emptyTable := new(Table)
	emptyTable = emptyTable.AddColumnFillMissingWithZero(List{})

	peach := List{
		elV(day0, 1), elV(day1, 2), elV(day2, 3),
	}

	t.Run("marshal", func(t *testing.T) {
		for _, tt := range []struct {
			Name     string
			List     List
			Expected string
		}{
			{Name: "nil", List: nil, Expected: `null`},
			{Name: "empty", List: List{}, Expected: `[]`},
			{Name: "some values", List: peach,
				// language=json
				Expected: `[
						{"time": "2022-10-20T00:00:00Z", "value": 1},
						{"time": "2022-10-21T00:00:00Z", "value": 2},
						{"time": "2022-10-24T00:00:00Z", "value": 3}
					]`,
			},
		} {
			t.Run(tt.Name, func(t *testing.T) {
				data, err := json.Marshal(tt.List)
				assert.NoError(t, err)
				assert.JSONEq(t, tt.Expected, string(data))
			})
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		for _, tt := range []struct {
			Name     string
			JSON     string
			Expected List
		}{
			{Name: "nil", Expected: nil, JSON: `null`},
			{Name: "empty", Expected: List{}, JSON: `[]`},
			{Name: "some values", Expected: peach,
				// language=json
				JSON: `[
						{"time": "2022-10-20T00:00:00Z", "value": 1},
						{"time": "2022-10-21T00:00:00Z", "value": 2},
						{"time": "2022-10-24T00:00:00Z", "value": 3}
					]`,
			},
		} {
			t.Run(tt.Name, func(t *testing.T) {
				var got List
				err := json.Unmarshal([]byte(tt.JSON), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, got)
			})
		}
	})
}
