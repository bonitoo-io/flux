package universe

import (
	"github.com/influxdata/flux/querytest"
	"testing"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/execute/executetest"
)

func TestElapsedOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"elapsed","kind":"elapsed","spec":{TimeColumn: "time"}}`)
	op := &flux.Operation{
		ID: "elapsed",
		Spec: &ElapsedOpSpec{
			TimeColumn: "time",
		},
	}
	querytest.OperationMarshalingTestHelper(t, data, op)
}

func TestElapsed_PassThrough(t *testing.T) {
	executetest.TransformationPassThroughTestHelper(t, func(d execute.Dataset, c execute.TableBuilderCache) execute.Transformation {
		s := NewElapsedTransformation(
			d,
			c,
			&ElapsedProcedureSpec{},
		)
		return s
	})
}


func TestElapsed_Process(t *testing.T) {
	testCases := []struct {
		name string
		spec *ElapsedProcedureSpec
		data []flux.Table
		want []*executetest.Table
	}{
		{
			name: "basic",
			spec: &ElapsedProcedureSpec{
				TimeColumn: execute.DefaultTimeColLabel,
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 2.0},
					{execute.Time(2), 1.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(2), -1.0},
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
					return NewElapsedTransformation(d, c, tc.spec)
				},
			)
		})
	}
}
