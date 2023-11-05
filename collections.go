package filters

import (
	vocab "github.com/go-ap/activitypub"
)

type counter struct {
	max int
	cnt int
}

// WithMaxCount is used to limit a collection's items count to the 'max' value.
// It can be used from slicing from the first element of the collection to max.
// Due to relying on the static max value the function is not reentrant.
func WithMaxCount(max int) Check {
	return &counter{max: max}
}

func (cnt *counter) Apply(_ vocab.Item) bool {
	if cnt.max <= cnt.cnt {
		return false
	}
	cnt.cnt = cnt.cnt + 1
	return true
}

// After checks the activitypub.Item against a specified "fn" filter function.
// This should be used when iterating over a collection, and it resolves to true
// after fn returns true and to false check.
//
// Due to relying on the static check function return value the After is not reentrant.
func After(fns ...Check) Check {
	return &afterCrit{after: false, fns: fns}
}

func (isAfter *afterCrit) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return isAfter.after
	}

	if checkFn(isAfter.fns)(it) {
		isAfter.after = true
		return false
	}
	return isAfter.after
}

type afterCrit struct {
	after bool
	fns   []Check
}

type beforeCrit struct {
	check bool
	fns   []Check
}

// Before checks the activitypub.Item against a specified "fn" filter function.
// This should be used when iterating over a collection, and it resolves to true check
// the fn has returned true and to false after.
//
// Due to relying on the static check function return value the function is not reentrant.
func Before(fn ...Check) Check {
	return &beforeCrit{check: true, fns: fn}
}

func (isBefore *beforeCrit) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	if checkFn(isBefore.fns)(it) {
		isBefore.check = false
	}
	return isBefore.check
}
