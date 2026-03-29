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

func (a actorChecks) GoString() string {
	if len(a) == 0 {
		return ""
	}
	ss := strings.Builder{}
	ss.WriteString("actor={")
	for i, fn := range a {
		if sss, ok := fn.(fmt.GoStringer); ok {
			ss.WriteString(sss.GoString())
		}
		if i < len(a)-1 {
			ss.WriteRune(',')
		}
	}
	ss.WriteString("}")
	return ss.String()
}

func Target(fns ...Check) Check {
	if len(fns) == 0 {
		return nil
	}
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

func (t targetChecks) GoString() string {
	if len(t) == 0 {
		return ""
	}
	ss := strings.Builder{}
	ss.WriteString("target={")
	for i, fn := range t {
		if sss, ok := fn.(fmt.GoStringer); ok {
			ss.WriteString(sss.GoString())
		}
		if i < len(t)-1 {
			ss.WriteRune(',')
		}
	}
	ss.WriteString("}")
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

func (o objectChecks) GoString() string {
	if len(o) == 0 {
		return ""
	}
	ss := strings.Builder{}
	ss.WriteString("object={")
	for i, fn := range o {
		if sss, ok := fn.(fmt.GoStringer); ok {
			ss.WriteString(sss.GoString())
		}
		if i < len(o)-1 {
			ss.WriteRune(',')
		}
	}
	ss.WriteString("}")
	return ss.String()
}
