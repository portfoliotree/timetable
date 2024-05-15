package timetable_test

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/portfoliotree/timetable"
)

type List = timetable.List[Value]

func fourInARow() List {
	return List{elT(day0), elT(day1), elT(day2), elT(day3)}
}

func TestList_Between(t *testing.T) {
	for _, tt := range []struct {
		Name       string
		Start, End time.Time
		List       List
		Then       func(t *testing.T, result List)
	}{
		{
			Name:  "empty",
			List:  List{},
			Start: date(day0), End: date(day1),
			Then: func(t *testing.T, result List) {
				assert.Len(t, result, 0)
			},
		},
		{
			Name: "nil",
			List: List(nil),
			Then: func(t *testing.T, result List) {
				assert.Len(t, result, 0)
			},
		},
		{
			Name:  "out of bounds",
			List:  fourInARow(),
			Start: date(dayBefore), End: date(dayAfter),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, fourInARow(), result)
			},
		},
		{
			Name:  "list out of order",
			List:  List{elT(day0), elT(day2), elT(day3), elT(day1)},
			Start: date(dayBefore), End: date(dayAfter),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, fourInARow(), result)
			},
		},
		{
			Name:  "range out of order",
			List:  fourInARow(),
			Start: date(day3), End: date(day0),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, fourInARow(), result)
			},
		},
		{
			Name:  "both before",
			List:  fourInARow(),
			Start: date(dayBefore), End: date(dayBefore).AddDate(0, 0, -1),
			Then: func(t *testing.T, result List) {
				assert.Len(t, result, 0)
			},
		},
		{
			Name:  "both after",
			List:  fourInARow(),
			Start: date(dayAfter), End: date(dayAfter).AddDate(0, 0, 1),
			Then: func(t *testing.T, result List) {
				assert.Len(t, result, 0)
			},
		},
		{
			Name:  "same day",
			Start: date(day1), End: date(day1),
			List: fourInARow(),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day1)}, result)
			},
		},
		{
			Name:  "days between",
			List:  fourInARow(),
			Start: date(day1), End: date(day3),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day1), elT(day2), elT(day3)}, result)
			},
		},
		{
			Name:  "one element",
			List:  List{elT(day0)},
			Start: date(day0), End: date(day0),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day0)}, result)
			},
		},
		{
			Name:  "two elements",
			List:  List{elT(day0), elT(day1)},
			Start: date(day0), End: date(day1),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day0), elT(day1)}, result)
			},
		},
		{
			Name:  "truncate the first elements exact match",
			List:  List{elT(day0), elT(day1), elT(day2)},
			Start: date(day1), End: date(day2),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day1), elT(day2)}, result)
			},
		},
		{
			Name:  "truncate the first elements index not found",
			Start: date(day1), End: date(day3),
			List: List{elT(day0), elT(day2), elT(day3)},
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day2), elT(day3)}, result)
			},
		},
		{
			Name:  "truncate the last elements exact match",
			List:  List{elT(day0), elT(day1), elT(day2)},
			Start: date(day1), End: date(day2),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day1), elT(day2)}, result)
			},
		},
		{
			Name:  "truncate the last elements index not found",
			List:  List{elT(day0), elT(day1), elT(day3)},
			Start: date(day0), End: date(day2),
			Then: func(t *testing.T, result List) {
				assert.Equal(t, List{elT(day0), elT(day1)}, result)
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			clone := slices.Clone(tt.List)

			result := tt.List.Between(tt.Start, tt.End)
			tt.Then(t, result)

			t.Run("table and list have same behavior", func(t *testing.T) {
				table := timetable.New(clone)
				updated := table.Between(tt.Start, tt.End)
				row, ok := updated.Column(0)
				assert.True(t, ok)
				tt.Then(t, row)
			})
		})
	}

	t.Run("end during weekend", func(t *testing.T) {
		t.Run("on sunday", func(t *testing.T) {
			twoWeeksOfReturns := List{
				elT("2021-04-23"),
				elT("2021-04-22"),
				elT("2021-04-21"),
				elT("2021-04-20"),
				elT("2021-04-19"),
				elT("2021-04-16"),
				elT("2021-04-15"),
				elT("2021-04-14"),
				elT("2021-04-13"),
				elT("2021-04-12"),
			}

			end := date("2021-04-18")
			assert.Equal(t, end.Weekday(), time.Sunday)
			start := date("2021-04-12")

			result := twoWeeksOfReturns.Between(end, start)

			assert.Equal(t, List{
				elT("2021-04-12"),
				elT("2021-04-13"),
				elT("2021-04-14"),
				elT("2021-04-15"),
				elT("2021-04-16"),
			}, result)
		})

		t.Run("on saturday", func(t *testing.T) {
			twoWeeksOfReturns := List{
				elT("2021-04-23"),
				elT("2021-04-22"),
				elT("2021-04-21"),
				elT("2021-04-20"),
				elT("2021-04-19"),
				elT("2021-04-16"),
				elT("2021-04-15"),
				elT("2021-04-14"),
				elT("2021-04-13"),
				elT("2021-04-12"),
			}

			end := date("2021-04-17")
			assert.Equal(t, end.Weekday(), time.Saturday)
			start := date("2021-04-12")

			result := twoWeeksOfReturns.Between(end, start)

			assert.Equal(t, List{
				elT("2021-04-12"),
				elT("2021-04-13"),
				elT("2021-04-14"),
				elT("2021-04-15"),
				elT("2021-04-16"),
			}, result)
		})
	})
}
