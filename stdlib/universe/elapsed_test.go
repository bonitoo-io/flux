package universe_test

import (
	"github.com/influxdata/flux/querytest"
	"testing"
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/execute/executetest"
	"github.com/influxdata/flux/stdlib/universe"
)

func TestElapsedOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"elapsed","kind":"elapsed","spec":{"timeColumn": "_time"}}`)
	op := &flux.Operation{
		ID: "elapsed",
		Spec: &universe.ElapsedOpSpec{
			TimeColumn: "_time",
		},
	}
	querytest.OperationMarshalingTestHelper(t, data, op)
}

func TestElapsed_PassThrough(t *testing.T) {
	executetest.TransformationPassThroughTestHelper(t, func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
		s := universe.NewElapsedTransformation(
			d,
			c,
			&universe.ElapsedProcedureSpec{},
		)
		return s
	})
}


func TestElapsed_Process(t *testing.T) {
	testCases := []struct {
		name string
		spec *universe.ElapsedProcedureSpec
		data []flux.Table
		want []*executetest.Table
	}{
		{
			name: "basic",
			spec: &universe.ElapsedProcedureSpec{
				Unit:       flux.Duration(time.Nanosecond),
				TimeColumn: execute.DefaultTimeColLabel,
				ColumnName: "elapsed",
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
				},
				Data: [][]interface{}{
					{execute.Time(1)},
					{execute.Time(2)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "elapsed", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(2), int64(execute.Time(2) - execute.Time(1))},
				},
			}},
		},
		{
			name: "a little less basic, but still simple",
			spec: &universe.ElapsedProcedureSpec{
				Unit:       flux.Duration(time.Second),
				TimeColumn: execute.DefaultTimeColLabel,
				ColumnName: "elapsed",
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
				},
				Data: [][]interface{}{
					{execute.Time(1)},
					{execute.Time(2)},
					{execute.Time(3)},
					{execute.Time(4)},
					{execute.Time(5)},
					{execute.Time(6)},
					{execute.Time(7)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "elapsed", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(2), int64(execute.Time(2) - execute.Time(1))},
					{execute.Time(3), int64(execute.Time(3) - execute.Time(2))},
					{execute.Time(4), int64(execute.Time(4) - execute.Time(3))},
					{execute.Time(5), int64(execute.Time(5) - execute.Time(4))},
					{execute.Time(6), int64(execute.Time(6) - execute.Time(5))},
					{execute.Time(7), int64(execute.Time(7) - execute.Time(6))},
				},
			}},
		},
		{
			name: "two columns: time, _value",
			spec: &universe.ElapsedProcedureSpec{
				Unit:       flux.Duration(time.Nanosecond),
				TimeColumn: execute.DefaultTimeColLabel,
				ColumnName: "elapsed",
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(2)},
					{execute.Time(2), int64(2)},
					{execute.Time(3), int64(2)},
					{execute.Time(4), int64(2)},
					{execute.Time(5), int64(2)},
					{execute.Time(6), int64(2)},
					{execute.Time(7), int64(2)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
					{Label: "elapsed", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(2), int64(2), int64(execute.Time(2) - execute.Time(1))},
					{execute.Time(3), int64(2), int64(execute.Time(3) - execute.Time(2))},
					{execute.Time(4), int64(2), int64(execute.Time(4) - execute.Time(3))},
					{execute.Time(5), int64(2), int64(execute.Time(5) - execute.Time(4))},
					{execute.Time(6), int64(2), int64(execute.Time(6) - execute.Time(5))},
					{execute.Time(7), int64(2), int64(execute.Time(7) - execute.Time(6))},
				},
			}},
		},
		{
			name: "three columns: time, _value, path",
			spec: &universe.ElapsedProcedureSpec{
				Unit:       flux.Duration(time.Nanosecond),
				TimeColumn: execute.DefaultTimeColLabel,
				ColumnName: "elapsed",
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "path", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(1), 2.0, "/"},
					{execute.Time(2), 1.0, "/"},
					{execute.Time(3), 3.6, "/"},
					{execute.Time(4), 9.7, "/"},
					{execute.Time(5), 13.1, "/"},
					{execute.Time(6), 10.2, "/"},
					{execute.Time(7), 5.4, "/"},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "path", Type: flux.TString},
					{Label: "elapsed", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(2), 1.0, "/", int64(execute.Time(2) - execute.Time(1))},
					{execute.Time(3), 3.6, "/", int64(execute.Time(3) - execute.Time(2))},
					{execute.Time(4), 9.7, "/", int64(execute.Time(4) - execute.Time(3))},
					{execute.Time(5), 13.1, "/", int64(execute.Time(5) - execute.Time(4))},
					{execute.Time(6), 10.2, "/", int64(execute.Time(6) - execute.Time(5))},
					{execute.Time(7), 5.4, "/", int64(execute.Time(7) - execute.Time(6))},
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
				nil,
				func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
					return universe.NewElapsedTransformation(d, c, tc.spec)
				},
			)
		})
	}
}
