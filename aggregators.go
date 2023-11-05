package filters

import vocab "github.com/go-ap/activitypub"

type notCrit []Check

func (n notCrit) Apply(it vocab.Item) bool {
	if len(n) == 0 {
		return false
	}
	f := n[0]
	if f == nil {
		return false
	}
	return !f.Apply(it)
}

func Not(fn Check) Check {
	return notCrit([]Check{fn})
}

type anyCrit []Check

func (a anyCrit) Apply(it vocab.Item) bool {
	for _, fn := range a {
		if fn == nil {
			continue
		}
		if fn.Apply(it) {
			return true
		}
	}
	return false
}

func Any(fns ...Check) Check {
	return anyCrit(fns)
}

type checkAll []Check

func (a checkAll) Apply(it vocab.Item) bool {
	if len(a) == 0 {
		return true
	}
	for _, fn := range a {
		if fn == nil {
			continue
		}
		if !fn.Apply(it) {
			return false
		}
	}
	return true
}
func All(fns ...Check) Check {
	return checkAll(fns)
}
