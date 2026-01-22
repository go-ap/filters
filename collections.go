package filters

import (
	"fmt"
	"strconv"
	"strings"

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

func (cnt *counter) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	if cnt.max <= cnt.cnt {
		return false
	}
	cnt.cnt = cnt.cnt + 1
	return true
}

func (cnt *counter) String() string {
	return "maxItems=" + strconv.Itoa(cnt.max)
}

// After checks the activitypub.Item against a specified "fn" filter function.
// This should be used when iterating over a collection, and it resolves to true
// after fn returns true and to false check.
//
// Due to relying on the static check function return value the After is not reentrant.
func After(fns ...Check) Check {
	return &afterCrit{check: false, fns: fns}
}

func (a afterCrit) String() string {
	ss := strings.Builder{}
	ss.WriteString("after.")
	for i, fn := range a.fns {
		if sss, ok := fn.(fmt.Stringer); ok {
			ss.WriteString(sss.String())
		}
		if i < len(a.fns)-1 {
			ss.WriteRune(',')
		}
	}
	return ss.String()
}

func (isAfter *afterCrit) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return isAfter.check
	}

	if checkFn(isAfter.fns)(it) {
		isAfter.check = true
		return false
	}
	return isAfter.check
}

type afterCrit struct {
	check bool
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

func (isBefore *beforeCrit) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	if checkFn(isBefore.fns)(it) {
		isBefore.check = false
	}
	return isBefore.check
}
