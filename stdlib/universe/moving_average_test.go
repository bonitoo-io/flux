package universe_test

import (
	"fmt"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/execute/executetest"
	"github.com/influxdata/flux/querytest"
	"github.com/influxdata/flux/stdlib/universe"
	"testing"
	"time"
)

func TestMovingAverageOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"movingAverage","kind":"movingAverage","spec":{"n":1}}`)
	op := &flux.Operation{
		ID: "movingAverage",
		Spec: &universe.MovingAverageOpSpec{
			N: 1,
		},
	}
	querytest.OperationMarshalingTestHelper(t, data, op)
}

func TestMovingAverage_PassThrough(t *testing.T) {
	executetest.TransformationPassThroughTestHelper(t, func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
		s := universe.NewMovingAverageTransformation(
			d,
			c,
			&universe.MovingAverageProcedureSpec{},
		)
		return s
	})
}

func TestMovingAverage_Process(t *testing.T) {
	testCases := []struct {
		name    string
		spec    *universe.MovingAverageProcedureSpec
		data    []flux.Table
		want    []*executetest.Table
		wantErr error
	}{
		{
			name: "float",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       2,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 2.0},
					{execute.Time(2), 4.0},
					{execute.Time(3), 5.0},
					{execute.Time(4), 9.0},
					{execute.Time(5), 8.0},
					{execute.Time(6), 11.0},
					{execute.Time(7), 15.0},
					{execute.Time(8), 12.0},
					{execute.Time(9), 5.0},
					{execute.Time(10), 7.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(2), 3.0},
					{execute.Time(3), 4.5},
					{execute.Time(4), 7.0},
					{execute.Time(5), 8.5},
					{execute.Time(6), 9.5},
					{execute.Time(7), 13.0},
					{execute.Time(8), 13.5},
					{execute.Time(9), 8.5},
					{execute.Time(10), 6.0},
				},
			}},
		},
		{
			name: "float with 3",
			spec: &universe.MovingAverageProcedureSpec{
				Columns:    []string{execute.DefaultValueColLabel},
				TimeColumn: execute.DefaultTimeColLabel,
				N:          3,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1 * time.Second), 2.0},
					{execute.Time(2 * time.Second), 4.0},
					{execute.Time(3 * time.Second), 5.0},
					{execute.Time(4 * time.Second), 9.0},
					{execute.Time(5 * time.Second), 8.0},
					{execute.Time(6 * time.Second), 11.0},
					{execute.Time(7 * time.Second), 15.0},
					{execute.Time(8 * time.Second), 12.0},
					{execute.Time(9 * time.Second), 5.0},
					{execute.Time(10 * time.Second), 7.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(3 * time.Second), 11.0 / 3},
					{execute.Time(4 * time.Second), 6.0},
					{execute.Time(5 * time.Second), 22.0 / 3},
					{execute.Time(6 * time.Second), 28.0 / 3},
					{execute.Time(7 * time.Second), 34.0 / 3},
					{execute.Time(8 * time.Second), 38.0 / 3},
					{execute.Time(9 * time.Second), 32.0 / 3},
					{execute.Time(10 * time.Second), 8.0},
				},
			}},
		},
		{
			name: "int",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       2,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(2)},
					{execute.Time(2), int64(4)},
					{execute.Time(3), int64(5)},
					{execute.Time(4), int64(9)},
					{execute.Time(5), int64(8)},
					{execute.Time(6), int64(11)},
					{execute.Time(7), int64(15)},
					{execute.Time(8), int64(12)},
					{execute.Time(9), int64(5)},
					{execute.Time(10), int64(7)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(2), 3.0},
					{execute.Time(3), 4.5},
					{execute.Time(4), 7.0},
					{execute.Time(5), 8.5},
					{execute.Time(6), 9.5},
					{execute.Time(7), 13.0},
					{execute.Time(8), 13.5},
					{execute.Time(9), 8.5},
					{execute.Time(10), 6.0},
				},
			}},
		},
		{
			name: "int with 3",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       3,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1 * time.Second), int64(2)},
					{execute.Time(2 * time.Second), int64(4)},
					{execute.Time(3 * time.Second), int64(5)},
					{execute.Time(4 * time.Second), int64(9)},
					{execute.Time(5 * time.Second), int64(8)},
					{execute.Time(6 * time.Second), int64(11)},
					{execute.Time(7 * time.Second), int64(15)},
					{execute.Time(8 * time.Second), int64(12)},
					{execute.Time(9 * time.Second), int64(5)},
					{execute.Time(10 * time.Second), int64(7)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(3 * time.Second), 11.0 / 3},
					{execute.Time(4 * time.Second), 6.0},
					{execute.Time(5 * time.Second), 22.0 / 3},
					{execute.Time(6 * time.Second), 28.0 / 3},
					{execute.Time(7 * time.Second), 34.0 / 3},
					{execute.Time(8 * time.Second), 38.0 / 3},
					{execute.Time(9 * time.Second), 32.0 / 3},
					{execute.Time(10 * time.Second), 8.0},
				},
			}},
		},
		{
			name: "uint",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       2,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(2)},
					{execute.Time(2), uint64(4)},
					{execute.Time(3), uint64(5)},
					{execute.Time(4), uint64(9)},
					{execute.Time(5), uint64(8)},
					{execute.Time(6), uint64(11)},
					{execute.Time(7), uint64(15)},
					{execute.Time(8), uint64(12)},
					{execute.Time(9), uint64(5)},
					{execute.Time(10), uint64(7)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(2), 3.0},
					{execute.Time(3), 4.5},
					{execute.Time(4), 7.0},
					{execute.Time(5), 8.5},
					{execute.Time(6), 9.5},
					{execute.Time(7), 13.0},
					{execute.Time(8), 13.5},
					{execute.Time(9), 8.5},
					{execute.Time(10), 6.0},
				},
			}},
		},
		{
			name: "uint with 3",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       3,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1 * time.Second), uint64(2)},
					{execute.Time(2 * time.Second), uint64(4)},
					{execute.Time(3 * time.Second), uint64(5)},
					{execute.Time(4 * time.Second), uint64(9)},
					{execute.Time(5 * time.Second), uint64(8)},
					{execute.Time(6 * time.Second), uint64(11)},
					{execute.Time(7 * time.Second), uint64(15)},
					{execute.Time(8 * time.Second), uint64(12)},
					{execute.Time(9 * time.Second), uint64(5)},
					{execute.Time(10 * time.Second), uint64(7)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(3 * time.Second), 11.0 / 3},
					{execute.Time(4 * time.Second), 6.0},
					{execute.Time(5 * time.Second), 22.0 / 3},
					{execute.Time(6 * time.Second), 28.0 / 3},
					{execute.Time(7 * time.Second), 34.0 / 3},
					{execute.Time(8 * time.Second), 38.0 / 3},
					{execute.Time(9 * time.Second), 32.0 / 3},
					{execute.Time(10 * time.Second), 8.0},
				},
			}},
		},
		{
			name: "int",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       2,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(2)},
					{execute.Time(2), int64(4)},
					{execute.Time(3), int64(5)},
					{execute.Time(4), int64(9)},
					{execute.Time(5), int64(8)},
					{execute.Time(6), int64(11)},
					{execute.Time(7), int64(15)},
					{execute.Time(8), int64(12)},
					{execute.Time(9), int64(5)},
					{execute.Time(10), int64(7)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(2), 3.0},
					{execute.Time(3), 4.5},
					{execute.Time(4), 7.0},
					{execute.Time(5), 8.5},
					{execute.Time(6), 9.5},
					{execute.Time(7), 13.0},
					{execute.Time(8), 13.5},
					{execute.Time(9), 8.5},
					{execute.Time(10), 6.0},
				},
			}},
		},
		{
			name: "int with 3",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       3,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1 * time.Second), int64(2)},
					{execute.Time(2 * time.Second), int64(4)},
					{execute.Time(3 * time.Second), int64(5)},
					{execute.Time(4 * time.Second), int64(9)},
					{execute.Time(5 * time.Second), int64(8)},
					{execute.Time(6 * time.Second), int64(11)},
					{execute.Time(7 * time.Second), int64(15)},
					{execute.Time(8 * time.Second), int64(12)},
					{execute.Time(9 * time.Second), int64(5)},
					{execute.Time(10 * time.Second), int64(7)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(3 * time.Second), 11.0 / 3},
					{execute.Time(4 * time.Second), 6.0},
					{execute.Time(5 * time.Second), 22.0 / 3},
					{execute.Time(6 * time.Second), 28.0 / 3},
					{execute.Time(7 * time.Second), 34.0 / 3},
					{execute.Time(8 * time.Second), 38.0 / 3},
					{execute.Time(9 * time.Second), 32.0 / 3},
					{execute.Time(10 * time.Second), 8.0},
				},
			}},
		},
		{
			name: "float with tags",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{execute.DefaultValueColLabel},
				N:       3,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(1 * time.Second), 2.0, "a"},
					{execute.Time(2 * time.Second), 4.0, "b"},
					{execute.Time(3 * time.Second), 5.0, "c"},
					{execute.Time(4 * time.Second), 9.0, "d"},
					{execute.Time(5 * time.Second), 8.0, "e"},
					{execute.Time(6 * time.Second), 11.0, "f"},
					{execute.Time(7 * time.Second), 15.0, "g"},
					{execute.Time(8 * time.Second), 12.0, "h"},
					{execute.Time(9 * time.Second), 5.0, "i"},
					{execute.Time(10 * time.Second), 7.0, "j"},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(3 * time.Second), 11.0 / 3, "c"},
					{execute.Time(4 * time.Second), 6.0, "d"},
					{execute.Time(5 * time.Second), 22.0 / 3, "e"},
					{execute.Time(6 * time.Second), 28.0 / 3, "f"},
					{execute.Time(7 * time.Second), 34.0 / 3, "g"},
					{execute.Time(8 * time.Second), 38.0 / 3, "h"},
					{execute.Time(9 * time.Second), 32.0 / 3, "i"},
					{execute.Time(10 * time.Second), 8.0, "j"},
				},
			}},
		},
		{
			name: "float with null values",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{"x", "y"},
				N:       1,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "x", Type: flux.TFloat},
					{Label: "y", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 2.0, nil},
					{execute.Time(2), nil, 10.0},
					{execute.Time(3), 8.0, 20.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "x", Type: flux.TFloat},
					{Label: "y", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(2), 2.0, nil},
					{execute.Time(3), 8.0, 15.0},
				},
			}},
		},
		{
			name: "nulls in time column",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{"x", "y"},
				N:       1,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "x", Type: flux.TFloat},
					{Label: "y", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{nil, 2.0, 3.0},
					{execute.Time(2), nil, 10.0},
					{nil, 8.0, 20.0},
					{execute.Time(4), 8.0, 20.0},
					{nil, 8.0, 20.0},
					{execute.Time(6), 10.0, 25.0},
					{nil, 8.0, 20.0},
				},
			}},
			wantErr: fmt.Errorf("moving average found null time in time column"),
		},
		{
			name: "times out of order",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{"x", "y"},
				N:       1,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "x", Type: flux.TFloat},
					{Label: "y", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(2), nil, 10.0},
					{execute.Time(4), 8.0, 20.0},
					{execute.Time(6), 10.0, 25.0},

					{execute.Time(3), nil, 10.0},
					{execute.Time(5), 8.0, 20.0},
					{execute.Time(7), 10.0, nil},
				},
			}},
			wantErr: fmt.Errorf("moving average found out-of-order times in time column"),
		},
		{
			name: "pass through with repeated times",
			spec: &universe.MovingAverageProcedureSpec{
				Columns: []string{"x"},
				N:       2,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "x", Type: flux.TFloat},
					{Label: "b", Type: flux.TBool},
					{Label: "s", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(2), 6.0, false, "bar"},
					{execute.Time(2), 1.0, false, "bar"},
					{execute.Time(4), 8.0, false, nil},
					{execute.Time(4), 9.0, true, "baz"},
					{execute.Time(6), 10.0, nil, "dog"},
					{execute.Time(6), 10.0, nil, "dog"},
					{execute.Time(6), 10.0, nil, "dog"},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "x", Type: flux.TFloat},
					{Label: "b", Type: flux.TBool},
					{Label: "s", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(4), 7.0, false, nil},
					{execute.Time(6), 9.0, nil, "dog"},
				},
			}},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.ProcessTestHelper(
				t,
				tc.data,
				tc.want,
				tc.wantErr,
				func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
					return universe.NewMovingAverageTransformation(d, c, tc.spec)
				},
			)
		})
	}
}
