package cli

import (
	"flag"
)

type (
	// SliceFlag extends implementations like StringSliceFlag and IntSliceFlag with support for using slices directly,
	// as Value and/or Destination.
	// See also SliceFlagTarget, MultiStringFlag, MultiFloat64Flag, MultiInt64Flag, MultiIntFlag.
	SliceFlag[T SliceFlagTarget[E], S ~[]E, E any] struct {
		Target      T
		Value       S
		Destination *S
	}

	// SliceFlagTarget models a target implementation for use with SliceFlag.
	// The three methods, SetValue, SetDestination, and GetDestination, are necessary to propagate Value and
	// Destination, where Value is propagated inwards (initially), and Destination is propagated outwards (on every
	// update).
	SliceFlagTarget[E any] interface {
		Flag
		RequiredFlag
		DocGenerationFlag
		VisibleFlag
		CategorizableFlag

		// SetValue should propagate the given slice to the target, ideally as a new value.
		// Note that a nil slice should nil/clear any existing value (modelled as ~[]E).
		SetValue(slice []E)
		// SetDestination should propagate the given slice to the target, ideally as a new value.
		// Note that a nil slice should nil/clear any existing value (modelled as ~*[]E).
		SetDestination(slice []E)
		// GetDestination should return the current value referenced by any destination, or nil if nil/unset.
		GetDestination() []E
	}

	// MultiStringFlag extends StringSliceFlag with support for using slices directly, as Value and/or Destination.
	// See also SliceFlag.
	MultiStringFlag = SliceFlag[*StringSliceFlag, []string, string]

	// MultiFloat64Flag extends Float64SliceFlag with support for using slices directly, as Value and/or Destination.
	// See also SliceFlag.
	MultiFloat64Flag = SliceFlag[*Float64SliceFlag, []float64, float64]

	// MultiInt64Flag extends Int64SliceFlag with support for using slices directly, as Value and/or Destination.
	// See also SliceFlag.
	MultiInt64Flag = SliceFlag[*Int64SliceFlag, []int64, int64]

	// MultiIntFlag extends IntSliceFlag with support for using slices directly, as Value and/or Destination.
	// See also SliceFlag.
	MultiIntFlag = SliceFlag[*IntSliceFlag, []int, int]

	flagValueHook struct {
		value Generic
		hook  func()
	}
)

var (
	// compile time assertions

	_ SliceFlagTarget[string]  = (*StringSliceFlag)(nil)
	_ SliceFlagTarget[string]  = (*SliceFlag[*StringSliceFlag, []string, string])(nil)
	_ SliceFlagTarget[string]  = (*MultiStringFlag)(nil)
	_ SliceFlagTarget[float64] = (*MultiFloat64Flag)(nil)
	_ SliceFlagTarget[int64]   = (*MultiInt64Flag)(nil)
	_ SliceFlagTarget[int]     = (*MultiIntFlag)(nil)

	_ Generic    = (*flagValueHook)(nil)
	_ Serializer = (*flagValueHook)(nil)
)

func (x *SliceFlag[T, S, E]) Apply(set *flag.FlagSet) error {
	x.Target.SetValue(x.convertSlice(x.Value))

	destination := x.Destination
	if destination == nil {
		x.Target.SetDestination(nil)

		return x.Target.Apply(set)
	}

	x.Target.SetDestination(x.convertSlice(*destination))

	return applyFlagValueHook(set, x.Target.Apply, func() {
		*destination = x.Target.GetDestination()
	})
}

func (x *SliceFlag[T, S, E]) convertSlice(slice S) []E {
	result := make([]E, len(slice))
	copy(result, slice)
	return result
}

func (x *SliceFlag[T, S, E]) SetValue(slice S) {
	x.Value = slice
}

func (x *SliceFlag[T, S, E]) SetDestination(slice S) {
	if slice != nil {
		x.Destination = &slice
	} else {
		x.Destination = nil
	}
}

func (x *SliceFlag[T, S, E]) GetDestination() S {
	if destination := x.Destination; destination != nil {
		return *destination
	}
	return nil
}

func (x *SliceFlag[T, S, E]) String() string         { return x.Target.String() }
func (x *SliceFlag[T, S, E]) Names() []string        { return x.Target.Names() }
func (x *SliceFlag[T, S, E]) IsSet() bool            { return x.Target.IsSet() }
func (x *SliceFlag[T, S, E]) IsRequired() bool       { return x.Target.IsRequired() }
func (x *SliceFlag[T, S, E]) TakesValue() bool       { return x.Target.TakesValue() }
func (x *SliceFlag[T, S, E]) GetUsage() string       { return x.Target.GetUsage() }
func (x *SliceFlag[T, S, E]) GetValue() string       { return x.Target.GetValue() }
func (x *SliceFlag[T, S, E]) GetDefaultText() string { return x.Target.GetDefaultText() }
func (x *SliceFlag[T, S, E]) GetEnvVars() []string   { return x.Target.GetEnvVars() }
func (x *SliceFlag[T, S, E]) IsVisible() bool        { return x.Target.IsVisible() }
func (x *SliceFlag[T, S, E]) GetCategory() string    { return x.Target.GetCategory() }

func (x *flagValueHook) Set(value string) error {
	if err := x.value.Set(value); err != nil {
		return err
	}
	x.hook()
	return nil
}

// String returns the string representation of the underlying flag value.
// Note: This simplified implementation always returns the value's string representation,
// unlike the standard library's flag package which might omit default values in help text.
func (x *flagValueHook) String() string {
	if x.value == nil {
		return ""
	}
	return x.value.String()
}

func (x *flagValueHook) Serialize() string {
	if value, ok := x.value.(Serializer); ok {
		return value.Serialize()
	}
	return x.String()
}

// applyFlagValueHook wraps calls apply then wraps flags to call a hook function on update and after initial apply.
func applyFlagValueHook(set *flag.FlagSet, apply func(set *flag.FlagSet) error, hook func()) error {
	if apply == nil || set == nil || hook == nil {
		panic(`invalid input`)
	}
	var tmp flag.FlagSet
	if err := apply(&tmp); err != nil {
		return err
	}
	tmp.VisitAll(func(f *flag.Flag) { set.Var(&flagValueHook{value: f.Value, hook: hook}, f.Name, f.Usage) })
	hook()
	return nil
}

// newSliceFlagValue is for implementing SliceFlagTarget.SetValue and SliceFlagTarget.SetDestination.
// It's e.g. as part of StringSliceFlag.SetValue, using the factory NewStringSlice.
func newSliceFlagValue[R any, S ~[]E, E any](factory func(defaults ...E) *R, defaults S) *R {
	if defaults == nil {
		return nil
	}
	return factory(defaults...)
}

// unwrapFlagValue strips any/all *flagValueHook wrappers.
func unwrapFlagValue(v flag.Value) flag.Value {
	for {
		h, ok := v.(*flagValueHook)
		if !ok {
			return v
		}
		v = h.value
	}
}

// NOTE: the methods below are in this file to make use of the build constraint

func (f *Float64SliceFlag) SetValue(slice []float64) {
	f.Value = newSliceFlagValue(NewFloat64Slice, slice)
}

func (f *Float64SliceFlag) SetDestination(slice []float64) {
	f.Destination = newSliceFlagValue(NewFloat64Slice, slice)
}

func (f *Float64SliceFlag) GetDestination() []float64 {
	if destination := f.Destination; destination != nil {
		return destination.Value()
	}
	return nil
}

func (f *Int64SliceFlag) SetValue(slice []int64) {
	f.Value = newSliceFlagValue(NewInt64Slice, slice)
}

func (f *Int64SliceFlag) SetDestination(slice []int64) {
	f.Destination = newSliceFlagValue(NewInt64Slice, slice)
}

func (f *Int64SliceFlag) GetDestination() []int64 {
	if destination := f.Destination; destination != nil {
		return destination.Value()
	}
	return nil
}

func (f *IntSliceFlag) SetValue(slice []int) {
	f.Value = newSliceFlagValue(NewIntSlice, slice)
}

func (f *IntSliceFlag) SetDestination(slice []int) {
	f.Destination = newSliceFlagValue(NewIntSlice, slice)
}

func (f *IntSliceFlag) GetDestination() []int {
	if destination := f.Destination; destination != nil {
		return destination.Value()
	}
	return nil
}

func (f *StringSliceFlag) SetValue(slice []string) {
	f.Value = newSliceFlagValue(NewStringSlice, slice)
}

func (f *StringSliceFlag) SetDestination(slice []string) {
	f.Destination = newSliceFlagValue(NewStringSlice, slice)
}

func (f *StringSliceFlag) GetDestination() []string {
	if destination := f.Destination; destination != nil {
		return destination.Value()
	}
	return nil
}
