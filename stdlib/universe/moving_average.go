package universe

import (
	"container/list"
	"fmt"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/interpreter"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
	"github.com/pkg/errors"
)

const MovingAverageKind = "movingAverage"

type MovingAverageOpSpec struct {
	N       int64    `json:"n"`
	Columns []string `json:"columns"`
}

func init() {
	movingAverageSignature := flux.FunctionSignature(
		map[string]semantic.PolyType{
			"n":       semantic.Int,
			"columns": semantic.NewArrayPolyType(semantic.String),
		},
		[]string{"n"},
	)

	flux.RegisterPackageValue("universe", MovingAverageKind, flux.FunctionValue(MovingAverageKind, createMovingAverageOpSpec, movingAverageSignature))
	flux.RegisterOpSpec(MovingAverageKind, newMovingAverageOp)
	plan.RegisterProcedureSpec(MovingAverageKind, newMovingAverageProcedure, MovingAverageKind)
	execute.RegisterTransformation(MovingAverageKind, createMovingAverageTransformation)
}

func createMovingAverageOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(MovingAverageOpSpec)

	if n, err := args.GetRequiredInt("n"); err != nil {
		return nil, err
	} else {
		spec.N = n
	}

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

func newMovingAverageOp() flux.OperationSpec {
	return new(MovingAverageOpSpec)
}

func (s *MovingAverageOpSpec) Kind() flux.OperationKind {
	return MovingAverageKind
}

type MovingAverageProcedureSpec struct {
	plan.DefaultCost
	N          int64    `json:"n"`
	Columns    []string `json:"columns"`
	TimeColumn string   `json:"timeColumn"`
}

func newMovingAverageProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*MovingAverageOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &MovingAverageProcedureSpec{
		N:       spec.N,
		Columns: spec.Columns,
	}, nil
}

func (s *MovingAverageProcedureSpec) Kind() plan.ProcedureKind {
	return MovingAverageKind
}

func (s *MovingAverageProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(MovingAverageProcedureSpec)
	*ns = *s
	if s.Columns != nil {
		ns.Columns = make([]string, len(s.Columns))
		copy(ns.Columns, s.Columns)
	}
	return ns
}

// TriggerSpec implements plan.TriggerAwareProcedureSpec
func (s *MovingAverageProcedureSpec) TriggerSpec() plan.TriggerSpec {
	return plan.NarrowTransformationTriggerSpec{}
}

func createMovingAverageTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*MovingAverageProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewMovingAverageTransformation(d, cache, s)
	return t, d, nil
}

type movingAverageTransformation struct {
	d     execute.Dataset
	cache execute.TableBuilderCache

	n       int64
	columns []string
}

func NewMovingAverageTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *MovingAverageProcedureSpec) *movingAverageTransformation {
	return &movingAverageTransformation{
		d:       d,
		cache:   cache,
		n:       spec.N,
		columns: spec.Columns,
	}
}

func (t *movingAverageTransformation) RetractTable(id execute.DatasetID, key flux.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *movingAverageTransformation) Process(id execute.DatasetID, tbl flux.Table) error {
	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return fmt.Errorf("moving average found duplicate table with key: %v", tbl.Key())
	}
	timeIdx := -1
	cols := tbl.Cols()
	doMovingAverage := make([]bool, len(cols))
	for j, c := range cols {
		found := false
		for _, label := range t.columns {
			if c.Label == label {
				found = true
				break
			}
		}
		if c.Label == execute.DefaultTimeColLabel {
			timeIdx = j
		}

		if found {
			mac := c
			mac.Type = flux.TFloat
			_, err := builder.AddCol(mac)
			if err != nil {
				return err
			}
			doMovingAverage[j] = true
		} else {
			_, err := builder.AddCol(c)
			if err != nil {
				return err
			}
		}
	}
	if timeIdx < 0 {
		return fmt.Errorf("no column %q exists", execute.DefaultTimeColLabel)
	}

	return tbl.Do(func(cr flux.ColReader) error {
		if cr.Len() == 0 {
			return nil
		}

		if cr.Times(timeIdx).NullN() > 0 {
			return fmt.Errorf("moving average found null time in time column")
		}

		for j, c := range cr.Cols() {
			var err error
			switch c.Type {
			case flux.TBool:
				err = t.passThroughBool(cr.Times(timeIdx), cr.Bools(j), builder, j)
			case flux.TInt:
				err = t.doInt(cr.Times(timeIdx), cr.Ints(j), builder, j, doMovingAverage[j])
			case flux.TUInt:
				err = t.doUInt(cr.Times(timeIdx), cr.UInts(j), builder, j, doMovingAverage[j])
			case flux.TFloat:
				err = t.doFloat(cr.Times(timeIdx), cr.Floats(j), builder, j, doMovingAverage[j])
			case flux.TString:
				err = t.passThroughString(cr.Times(timeIdx), cr.Strings(j), builder, j)
			case flux.TTime:
				err = t.passThroughTime(cr.Times(timeIdx), cr.Times(j), builder, j)
			}

			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *movingAverageTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}

func (t *movingAverageTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}

func (t *movingAverageTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}

const movingAverageUnsortedTimeErr = "moving average found out-of-order times in time column"

func (t *movingAverageTransformation) passThroughBool(ts *array.Int64, vs *array.Boolean, b execute.TableBuilder, bj int) error {
	i := 0
	pTime := execute.Time(ts.Value(i))
	i++
	count := int64(1)

	// Process the rest of the rows
	l := vs.Len()
	for ; i < l; i++ {
		cTime := execute.Time(ts.Value(i))
		if cTime < pTime {
			return errors.New(movingAverageUnsortedTimeErr)
		}

		if cTime == pTime {
			continue
		}

		count++

		if count < t.n {
			continue
		}

		// We have a valid time for this row.  Code below should not exit
		// the loop early so pTime can be set for the next iteration.

		if vs.IsValid(i) {
			if err := b.AppendBool(bj, vs.Value(i)); err != nil {
				return err
			}
		} else {
			if err := b.AppendNil(bj); err != nil {
				return err
			}
		}
		pTime = cTime
	}

	return nil
}

func (t *movingAverageTransformation) doInt(ts *array.Int64, vs *array.Int64, b execute.TableBuilder, bj int, doMovingAverage bool) error {
	i := 0
	var pValue int64
	validPValue := false
	window := list.New()
	sum := int64(0)
	nullCount := float64(0)

	pTime := execute.Time(ts.Value(i))
	if vs.IsValid(i) {
		pValue = vs.Value(i)
		validPValue = true
		window.PushBack(pValue)
		sum += pValue
	}
	i++

	l := vs.Len()
	for ; i < l; i++ {
		cTime := execute.Time(ts.Value(i))
		if cTime < pTime {
			return errors.New(movingAverageUnsortedTimeErr)
		}

		if cTime == pTime {
			continue
		}

		if vs.IsNull(i) {
			window.PushBack(nil)
			nullCount++
		} else {
			window.PushBack(vs.Value(i))
		}
		sum += vs.Value(i)

		if int64(window.Len()) < t.n {
			continue
		}

		if !doMovingAverage {
			if vs.IsValid(i) {
				if err := b.AppendInt(bj, vs.Value(i)); err != nil {
					return err
				}
			} else {
				if err := b.AppendNil(bj); err != nil {
					return err
				}
			}
		} else if !validPValue {
			if err := b.AppendNil(bj); err != nil {
				return err
			}

			if vs.IsValid(i) {
				pTime = cTime
				validPValue = true
			}
		} else {
			e := window.Front()
			average := float64(sum) / (float64(window.Len()) - nullCount)
			if e.Value != nil {
				sum -= e.Value.(int64)
			}
			window.Remove(e)
			if err := b.AppendFloat(bj, average); err != nil {
				return err
			}
		}
		pTime = cTime
	}
	return nil
}

func (t *movingAverageTransformation) doUInt(ts *array.Int64, vs *array.Uint64, b execute.TableBuilder, bj int, doMovingAverage bool) error {
	i := 0
	var pValue uint64
	validPValue := false
	window := list.New()
	sum := uint64(0)
	nullCount := float64(0)

	pTime := execute.Time(ts.Value(i))
	if vs.IsValid(i) {
		pValue = vs.Value(i)
		validPValue = true
		window.PushBack(pValue)
		sum += pValue
	}
	i++

	l := vs.Len()
	for ; i < l; i++ {
		cTime := execute.Time(ts.Value(i))
		if cTime < pTime {
			return errors.New(movingAverageUnsortedTimeErr)
		}

		if cTime == pTime {
			continue
		}

		if vs.IsNull(i) {
			window.PushBack(nil)
			nullCount++
		} else {
			window.PushBack(vs.Value(i))
		}
		sum += vs.Value(i)

		if int64(window.Len()) < t.n {
			continue
		}

		if !doMovingAverage {
			if vs.IsValid(i) {
				if err := b.AppendUInt(bj, vs.Value(i)); err != nil {
					return err
				}
			} else {
				if err := b.AppendNil(bj); err != nil {
					return err
				}
			}
		} else if !validPValue {
			if err := b.AppendNil(bj); err != nil {
				return err
			}

			if vs.IsValid(i) {
				pTime = cTime
				validPValue = true
			}
		} else {
			e := window.Front()
			average := float64(sum) / (float64(window.Len()) - nullCount)
			if e.Value != nil {
				sum -= e.Value.(uint64)
			}
			window.Remove(e)
			if err := b.AppendFloat(bj, average); err != nil {
				return err
			}
		}
		pTime = cTime
	}
	return nil
}

func (t *movingAverageTransformation) doFloat(ts *array.Int64, vs *array.Float64, b execute.TableBuilder, bj int, doMovingAverage bool) error {
	i := 0
	var pValue float64
	validPValue := false
	window := list.New()
	sum := float64(0)
	nullCount := float64(0)

	pTime := execute.Time(ts.Value(i))
	if vs.IsValid(i) {
		pValue = vs.Value(i)
		validPValue = true
		window.PushBack(pValue)
		sum += pValue
	}
	i++

	l := vs.Len()
	for ; i < l; i++ {
		cTime := execute.Time(ts.Value(i))
		if cTime < pTime {
			return errors.New(movingAverageUnsortedTimeErr)
		}

		if cTime == pTime {
			continue
		}

		if vs.IsNull(i) {
			window.PushBack(nil)
			nullCount++
		} else {
			window.PushBack(vs.Value(i))
		}
		sum += vs.Value(i)

		if int64(window.Len()) < t.n {
			continue
		}

		if !doMovingAverage {
			if vs.IsValid(i) {
				if err := b.AppendFloat(bj, vs.Value(i)); err != nil {
					return err
				}
			} else {
				if err := b.AppendNil(bj); err != nil {
					return err
				}
			}
		} else if !validPValue {
			if err := b.AppendNil(bj); err != nil {
				return err
			}

			if vs.IsValid(i) {
				pTime = cTime
				validPValue = true
			}
		} else {
			e := window.Front()
			average := float64(sum) / (float64(window.Len()) - nullCount)
			if e.Value != nil {
				sum -= e.Value.(float64)
			}
			window.Remove(e)
			if err := b.AppendFloat(bj, average); err != nil {
				return err
			}
		}
		pTime = cTime
	}
	return nil
}

func (t *movingAverageTransformation) passThroughString(ts *array.Int64, vs *array.Binary, b execute.TableBuilder, bj int) error {
	i := 0
	pTime := execute.Time(ts.Value(i))
	i++
	count := int64(1)

	// Process the rest of the rows
	l := vs.Len()
	for ; i < l; i++ {
		cTime := execute.Time(ts.Value(i))
		if cTime < pTime {
			return errors.New(movingAverageUnsortedTimeErr)
		}

		if cTime == pTime {
			// Only use the first value found if a time value is the same as
			// the previous row.
			continue
		}

		count++

		if int64(count) < t.n {
			continue
		}

		// We have a valid time for this row.  Code below should not exit
		// the loop early so pTime can be set for the next iteration.

		if vs.IsValid(i) {
			if err := b.AppendString(bj, string(vs.Value(i))); err != nil {
				return err
			}
		} else {
			if err := b.AppendNil(bj); err != nil {
				return err
			}
		}

		pTime = cTime
	}

	return nil
}

func (t *movingAverageTransformation) passThroughTime(ts *array.Int64, vs *array.Int64, b execute.TableBuilder, bj int) error {
	i := 0
	pTime := execute.Time(ts.Value(i))
	i++
	count := int64(1)

	// Process the rest of the rows
	l := vs.Len()
	for ; i < l; i++ {
		cTime := execute.Time(ts.Value(i))
		if cTime < pTime {
			return errors.New(movingAverageUnsortedTimeErr)
		}

		if cTime == pTime {
			// Only use the first value found if a time value is the same as
			// the previous row.
			continue
		}

		count++

		if count < t.n {
			continue
		}

		// We have a valid time for this row.  Code below should not exit
		// the loop early so pTime can be set for the next iteration.

		if vs.IsValid(i) {
			if err := b.AppendTime(bj, execute.Time(vs.Value(i))); err != nil {
				return err
			}
		} else {
			if err := b.AppendNil(bj); err != nil {
				return err
			}
		}

		pTime = cTime
	}

	return nil
}
