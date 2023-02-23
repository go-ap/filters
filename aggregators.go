package filters

import vocab "github.com/go-ap/activitypub"

func Any(fns ...Fn) Fn {
	return func(it vocab.Item) bool {
		for _, fn := range fns {
			if fn(it) {
				return true
			}
		}
		return false
	}
}

func All(fns ...Fn) Fn {
	return func(it vocab.Item) bool {
		for _, fn := range fns {
			if !fn(it) {
				return false
			}
		}
		return true
	}
}

func NameIn(names ...string) Fn {
	ors := make([]Fn, 0)
	for _, n := range names {
		ors = append(ors, NameIs(n))
	}
	return Any(ors...)
}
