package universe

import (
	"fmt"
	"github.com/influxdata/flux/interpreter"
	"math"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
)

const ElapsedKind = "elapsed"

type ElapsedOpSpec struct {
	Columns     []string `json:"columns"`
}

func init() {
	elapsedSignature := flux.FunctionSignature(
		map[string]semantic.PolyType{
			"columns":     semantic.NewArrayPolyType(semantic.String),
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

	if cols, ok, err := args.GetArray("columns", semantic.String); err != nil {
		return nil, err
	} else if ok {
		columns, err := interpreter.ToStringArray(cols)
		if err != nil {
			return nil, err
		}
		spec.Columns = columns
	} else {
		spec.Columns = []string{execute.DefaultValueColLabel}
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
	Columns     []string `json:"columns"`
}

func newElapsedProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*ElapsedOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &ElapsedProcedureSpec{
		Columns:     spec.Columns,
	}, nil
}

func (s *ElapsedProcedureSpec) Kind() plan.ProcedureKind {
	return ElapsedKind
}
func (s *ElapsedProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(ElapsedProcedureSpec)
	*ns = *s
	if s.Columns != nil {
		ns.Columns = make([]string, len(s.Columns))
		copy(ns.Columns, s.Columns)
	}
	return ns
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
	columns     []string
}

func NewElapsedTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *ElapsedProcedureSpec) *elapsedTransformation {
	return &elapsedTransformation{
		d:           d,
		cache:       cache,
		columns:     spec.Columns,
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
