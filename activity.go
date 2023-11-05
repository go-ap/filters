package filters

import vocab "github.com/go-ap/activitypub"

type actorCrit []Check

func (a actorCrit) Apply(it vocab.Item) bool {
	act, err := vocab.ToIntransitiveActivity(it)
	if err != nil {
		return false
	}
	return All(a...).Apply(act.Actor)
}

func Actor(fns ...Check) Check {
	return actorCrit(fns)
}

type targetCrit []Check

func (t targetCrit) Apply(it vocab.Item) bool {
	act, err := vocab.ToIntransitiveActivity(it)
	if err != nil {
		return false
	}
	return All(t...).Apply(act.Target)
}
func Target(fns ...Check) Check {
	return targetCrit(fns)
}

type objectCrit []Check

func (o objectCrit) Apply(it vocab.Item) bool {
	act, err := vocab.ToActivity(it)
	if err != nil {
		return false
	}
	return All(o...).Apply(act.Object)
}

func Object(fns ...Check) Check {
	return objectCrit(fns)
}
