package filters

import vocab "github.com/go-ap/activitypub"

func Tag(fns ...Check) Check {
	return tagChecks(fns)
}

type tagChecks []Check

func (a tagChecks) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	ob, err := vocab.ToObject(it)
	if err != nil {
		return false
	}
	return All(a...).Match(ob.Tag)
}
