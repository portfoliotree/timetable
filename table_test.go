package timetable_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/timetable"
)

type Table = timetable.Compact[Value]

func TestNew(t *testing.T) {
	for _, tt := range []struct {
		Name  string
		Lists []List
		Then  func(t *testing.T, table *Table)
	}{
		{
			Name:  "with nil input",
			Lists: []List(nil),
			Then: func(t *testing.T, table *Table) {
				assert.Len(t, table.Times(), 0)
				assert.Equal(t, table.NumberOfColumns(), 0)
			},
		},
		{
			Name:  "with an empty list",
			Lists: []List{{}},
			Then: func(t *testing.T, table *Table) {
				assert.Len(t, table.Times(), 0)
				assert.Equal(t, table.NumberOfColumns(), 1)
			},
		},
		{
			Name:  "with one list",
			Lists: []List{{elT(day1), elT(day2)}},
			Then: func(t *testing.T, table *Table) {
				assert.Equal(t, 2, table.NumberOfRows())
				assert.Equal(t, 1, table.NumberOfColumns())
				assert.Equal(t, 2, table.NumberOfCells())
				assert.Equal(t, []time.Time{date(day1), date(day2)}, table.UnderlyingTimes())
				assert.Equal(t, [][]Value{{0, 0}}, table.UnderlyingValues())
			},
		},
		{
			Name: "with aligned single elements",
			Lists: []List{
				{elV(day0, 1)},
				{elV(day0, 2)},
				{elV(day0, 3)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Len(t, table.Times(), 1, "it returns one row")
				assert.Equal(t, 3, table.NumberOfColumns(), "it gives the correct column count")
				assert.Equal(t, 1, table.NumberOfRows(), "it gives the correct row count")
				assert.Equal(t, 3, table.NumberOfCells())
				values := table.Values()
				assert.Equal(t, [][]Value{{1}, {2}, {3}}, values)
			},
		},
		{
			Name: "with an additional empty list",
			Lists: []List{
				{elV(day0, 1)},
				{},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Len(t, table.Times(), 0)
				assert.Equal(t, 0, table.NumberOfRows())
				assert.Equal(t, 2, table.NumberOfColumns())
			},
		},
		{
			Name: "with three shuffled matched elements per list",
			Lists: []List{
				{elV(day2, 1), elV(day1, 10), elV(day0, 100)},
				{elV(day1, 20), elV(day2, 2), elV(day0, 200)},
				{elV(day2, 3), elV(day0, 300), elV(day1, 30)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Len(t, table.Times(), 3)
				assert.Equal(t, 3, table.NumberOfRows())
				assert.Equal(t, 3, table.NumberOfColumns())
				assert.Equal(t, 9, table.NumberOfCells())
				values := table.UnderlyingValues()
				assert.Equal(t, [][]Value{
					{100, 10, 1},
					{200, 20, 2},
					{300, 30, 3},
				}, values)
				assert.Equal(t, []time.Time{
					date(day0), date(day1), date(day2),
				}, table.UnderlyingTimes())

				if row, ok := table.Row(date(day0)); assert.True(t, ok) {
					assert.Equal(t, []Value{100, 200, 300}, row)
				}
				if row, ok := table.Row(date(day1)); assert.True(t, ok) {
					assert.Equal(t, []Value{10, 20, 30}, row)
				}
				if row, ok := table.Row(date(day2)); assert.True(t, ok) {
					assert.Equal(t, []Value{1, 2, 3}, row)
				}
				if row, ok := table.Column(0); assert.True(t, ok) {
					assert.Equal(t, List{elV(day0, 100), elV(day1, 10), elV(day2, 1)}, row)
				}
				if row, ok := table.Column(1); assert.True(t, ok) {
					assert.Equal(t, List{elV(day0, 200), elV(day1, 20), elV(day2, 2)}, row)
				}
				if row, ok := table.Column(2); assert.True(t, ok) {
					assert.Equal(t, List{elV(day0, 300), elV(day1, 30), elV(day2, 3)}, row)
				}
			},
		},
		{
			Name: "with different dates",
			Lists: []List{
				{elV(day1, 1)},
				{elV(day2, 2)},
				{elV(day3, 3)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Equal(t, 0, table.NumberOfRows())
				assert.Equal(t, 3, table.NumberOfColumns())
				assert.Len(t, table.Times(), 0)
				assert.Equal(t, [][]Value{{}, {}, {}}, table.Values())
			},
		},
		{
			Name: "with a second column with more history",
			Lists: []List{
				{elV(day1, 2)},
				{elV(day0, 10), elV(day1, 20)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Equal(t, 1, table.NumberOfRows())
				assert.Equal(t, 2, table.NumberOfColumns())
				assert.Equal(t, []time.Time{date(day1)}, table.Times())
				assert.Equal(t, [][]Value{
					{2},
					{20},
				}, table.Values())
			},
		},
		{
			Name: "with a second column with less history",
			Lists: []List{
				{elT(day0), elV(day1, 2)},
				{elV(day1, 3)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Len(t, table.Times(), 1)
				assert.Equal(t, 1, table.NumberOfRows())
				assert.Equal(t, 2, table.NumberOfColumns())
				assert.Equal(t, []time.Time{date(day1)}, table.Times())
				assert.Equal(t, [][]Value{{2}, {3}}, table.Values())
			},
		},
		{
			Name: "with no overlap",
			Lists: []List{
				{elT(day0), elV(day2, 2)},
				{elV(day1, 3)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Len(t, table.Times(), 0)
				assert.Equal(t, 0, table.NumberOfRows())
				assert.Equal(t, 2, table.NumberOfColumns())
				assert.Equal(t, [][]Value{{}, {}}, table.Values())
			},
		},
		{
			Name: "with second asset having more data than the first overlap",
			Lists: []List{
				{elV(day0, 1), elV(day1, 2), elV(day3, 4)},
				{elV(day0, 10), elV(day1, 20), elV(day2, 30), elV(day3, 40)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Equal(t, 4, table.NumberOfRows())
				assert.Equal(t, 2, table.NumberOfColumns())
				assert.Equal(t, [][]Value{
					{1, 2, 0, 4},
					{10, 20, 30, 40},
				}, table.Values())
			},
		},
		{
			Name: "with second asset having less data than the first overlap",
			Lists: []List{
				{elV(day0, 10), elV(day1, 20), elV(day2, 30), elV(day3, 40)},
				{elV(day0, 1), elV(day1, 2), elV(day3, 4)},
			},
			Then: func(t *testing.T, table *Table) {
				assert.Equal(t, 4, table.NumberOfRows())
				assert.Equal(t, 2, table.NumberOfColumns())
				assert.Equal(t, [][]Value{
					{10, 20, 30, 40},
					{1, 2, 0, 4},
				}, table.Values())
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			result := timetable.New(tt.Lists...)
			tt.Then(t, result)
		})
	}
}

func TestCompact_JSON(t *testing.T) {
	emptyTable := new(Table)
	emptyTable = emptyTable.AddColumnFillMissingWithZero(List{})

	peach := timetable.New(List{
		elV(day0, 1), elV(day1, 2), elV(day2, 3),
	}, List{
		elV(day0, 10), elV(day1, 20), elV(day2, 30),
	})

	t.Run("marshal", func(t *testing.T) {
		for _, tt := range []struct {
			Name     string
			Table    *Table
			Expected string
		}{
			{Name: "nil", Table: nil, Expected: `null`},
			{Name: "zero", Table: new(Table), Expected: `{"times":null,"values":null}`},
			{Name: "empty", Table: emptyTable, Expected: `{"times":[],"values":[[]]}`},
			{Name: "some values", Table: peach,
				// language=json
				Expected: `{
					"times":[
						"2022-10-20T00:00:00Z",
						"2022-10-21T00:00:00Z",
						"2022-10-24T00:00:00Z"
					],
					"values":[
						[1,2,3],
						[10,20,30]
					]
				}`,
			},
		} {
			t.Run(tt.Name, func(t *testing.T) {
				data, err := json.Marshal(tt.Table)
				assert.NoError(t, err)
				assert.JSONEq(t, tt.Expected, string(data))
			})
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		for _, tt := range []struct {
			Name     string
			JSON     string
			Expected *Table
		}{
			{Name: "nil", Expected: nil, JSON: `null`},
			{Name: "zero", Expected: new(Table), JSON: `{"times":null,"values":null}`},
			{Name: "empty", Expected: emptyTable, JSON: `{"times":[],"values":[[]]}`},
			{Name: "reverse order", Expected: peach,
				// language=json
				JSON: `{
					"times":[
						"2022-10-24T00:00:00Z",
						"2022-10-21T00:00:00Z",
						"2022-10-20T00:00:00Z"
					],
					"values":[
						[3,2,1],
						[30,20,10]
					]
				}`,
			},
			{Name: "in order", Expected: peach,
				// language=json
				JSON: `{
					"times":[
						"2022-10-20T00:00:00Z",
						"2022-10-21T00:00:00Z",
						"2022-10-24T00:00:00Z"
					],
					"values":[
						[1,2,3],
						[10,20,30]
					]
				}`,
			},
		} {
			t.Run(tt.Name, func(t *testing.T) {
				var got *Table
				err := json.Unmarshal([]byte(tt.JSON), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, got)
			})
		}
	})
}

func TestCompact_Row(t *testing.T) {
	t.Run("row not found", func(t *testing.T) {
		table := timetable.New(List{
			elV(day0, 1),
			elV(day1, 2),
		})
		row, ok := table.Row(date(day2))
		assert.False(t, ok)
		assert.Len(t, row, 0)
	})
}

func TestCompact_Column(t *testing.T) {
	t.Run("column out of range", func(t *testing.T) {
		table := timetable.New(List{
			elV(day0, 1),
			elV(day1, 2),
		})
		column, ok := table.Column(2)
		assert.False(t, ok)
		assert.Len(t, column, 0)
	})

	t.Run("negative range", func(t *testing.T) {
		table := timetable.New(List{
			elV(day0, 1),
			elV(day1, 2),
		})
		column, ok := table.Column(-1)
		assert.False(t, ok)
		assert.Len(t, column, 0)
	})
}
