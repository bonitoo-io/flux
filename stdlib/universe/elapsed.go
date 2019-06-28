package universe

import (
	"fmt"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
)

const ElapsedKind = "elapsed"

type ElapsedOpSpec struct {
	TimeColumn	string	 `json:"timeColumn"`
}

func init() {
	elapsedSignature := flux.FunctionSignature(
		map[string]semantic.PolyType{
			"timeColumn":  semantic.String,
		},
		nil,
	)

	flux.RegisterPackageValue("universe", ElapsedKind, flux.FunctionValue(ElapsedKind, createElapsedOpSpec, elapsedSignature))
	flux.RegisterOpSpec(ElapsedKind, newElapsedOp)
	plan.RegisterProcedureSpec(ElapsedKind, newElapsedProcedure, ElapsedKind)
	execute.RegisterTransformation(ElapsedKind, createElapsedTransformation)
}

func createElapsedOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	err := a.AddParentFromArgs(args)
	if err != nil {
		return nil, err
	}

	spec := new(ElapsedOpSpec)

	if nn, ok, err := args.GetString("timeColumn"); err != nil {
		return nil, err
	} else if ok {
		spec.TimeColumn = nn
	} else {
		spec.TimeColumn = execute.DefaultTimeColLabel
	}

	return spec, nil
}

func newElapsedOp() flux.OperationSpec {
	return new(ElapsedOpSpec)
}

func (s *ElapsedOpSpec) Kind() flux.OperationKind {
	return ElapsedKind
}

type ElapsedProcedureSpec struct {
	plan.DefaultCost
	TimeColumn	string	 `json:"timeColumn"`
}

func newElapsedProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*ElapsedOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &ElapsedProcedureSpec{
		TimeColumn:  spec.TimeColumn,
	}, nil
}

func (s *ElapsedProcedureSpec) Kind() plan.ProcedureKind {
	return ElapsedKind
}

func (s *ElapsedProcedureSpec) Copy() plan.ProcedureSpec {
	return &ElapsedProcedureSpec{
		TimeColumn: s.TimeColumn,
	}
}

func createElapsedTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*ElapsedProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewElapsedTransformation(d, cache, s)
	return t, d, nil
}

type elapsedTransformation struct {
	d     execute.Dataset
	cache execute.TableBuilderCache

	timeColumn	string
}

func NewElapsedTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *ElapsedProcedureSpec) *elapsedTransformation {
	return &elapsedTransformation{
		d:           d,
		cache:       cache,

		timeColumn:  spec.TimeColumn,
	}
}

func (t *elapsedTransformation) RetractTable(id execute.DatasetID, key flux.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *elapsedTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}

func (t *elapsedTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}

func (t *elapsedTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}

func (t *elapsedTransformation) Process(id execute.DatasetID, tbl flux.Table) error {
	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return fmt.Errorf("found duplicate table with key: %v", tbl.Key())
	}
	cols := tbl.Cols()
	for _, c := range cols {
		found := c.Label == t.timeColumn

		if found {
			var typ flux.ColType
			if c.Type == flux.TTime {
				typ = flux.TFloat
			}

			_, err := builder.AddCol(c)
			if err != nil {
				return err
			}

			if _, err := builder.AddCol(flux.ColMeta{
				Label: "elapsed",
				Type:  typ,
			}); err != nil {
				return err
			}
		}
	}

	return tbl.Do(func(cr flux.ColReader) error {
		l := cr.Len()

		if l != 0 {
			for j, c := range cols {
				if c.Type == flux.TTime {
						ts := cr.Times(j)
						prevTime := float64(execute.Time(ts.Value(0)))
						currTime := 0.0
						for i := 1; i < l; i++ {
							pTime := execute.Time(ts.Value(i))
							currTime = float64(pTime)

							if err := builder.AppendTime(0, pTime); err != nil {
								return err
							}

							if err := builder.AppendFloat(1, currTime - prevTime); err != nil {
								return err
							}
							prevTime = currTime
						}
				}
			}
		}

		// Now that we skipped the first row, start at 0 for the rest of the batches
		return nil
	})
}