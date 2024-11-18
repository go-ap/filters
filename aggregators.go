package filters

import vocab "github.com/go-ap/activitypub"

type notCrit []Check

func (n notCrit) Match(it vocab.Item) bool {
	if len(n) == 0 {
		return false
	}
	f := n[0]
	if f == nil {
		return false
	}
	return !f.Match(it)
}

// Not negates the result of a Check function.
// It is equivalent to a unary NOT operator.
func Not(fn Check) Check {
	return notCrit([]Check{fn})
}

type checkAny []Check

func (a checkAny) Match(it vocab.Item) bool {
	for _, fn := range a {
		if fn == nil {
			continue
		}
		if fn.Match(it) {
			return true
		}
	}
	return false
}

// Any aggregates a list of individual Check functions into a single Check
// which resolves to false if all the individual members resolve as false,
// and true if any of them resolves as true.
// It is equivalent to a sequence of OR operators.
func Any(fns ...Check) Check {
	if len(fns) == 1 {
		return fns[0]
	}
	return checkAny(fns)
}

type checkAll []Check

func (a checkAll) Match(it vocab.Item) bool {
	if len(a) == 0 {
		return true
	}
	for _, fn := range a {
		if fn == nil {
			continue
		}
		if !fn.Match(it) {
			return false
		}
	}
	return true
}

// All aggregates a list of individual Check functions into a single Check
// which resolves true if all individual members resolve as true, and false otherwise.
// It is equivalent to a sequence of AND operators.
func All(fns ...Check) Check {
	if len(fns) == 1 {
		return fns[0]
	}
	return checkAll(fns)
}
