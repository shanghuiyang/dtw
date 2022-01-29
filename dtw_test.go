package dtw

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type point struct {
	x float64
	y float64
}

func TestDTW(t *testing.T) {
	testCases := []struct {
		desc    string
		s       interface{}
		t       interface{}
		df      DistanceFunc
		dist    float64
		path    [][2]int
		noError bool
	}{
		{
			desc: "int series",
			s:    []int{1, 3, 4, 9, 8, 2, 1, 5, 7, 3},
			t:    []int{1, 6, 2, 3, 0, 9, 4, 3, 6, 3},
			df: func(x, y interface{}) float64 {
				xx := x.(int)
				yy := y.(int)
				return math.Abs(float64(xx - yy))
			},
			dist: 15,
			path: [][2]int{
				{0, 0}, {1, 1}, {1, 2}, {1, 3}, {2, 4}, {3, 5}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 8}, {9, 9},
			},
			noError: true,
		},
		{
			desc: "float series",
			s:    []float64{1, 3, 4, 9, 8, 2, 1, 5, 7, 3},
			t:    []float64{1, 6, 2, 3, 0, 9, 4, 3, 6, 3},
			df: func(x, y interface{}) float64 {
				xx := x.(float64)
				yy := y.(float64)
				return math.Abs(xx - yy)
			},
			dist: 15,
			path: [][2]int{
				{0, 0}, {1, 1}, {1, 2}, {1, 3}, {2, 4}, {3, 5}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 8}, {9, 9},
			},
			noError: true,
		},
		{
			desc: "2D series",
			s:    []point{{0, 0}, {1, 0}, {2, 0}},
			t:    []point{{0, 1}, {1, 1}, {2, 1}},
			df: func(x, y interface{}) float64 {
				p1 := x.(point)
				p2 := y.(point)
				return math.Sqrt((p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y))
			},
			dist: 3,
			path: [][2]int{
				{0, 0}, {1, 1}, {2, 2},
			},
			noError: true,
		},
		{
			desc: "same series",
			s:    []int{1, 2, 3},
			t:    []int{1, 2, 3},
			df: func(x, y interface{}) float64 {
				xx := x.(int)
				yy := y.(int)
				return math.Abs(float64(xx - yy))
			},
			dist: 0,
			path: [][2]int{
				{0, 0}, {1, 1}, {2, 2},
			},
			noError: true,
		},
		{
			desc: "series with different length",
			s:    []int{1, 2, 3},
			t:    []int{1, 2, 3, 4},
			df: func(x, y interface{}) float64 {
				xx := x.(int)
				yy := y.(int)
				return math.Abs(float64(xx - yy))
			},
			dist: 1,
			path: [][2]int{
				{0, 0}, {1, 1}, {2, 2}, {3, 2},
			},
			noError: true,
		},
		{
			desc: "empty series",
			s:    []int{},
			t:    []int{},
			df: func(x, y interface{}) float64 {
				return 0
			},
			dist:    0,
			path:    [][2]int{},
			noError: false,
		},
		{
			desc:    "invalid distance func",
			s:       []int{1, 2},
			t:       []int{1, 3},
			df:      nil,
			dist:    0,
			path:    [][2]int{},
			noError: false,
		},
	}

	for _, test := range testCases {
		dtw := New()
		d, err := dtw.Distance(test.s, test.t, test.df)
		if test.noError {
			assert.NoError(t, err)
			assert.InDelta(t, test.dist, d, 0.01)
			path := dtw.Path()
			assert.Equal(t, test.path, path)
			assert.NotPanics(t, func() { dtw.Draw(os.Stdout) })
			continue
		}
		assert.Error(t, err)
	}
}
