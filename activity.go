package filters

import (
	"fmt"
	"strings"

	vocab "github.com/go-ap/activitypub"
)

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

func (a actorChecks) String() string {
	ss := strings.Builder{}
	ss.WriteString("actor.")
	for i, fn := range a {
		if sss, ok := fn.(fmt.Stringer); ok {
			ss.WriteString(sss.String())
		}
		if i < len(a)-1 {
			ss.WriteRune(',')
		}
	}
	return ss.String()
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

func (t targetChecks) String() string {
	ss := strings.Builder{}
	ss.WriteString("target.")
	for i, fn := range t {
		if sss, ok := fn.(fmt.Stringer); ok {
			ss.WriteString(sss.String())
		}
		if i < len(t)-1 {
			ss.WriteRune(',')
		}
	}
	return ss.String()
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

func (o objectChecks) String() string {
	ss := strings.Builder{}
	ss.WriteString("object.")
	for i, fn := range o {
		if sss, ok := fn.(fmt.Stringer); ok {
			ss.WriteString(sss.String())
		}
		if i < len(o)-1 {
			ss.WriteRune(',')
		}
	}
	return ss.String()
}
