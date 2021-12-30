package dtw

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// DistanceFunc is the function used to calculate the distance between x and y
type DistanceFunc func(x, y interface{}) float64

// Dtw ...
type Dtw struct {
	s      reflect.Value
	t      reflect.Value
	df     DistanceFunc
	matrix [][]float64
}

// New ...
func New() *Dtw {
	return &Dtw{}
}

// Distance calculates the DTW between series s and t.
func (dtw *Dtw) Distance(s, t interface{}, f DistanceFunc) (float64, error) {
	if f == nil {
		return 0, errors.New("invalid distance func")
	}
	dtw.df = f

	if reflect.TypeOf(s).Kind() != reflect.Slice {
		return 0, errors.New("series s is not a slice")
	}
	if reflect.TypeOf(t).Kind() != reflect.Slice {
		return 0, errors.New("series t is not a slice")
	}

	ss := reflect.ValueOf(s)
	tt := reflect.ValueOf(t)

	if ss.Len() == 0 {
		return 0, errors.New("s series is empty")
	}

	if tt.Len() == 0 {
		return 0, errors.New("t series is empty")
	}

	// get dtw.s be the longer series
	dtw.s, dtw.t = ss, tt
	if ss.Len() < tt.Len() {
		dtw.s, dtw.t = tt, ss
	}

	dtw.initMatrix()
	m, n := dtw.s.Len(), dtw.t.Len()
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			dtw.matrix[r][c] = dtw.dist(r, c)
		}
	}

	return dtw.matrix[m-1][n-1], nil
}

// Path ...
func (dtw *Dtw) Path() [][2]int {
	m, n := dtw.s.Len(), dtw.t.Len()
	r, c := m-1, n-1

	path := [][2]int{}
	path = append(path, [2]int{r, c})
	for r > 0 && c > 0 {
		min := dtw.matrix[r-1][c-1]
		rr, cc := r-1, c-1
		if dtw.matrix[r-1][c] < min {
			rr, cc = r-1, c
			min = dtw.matrix[r-1][c]
		}
		if dtw.matrix[r][c-1] < min {
			rr, cc = r, c-1
		}
		r, c = rr, cc
		path = append(path, [2]int{r, c})
	}

	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// Draw ...
func (dtw *Dtw) Draw(w io.Writer) {
	m, n := dtw.s.Len(), dtw.t.Len()

	matrix := make([][]string, m)
	for i := 0; i < m; i++ {
		matrix[i] = make([]string, n)
	}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			s := fmt.Sprintf("%.2f ", dtw.matrix[r][c])
			matrix[r][c] = fmt.Sprintf("%15s", s)
		}
	}

	path := dtw.Path()
	for _, p := range path {
		r, c := p[0], p[1]
		s := fmt.Sprintf("[%.2f]", dtw.matrix[r][c])
		matrix[r][c] = fmt.Sprintf("%15s", s)
	}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			fmt.Fprintf(w, "%v", matrix[r][c])
		}
		fmt.Fprintln(w)
	}
}

func (dtw *Dtw) initMatrix() {
	m := dtw.s.Len()
	n := dtw.t.Len()
	matrix := make([][]float64, m)
	for i := 0; i < m; i++ {
		matrix[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = math.Inf(1)
		}
	}
	dtw.matrix = matrix
}

func (dtw *Dtw) min(v1, v2, v3 float64) float64 {
	min := v1
	if v2 < min {
		min = v2
	}
	if v3 < min {
		min = v3
	}
	return min
}

func (dtw *Dtw) dist(r, c int) float64 {
	dist := dtw.df(dtw.s.Index(r).Interface(), dtw.t.Index(c).Interface())
	if r == 0 && c == 0 {
		return dist
	}
	if r == 0 && c > 0 {
		return dist + dtw.matrix[r][c-1]
	}
	if c == 0 && r > 0 {
		return dist + dtw.matrix[r-1][c]
	}
	return dist + dtw.min(dtw.matrix[r-1][c-1], dtw.matrix[r-1][c], dtw.matrix[r][c-1])
}
