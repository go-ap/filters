package filters

import vocab "github.com/go-ap/activitypub"

func Actor(fns ...Check) Check {
	return actorChecks(fns)
}

type actorChecks []Check

func (a actorChecks) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	act, err := vocab.ToIntransitiveActivity(it)
	if err != nil {
		return false
	}
	return All(a...).Match(act.Actor)
}

func Target(fns ...Check) Check {
	return targetChecks(fns)
}

type targetChecks []Check

func (t targetChecks) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	act, err := vocab.ToIntransitiveActivity(it)
	if err != nil {
		return false
	}
	return All(t...).Match(act.Target)
}

func Object(fns ...Check) Check {
	return objectChecks(fns)
}

type objectChecks []Check

func (o objectChecks) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	act, err := vocab.ToActivity(it)
	if err != nil {
		return false
	}
	return All(o...).Match(act.Object)
}
