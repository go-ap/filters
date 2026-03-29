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
	if len(a) == 0 {
		a = append(a, NilItem)
	}
	ob, err := vocab.ToObject(it)
	if err != nil {
		return false
	}
	// NOTE(marius): The tag property is likely to be an item collection
	// so we match if any of the items matches.
	match := All(a...).Match(ob.Tag)
	if match {
		return match
	}
	_ = vocab.OnItem(ob.Tag, func(item vocab.Item) error {
		match = match || All(a...).Match(item)
		return nil
	})
	return match
}
