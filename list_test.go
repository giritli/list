package list_test

import (
	"fmt"
	"github.com/giritli/list"
	"reflect"
	"strconv"
	"testing"
)

func TestConvert(t *testing.T) {
	type testCase struct {
		name   string
		before list.Of[int]
		fn     func(int) string
		after  list.Of[string]
	}

	testCases := []testCase{
		{
			"convert int list to string",
			list.Of[int]{1, 2, 3, 4},
			func(i int) string {
				return strconv.Itoa(i)
			},
			list.Of[string]{"1", "2", "3", "4"},
		},
		{
			"convert and multiply int list to string",
			list.Of[int]{1, 2, 3, 4},
			func(i int) string {
				return strconv.Itoa(i * i)
			},
			list.Of[string]{"1", "4", "9", "16"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := list.Convert(tt.before, tt.fn)
			if !reflect.DeepEqual(tt.after, got) {
				t.Errorf("wanted %v, got %v", tt.after, got)
			}
		})
	}
}

func TestReduceInto(t *testing.T) {
	type testCase struct {
		name string
		fn   func() any
		want any
	}

	type testStruct struct {
		total int
	}

	testCases := []testCase{
		{
			"convert int to string",
			func() any {
				before := list.Of[int]{1, 2, 3, 4, 5}
				after := list.ReduceInto(before, func(to string, from int) string {
					return fmt.Sprintf("%s%d", to, from)
				})

				return after
			},
			"12345",
		},
		{
			"sum struct fields",
			func() any {
				before := list.Of[testStruct]{
					{1}, {7}, {2},
				}
				after := list.ReduceInto(before, func(to int, from testStruct) int {
					return to + from.total
				})

				return after
			},
			10,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("wanted %v, got %v", tt.want, got)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type testCase struct {
		name     string
		items    list.Of[int]
		item     int
		contains bool
	}

	testCases := []testCase{
		{
			"check number exists in list",
			list.Of[int]{1, 2, 3, 4, 5},
			1,
			true,
		},
		{
			"check number does not exist in list",
			list.Of[int]{1, 2, 3, 4, 5},
			6,
			false,
		},
		{
			"check number exists more than once",
			list.Of[int]{1, 2, 3, 4, 5, 4},
			4,
			true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.items.Contains(tt.item, func(a, b int) bool {
				return a == b
			})
			if tt.contains != got {
				t.Errorf("wanted %v, got %v", tt.contains, got)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	type testCase struct {
		name   string
		before list.Of[int]
		size   uint
		after  []list.Of[int]
	}

	testCases := []testCase{
		{
			"empty list",
			list.Of[int]{},
			0,
			[]list.Of[int]{{}},
		},
		{
			"max list",
			list.Of[int]{1, 2, 3, 4, 5},
			0,
			[]list.Of[int]{{1, 2, 3, 4, 5}},
		},
		{
			"size of 1",
			list.Of[int]{1, 2, 3, 4, 5},
			1,
			[]list.Of[int]{{1}, {2}, {3}, {4}, {5}},
		},
		{
			"size of 2",
			list.Of[int]{1, 2, 3, 4, 5},
			2,
			[]list.Of[int]{{1, 2}, {3, 4}, {5}},
		},
		{
			"no remainders",
			list.Of[int]{1, 2, 3, 4, 5, 6},
			2,
			[]list.Of[int]{{1, 2}, {3, 4}, {5, 6}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.before.Chunk(tt.size)
			if !reflect.DeepEqual(tt.after, got) {
				t.Errorf("wanted %v, got %v", tt.after, got)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type testCase struct {
		name   string
		before list.Of[int]
		fn     func(int) bool
		after  list.Of[int]
	}

	testCases := []testCase{
		{
			"true",
			list.Of[int]{1, 2, 3, 4, 5},
			func(i int) bool {
				return true
			},
			list.Of[int]{1, 2, 3, 4, 5},
		},
		{
			"false",
			list.Of[int]{1, 2, 3, 4, 5},
			func(i int) bool {
				return false
			},
			list.Of[int]{},
		},
		{
			"even",
			list.Of[int]{1, 2, 3, 4, 5},
			func(i int) bool {
				return i%2 == 0
			},
			list.Of[int]{2, 4},
		},
		{
			"greater than 3",
			list.Of[int]{1, 2, 3, 4, 5},
			func(i int) bool {
				return i > 3
			},
			list.Of[int]{4, 5},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.before.Filter(tt.fn)
			if !reflect.DeepEqual(tt.after, got) {
				t.Errorf("wanted %v, got %v", tt.after, got)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	type testCase struct {
		name   string
		before list.Of[int]
		at     uint
		items  []int
		after  list.Of[int]
	}

	testCases := []testCase{
		{
			"no insert",
			list.Of[int]{1, 2, 3, 4, 5},
			0,
			[]int{},
			list.Of[int]{1, 2, 3, 4, 5},
		},
		{
			"insert before",
			list.Of[int]{1, 2, 3, 4, 5},
			0,
			[]int{9, 8},
			list.Of[int]{9, 8, 1, 2, 3, 4, 5},
		},
		{
			"insert after",
			list.Of[int]{1, 2, 3, 4, 5},
			5,
			[]int{9, 8},
			list.Of[int]{1, 2, 3, 4, 5, 9, 8},
		},
		{
			"insert after beyond bounds",
			list.Of[int]{1, 2, 3, 4, 5},
			9,
			[]int{9, 8},
			list.Of[int]{1, 2, 3, 4, 5, 9, 8},
		},
		{
			"insert middle",
			list.Of[int]{1, 2, 3, 4, 5},
			2,
			[]int{9, 8},
			list.Of[int]{1, 2, 9, 8, 3, 4, 5},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.before.Insert(tt.at, tt.items...)
			if !reflect.DeepEqual(tt.after, got) {
				t.Errorf("wanted %v, got %v", tt.after, got)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type testCase struct {
		name   string
		before list.Of[int]
		fn     func(int) int
		after  list.Of[int]
	}

	testCases := []testCase{
		{
			"numbers multiply",
			list.Of[int]{1, 2, 3, 4, 5},
			func(i int) int {
				return i * 2
			},
			list.Of[int]{2, 4, 6, 8, 10},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.before.Map(tt.fn)
			if !reflect.DeepEqual(tt.after, got) {
				t.Errorf("wanted %v, got %v", tt.after, got)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type testCase struct {
		name   string
		before list.Of[int]
		fn     func(a, b int) bool
		after  list.Of[int]
	}

	testCases := []testCase{
		{
			"numbers asc",
			list.Of[int]{9, 1, 5, 7, 3, 2},
			func(a, b int) bool {
				return a < b
			},
			list.Of[int]{1, 2, 3, 5, 7, 9},
		},
		{
			"numbers desc",
			list.Of[int]{9, 1, 5, 7, 3, 2},
			func(a, b int) bool {
				return a > b
			},
			list.Of[int]{9, 7, 5, 3, 2, 1},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.before.Sort(tt.fn)
			if !reflect.DeepEqual(tt.after, got) {
				t.Errorf("wanted %v, got %v", tt.after, got)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type testCase struct {
		name   string
		before list.Of[int]
		after  list.Of[int]
	}

	testCases := []testCase{
		{
			"all unique",
			list.Of[int]{1, 2, 3, 4, 5},
			list.Of[int]{1, 2, 3, 4, 5},
		},
		{
			"some unique",
			list.Of[int]{1, 1, 2, 3, 4, 5, 4},
			list.Of[int]{1, 2, 3, 4, 5},
		},
		{
			"some unique out of order",
			list.Of[int]{1, 1, 2, 3, 6, 5, 4, 4},
			list.Of[int]{1, 2, 3, 6, 5, 4},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.before.Unique()
			if !reflect.DeepEqual(tt.after, got) {
				t.Errorf("wanted %v, got %v", tt.after, got)
			}
		})
	}
}
