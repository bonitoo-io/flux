package universe

import (
	"fmt"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/arrow"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/interpreter"
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
	elapsedvalues := make([]*elapsed, len(cols))
	for j, c := range cols {
		found := c.Label == t.timeColumn

		if found {
			var typ flux.ColType
			switch c.Type {
			case flux.TInt, flux.TUInt:
				typ = flux.TInt
			case flux.TFloat:
				typ = flux.TFloat
			}
			if _, err := builder.AddCol(flux.ColMeta{
				Label: c.Label,
				Type:  typ,
			}); err != nil {
				return err
			}
			elapsedvalues[j] = newElapsed(j)
		} else {
			_, err := builder.AddCol(c)
			if err != nil {
				return err
			}
		}
	}

	firstIdx := 1
	return tbl.Do(func(cr flux.ColReader) error {
		l := cr.Len()

		if l != 0 {
			for j, c := range cols {
				d := elapsedvalues[j]
				switch c.Type {
				case flux.TBool:
					s := arrow.BoolSlice(cr.Bools(j), firstIdx, l)
					if err := builder.AppendBools(j, s); err != nil {
						s.Release()
						return err
					}
					s.Release()
				case flux.TInt:
					if d != nil {
						for i := 0; i < l; i++ {
							if vs := cr.Ints(j); vs.IsValid(i) {
								if v, first := d.updateInt(vs.Value(i)); !first {
									if d.nonNegative && v < 0 {
										if err := builder.AppendNil(j); err != nil {
											return err
										}
									} else {
										if err := builder.AppendInt(j, v); err != nil {
											return err
										}
									}
								}
							} else if err := builder.AppendNil(j); err != nil {
								return err
							}
						}
					} else {
						s := arrow.IntSlice(cr.Ints(j), firstIdx, l)
						if err := builder.AppendInts(j, s); err != nil {
							s.Release()
							return err
						}
						s.Release()
					}
				case flux.TUInt:
					if d != nil {
						for i := 0; i < l; i++ {
							if vs := cr.UInts(j); vs.IsValid(i) {
								if v, first := d.updateUInt(vs.Value(i)); !first {
									if d.nonNegative && v < 0 {
										if err := builder.AppendNil(j); err != nil {
											return err
										}
									} else {
										if err := builder.AppendInt(j, v); err != nil {
											return err
										}
									}
								}
							} else if err := builder.AppendNil(j); err != nil {
								return err
							}
						}
					} else {
						s := arrow.UintSlice(cr.UInts(j), firstIdx, l)
						if err := builder.AppendUInts(j, s); err != nil {
							s.Release()
							return err
						}
						s.Release()
					}
				case flux.TFloat:
					if d != nil {
						for i := 0; i < l; i++ {
							if vs := cr.Floats(j); vs.IsValid(i) {
								if v, first := d.updateFloat(vs.Value(i)); !first {
									if d.nonNegative && v < 0 {
										if err := builder.AppendNil(j); err != nil {
											return err
										}
									} else {
										if err := builder.AppendFloat(j, v); err != nil {
											return err
										}
									}
								}
							} else if err := builder.AppendNil(j); err != nil {
								return err
							}
						}
					} else {
						s := arrow.FloatSlice(cr.Floats(j), firstIdx, l)
						if err := builder.AppendFloats(j, s); err != nil {
							s.Release()
							return err
						}
						s.Release()
					}
				case flux.TString:
					s := arrow.StringSlice(cr.Strings(j), firstIdx, l)
					if err := builder.AppendStrings(j, s); err != nil {
						s.Release()
						return err
					}
					s.Release()
				case flux.TTime:
					s := arrow.IntSlice(cr.Times(j), firstIdx, l)
					if err := builder.AppendTimes(j, s); err != nil {
						s.Release()
						return err
					}
					s.Release()
				}
			}
		}

		// Now that we skipped the first row, start at 0 for the rest of the batches
		firstIdx = 0
		return nil
	})
}