package filters

import vocab "github.com/go-ap/activitypub"

type actorChecks []Check

func (a actorChecks) Apply(it vocab.Item) bool {
	act, err := vocab.ToIntransitiveActivity(it)
	if err != nil {
		return false
	}
	return All(a...).Apply(act.Actor)
}

func Actor(fns ...Check) Check {
	return actorChecks(fns)
}

type targetChecks []Check

func (t targetChecks) Apply(it vocab.Item) bool {
	act, err := vocab.ToIntransitiveActivity(it)
	if err != nil {
		return false
	}
	return All(t...).Apply(act.Target)
}
func Target(fns ...Check) Check {
	return targetChecks(fns)
}

type objectChecks []Check

func (o objectChecks) Apply(it vocab.Item) bool {
	act, err := vocab.ToActivity(it)
	if err != nil {
		return false
	}
	return All(o...).Apply(act.Object)
}

func Object(fns ...Check) Check {
	return objectChecks(fns)
}
