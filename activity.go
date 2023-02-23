package filters

import vocab "github.com/go-ap/activitypub"

func Actor(fns ...Fn) Fn {
	return func(it vocab.Item) bool {
		act, err := vocab.ToIntransitiveActivity(it)
		if err != nil {
			return false
		}
		return All(fns...)(act.Actor)
	}
}

func Target(fns ...Fn) Fn {
	return func(it vocab.Item) bool {
		act, err := vocab.ToIntransitiveActivity(it)
		if err != nil {
			return false
		}
		return All(fns...)(act.Target)
	}
}

func Object(fns ...Fn) Fn {
	return func(it vocab.Item) bool {
		act, err := vocab.ToActivity(it)
		if err != nil {
			return false
		}
		return All(fns...)(act.Object)
	}
}
