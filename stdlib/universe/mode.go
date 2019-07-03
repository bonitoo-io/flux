package universe

import (
	"fmt"
	"sort"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/values"
)

const ModeKind = "mode"

type ModeOpSpec struct {
	Column string `json:"column"`
}

func init() {
	modeSignature := flux.FunctionSignature(
		map[string]semantic.PolyType{
			"column": semantic.String,
		},
		nil,
	)

	flux.RegisterPackageValue("universe", ModeKind, flux.FunctionValue(ModeKind, createModeOpSpec, modeSignature))
	flux.RegisterOpSpec(ModeKind, newModeOp)
	plan.RegisterProcedureSpec(ModeKind, newModeProcedure, ModeKind)
	execute.RegisterTransformation(ModeKind, createModeTransformation)
}

func createModeOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(ModeOpSpec)

	if col, ok, err := args.GetString("column"); err != nil {
		return nil, err
	} else if ok {
		spec.Column = col
	} else {
		spec.Column = execute.DefaultValueColLabel
	}
	return spec, nil
}

func newModeOp() flux.OperationSpec {
	return new(ModeOpSpec)
}

func (s *ModeOpSpec) Kind() flux.OperationKind {
	return ModeKind
}

type ModeProcedureSpec struct {
	plan.DefaultCost
	Column string
}

func newModeProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*ModeOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	return &ModeProcedureSpec{
		Column: spec.Column,
	}, nil
}

func (s *ModeProcedureSpec) Kind() plan.ProcedureKind {
	return ModeKind
}
func (s *ModeProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(ModeProcedureSpec)

	*ns = *s

	return ns
}

// TriggerSpec implements plan.TriggerAwareProcedureSpec
func (s *ModeProcedureSpec) TriggerSpec() plan.TriggerSpec {
	return plan.NarrowTransformationTriggerSpec{}
}

func createModeTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*ModeProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewModeTransformation(d, cache, s)
	return t, d, nil
}

type modeTransformation struct {
	d     execute.Dataset
	cache execute.TableBuilderCache

	column string
}

func NewModeTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *ModeProcedureSpec) *modeTransformation {
	return &modeTransformation{
		d:      d,
		cache:  cache,
		column: spec.Column,
	}
}

func (t *modeTransformation) RetractTable(id execute.DatasetID, key flux.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *modeTransformation) Process(id execute.DatasetID, tbl flux.Table) error {
	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return fmt.Errorf("mode found duplicate table with key: %v", tbl.Key())
	}

	colIdx := execute.ColIdx(t.column, tbl.Cols())
	if colIdx < 0 {
		// doesn't exist in this table, so add an empty value
		if err := execute.AddTableKeyCols(tbl.Key(), builder); err != nil {
			return err
		}
		colIdx, err := builder.AddCol(flux.ColMeta{
			Label: execute.DefaultValueColLabel,
			Type:  flux.TString,
		})
		if err != nil {
			return err
		}

		if err := builder.AppendString(colIdx, ""); err != nil {
			return err
		}
		if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
			return err
		}
		// TODO: hack required to ensure data flows downstream
		return tbl.Do(func(flux.ColReader) error {
			return nil
		})
	}

	col := tbl.Cols()[colIdx]

	if err := execute.AddTableKeyCols(tbl.Key(), builder); err != nil {
		return err
	}
	colIdx, err := builder.AddCol(flux.ColMeta{
		Label: execute.DefaultValueColLabel,
		Type:  col.Type,
	})
	if err != nil {
		return err
	}

	if tbl.Key().HasCol(t.column) {
		j := execute.ColIdx(t.column, tbl.Key().Cols())

		if err := builder.AppendValue(colIdx, tbl.Key().Value(j)); err != nil {
			return err
		}

		if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
			return err
		}
		// TODO: hack required to ensure data flows downstream
		return tbl.Do(func(flux.ColReader) error {
			return nil
		})
	}

	var (
		numNil int64
	)

	switch col.Type {
	case flux.TBool:
		boolMode := make(map[bool]int64)
		return tbl.Do(func(cr flux.ColReader) error {
			l := cr.Len()
			err := t.doBool(cr, tbl, boolMode, l, numNil, builder, colIdx)

			if err != nil {
				return err
			}

			return nil
		})
	case flux.TInt:
		intMode := make(map[int64]int64)
		return tbl.Do(func(cr flux.ColReader) error {
			l := cr.Len()
			err := t.doInt(cr, tbl, intMode, l, numNil, builder, colIdx)

			if err != nil {
				return err
			}

			return nil
		})
	case flux.TUInt:
		uintMode := make(map[uint64]int64)
		return tbl.Do(func(cr flux.ColReader) error {
			l := cr.Len()
			err := t.doUInt(cr, tbl, uintMode, l, numNil, builder, colIdx)

			if err != nil {
				return err
			}

			return nil
		})
	case flux.TFloat:
		floatMode := make(map[float64]int64)
		return tbl.Do(func(cr flux.ColReader) error {
			l := cr.Len()
			err := t.doFloat(cr, tbl, floatMode, l, numNil, builder, colIdx)

			if err != nil {
				return err
			}

			return nil
		})
	case flux.TString:
		stringMode := make(map[string]int64)
		return tbl.Do(func(cr flux.ColReader) error {
			l := cr.Len()
			err := t.doString(cr, tbl, stringMode, l, numNil, builder, colIdx)

			if err != nil {
				return err
			}

			return nil
		})
	case flux.TTime:
		timeMode := make(map[execute.Time]int64)
		return tbl.Do(func(cr flux.ColReader) error {
			l := cr.Len()
			err := t.doTime(cr, tbl, timeMode, l, numNil, builder, colIdx)

			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
		return err
	}

	return nil
}

func (t *modeTransformation) doString(cr flux.ColReader, tbl flux.Table, stringMode map[string]int64, l int, numNil int64, builder execute.TableBuilder, colIdx int) error {
	j := execute.ColIdx(t.column, tbl.Cols())
	for i := 0; i < l; i++ {
		if cr.Strings(j).IsNull(i) {
			numNil++
			continue
		}
		v := cr.Strings(j).ValueString(i)
		stringMode[v]++
	}

	max, total := int64(0), int64(0)
	for val := range stringMode {
		if stringMode[val] > max {
			max, total = stringMode[val], 1
		} else if stringMode[val] == max {
			total++
		}
	}

	storedVals := make([]string, 0, total)
	for val := range stringMode {
		if stringMode[val] == max {
			storedVals = append(storedVals, val)
		}
	}

	/*
		// byte-slice method:
		for k := range stringMode {
			if stringMode[k] > max {
				storedVals = storedVals[:0]
				storedVals = append(storedVals, k)
				max = stringMode[k]
			} else if stringMode[k] == max {
				storedVals = append(storedVals, k)
			}
		}
	*/

	if numNil > max {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if len(storedVals) == len(stringMode) && (numNil == 0 || numNil == max) {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if numNil == max {
		sort.Strings(storedVals)
		for j := range storedVals {
			if err := builder.AppendString(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else {
		sort.Strings(storedVals)
		for j := range storedVals {
			if err := builder.AppendString(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
	}

	if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
		return err
	}

	return nil
}

func (t *modeTransformation) doBool(cr flux.ColReader, tbl flux.Table, boolMode map[bool]int64, l int, numNil int64, builder execute.TableBuilder, colIdx int) error {
	j := execute.ColIdx(t.column, tbl.Cols())
	for i := 0; i < l; i++ {
		if cr.Bools(j).IsNull(i) {
			numNil++
			continue
		}
		v := cr.Bools(j).Value(i)
		boolMode[v]++
	}

	max, total := int64(0), int64(0)
	for val := range boolMode {
		if boolMode[val] > max {
			max, total = boolMode[val], 1
		} else if boolMode[val] == max {
			total++
		}
	}

	storedVals := make([]bool, 0, total)
	for val := range boolMode {
		if boolMode[val] == max {
			storedVals = append(storedVals, val)
		}
	}

	// if more nils than the max occurrences of another value, mode is nil
	if numNil > max {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if len(storedVals) == len(boolMode) && (numNil == 0 || numNil == max) { // if no nils and all of them have max num of occurrences, no mode
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if numNil == max { // if max is the same as nil, mode is nil and the other max occurrences
		if len(storedVals) == 3 {
			if err := builder.AppendBool(colIdx, true); err != nil {
				return err
			}
			if err := builder.AppendBool(colIdx, false); err != nil {
				return err
			}
		} else {
			for j := range storedVals {
				if err := builder.AppendBool(colIdx, storedVals[j]); err != nil {
					return err
				}
			}
		}
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else {
		if len(storedVals) == 2 {
			if err := builder.AppendBool(colIdx, true); err != nil {
				return err
			}
			if err := builder.AppendBool(colIdx, false); err != nil {
				return err
			}
		} else {
			for j := range storedVals {
				if err := builder.AppendBool(colIdx, storedVals[j]); err != nil {
					return err
				}
			}
		}
	}

	if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
		return err
	}

	return nil
}

func (t *modeTransformation) doInt(cr flux.ColReader, tbl flux.Table, intMode map[int64]int64, l int, numNil int64, builder execute.TableBuilder, colIdx int) error {
	j := execute.ColIdx(t.column, tbl.Cols())
	for i := 0; i < l; i++ {
		if cr.Ints(j).IsNull(i) {
			numNil++
			continue
		}
		v := cr.Ints(j).Value(i)
		intMode[v]++
	}

	max, total := int64(0), int64(0)
	for val := range intMode {
		if intMode[val] > max {
			max, total = intMode[val], 1
		} else if intMode[val] == max {
			total++
		}
	}

	storedVals := make([]int64, 0, total)
	for val := range intMode {
		if intMode[val] == max {
			storedVals = append(storedVals, val)
		}
	}

	if numNil > max {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if len(storedVals) == len(intMode) && (numNil == 0 || numNil == max) {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if numNil == max {
		sort.Slice(storedVals, func(i, j int) bool { return storedVals[i] < storedVals[j] })
		for j := range storedVals {
			if err := builder.AppendInt(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else {
		sort.Slice(storedVals, func(i, j int) bool { return storedVals[i] < storedVals[j] })
		for j := range storedVals {
			if err := builder.AppendInt(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
	}
	if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
		return err
	}

	return nil
}

func (t *modeTransformation) doUInt(cr flux.ColReader, tbl flux.Table, uintMode map[uint64]int64, l int, numNil int64, builder execute.TableBuilder, colIdx int) error {
	j := execute.ColIdx(t.column, tbl.Cols())
	for i := 0; i < l; i++ {
		if cr.UInts(j).IsNull(i) {
			numNil++
			continue
		}
		v := cr.UInts(j).Value(i)
		uintMode[v]++
	}

	max, total := int64(0), int64(0)
	for val := range uintMode {
		if uintMode[val] > max {
			max, total = uintMode[val], 1
		} else if uintMode[val] == max {
			total++
		}
	}

	storedVals := make([]uint64, 0, total)
	for val := range uintMode {
		if uintMode[val] == max {
			storedVals = append(storedVals, val)
		}
	}

	if numNil > max {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if len(storedVals) == len(uintMode) && (numNil == 0 || numNil == max) {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if numNil == max {
		sort.Slice(storedVals, func(i, j int) bool { return storedVals[i] < storedVals[j] })
		for j := range storedVals {
			if err := builder.AppendUInt(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else {
		sort.Slice(storedVals, func(i, j int) bool { return storedVals[i] < storedVals[j] })
		for j := range storedVals {
			if err := builder.AppendUInt(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
	}
	if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
		return err
	}

	return nil
}

func (t *modeTransformation) doFloat(cr flux.ColReader, tbl flux.Table, floatMode map[float64]int64, l int, numNil int64, builder execute.TableBuilder, colIdx int) error {
	j := execute.ColIdx(t.column, tbl.Cols())
	for i := 0; i < l; i++ {
		if cr.Floats(j).IsNull(i) {
			numNil++
			continue
		}
		v := cr.Floats(j).Value(i)
		floatMode[v]++
	}

	max, total := int64(0), int64(0)
	for val := range floatMode {
		if floatMode[val] > max {
			max, total = floatMode[val], 1
		} else if floatMode[val] == max {
			total++
		}
	}

	storedVals := make([]float64, 0, total)
	for val := range floatMode {
		if floatMode[val] == max {
			storedVals = append(storedVals, val)
		}
	}

	if numNil > max {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if len(storedVals) == len(floatMode) && (numNil == 0 || numNil == max) {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if numNil == max {
		sort.Float64s(storedVals)
		for j := range storedVals {
			if err := builder.AppendFloat(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else {
		sort.Float64s(storedVals)
		for j := range storedVals {
			if err := builder.AppendFloat(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
	}
	if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
		return err
	}

	return nil
}

func (t *modeTransformation) doTime(cr flux.ColReader, tbl flux.Table, timeMode map[execute.Time]int64, l int, numNil int64, builder execute.TableBuilder, colIdx int) error {
	j := execute.ColIdx(t.column, tbl.Cols())
	for i := 0; i < l; i++ {
		if cr.Times(j).IsNull(i) {
			numNil++
			continue
		}
		v := values.Time(cr.Times(j).Value(i))
		timeMode[v]++
	}

	max, total := int64(0), int64(0)
	for val := range timeMode {
		if timeMode[val] > max {
			max, total = timeMode[val], 1
		} else if timeMode[val] == max {
			total++
		}
	}

	storedVals := make([]execute.Time, 0, total)
	for val := range timeMode {
		if timeMode[val] == max {
			storedVals = append(storedVals, val)
		}
	}

	if numNil > max {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if len(storedVals) == len(timeMode) && (numNil == 0 || numNil == max) {
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else if numNil == max {
		sort.Slice(storedVals, func(i, j int) bool { return storedVals[i] < storedVals[j] })
		for j := range storedVals {
			if err := builder.AppendTime(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
		if err := builder.AppendNil(colIdx); err != nil {
			return err
		}
	} else {
		sort.Slice(storedVals, func(i, j int) bool { return storedVals[i] < storedVals[j] })
		for j := range storedVals {
			if err := builder.AppendTime(colIdx, storedVals[j]); err != nil {
				return err
			}
		}
	}
	if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
		return err
	}

	return nil
}

func (t *modeTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}
func (t *modeTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}
func (t *modeTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
