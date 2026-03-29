package filters

import (
	"reflect"
	"testing"

	vocab "github.com/go-ap/activitypub"
	"github.com/google/go-cmp/cmp"
)

func TestMaxCount(t *testing.T) {
	tests := []struct {
		name string
		fns  []Check
		want int
	}{
		{
			name: "empty",
			fns:  nil,
			want: -1,
		},
		{
			name: "one check max 100",
			fns: Checks{
				&counter{max: 100},
			},
			want: 100,
		},
		{
			name: "multiple checks max 666",
			fns: Checks{
				All(HasType(vocab.PersonType)),
				&counter{max: 666},
			},
			want: 666,
		},
		{
			name: "all check with max 666",
			fns: Checks{
				All(&counter{max: 666}),
			},
			want: 666,
		},
		{
			name: "all checks with max 665 and additional filter",
			fns: Checks{
				All(HasType(vocab.PersonType), &counter{max: 665}),
			},
			want: 665,
		},
		{
			name: "random all checks, and max 5 and additional filter",
			fns: Checks{
				All(HasType(vocab.PersonType)),
				WithMaxCount(5),
				SameIRI("https://example.com"),
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxCount(tt.fns...); got != tt.want {
				t.Errorf("MaxCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterChecks(t *testing.T) {
	tests := []struct {
		name string
		args Checks
		want Checks
	}{
		{
			name: "empty",
		},
		{
			name: "empty when passing only a maxCount",
			args: Checks{WithMaxCount(5)},
			want: Checks{},
		},
		{
			name: "remove a maxCount",
			args: Checks{
				All(HasType(vocab.PersonType)),
				WithMaxCount(5),
				SameIRI("https://example.com"),
			},
			want: Checks{
				All(HasType(vocab.PersonType)),
				SameIRI("https://example.com"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterChecks(tt.args...); !cmp.Equal(got, tt.want) {
				t.Errorf("FilterChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSameIDChecks(t *testing.T) {
	tests := []struct {
		name   string
		checks []Check
		want   Checks
	}{
		{
			name: "empty",
			want: Checks{},
		},
		{
			name:   "has SameID",
			checks: Checks{SameID("http://example.com")},
			want:   Checks{SameID("http://example.com")},
		},
		{
			name:   "does not have SameID",
			checks: Checks{NilItem, NilID},
			want:   Checks{},
		},
		{
			name:   "All check with SameID",
			checks: Checks{All(SameID("http://example.com"), NilItem)},
			want:   Checks{All(SameID("http://example.com"))},
		},
		{
			name:   "Any check with SameID",
			checks: Checks{Any(NilInReplyTo, SameID("http://example.com/1"))},
			want:   Checks{Any(SameID("http://example.com/1"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameIDChecks(tt.checks...); !cmp.Equal(got, tt.want) {
				t.Errorf("SameIDChecks() = %s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestIDChecks(t *testing.T) {
	tests := []struct {
		name   string
		checks []Check
		want   Checks
	}{
		{
			name: "empty",
			want: Checks{},
		},
		{
			name:   "has SameID",
			checks: Checks{SameID("http://example.com")},
			want:   Checks{SameID("http://example.com")},
		},
		{
			name:   "has NilID",
			checks: Checks{NilID},
			want:   Checks{NilID},
		},
		{
			name:   "has IDLike",
			checks: Checks{IDLike("https://example.com/")},
			want:   Checks{IDLike("https://example.com/")},
		},
		{
			name:   "has multiple ID Checks",
			checks: Checks{SameID("http://example.com"), IDLike("https://example.com/"), NilID},
			want:   Checks{SameID("http://example.com"), IDLike("https://example.com/"), NilID},
		},
		{
			name:   "does not have SameID",
			checks: Checks{NilItem, SameAttributedTo("https://example.com/~jdoe")},
			want:   Checks{},
		},
		{
			name:   "All check with SameID",
			checks: Checks{All(SameID("http://example.com"), NilItem, NilID, IDLike("https://"))},
			want:   Checks{All(SameID("http://example.com"), NilID, IDLike("https://"))},
		},
		{
			name:   "Any check with SameID",
			checks: Checks{Any(NameIs("example"), IDLike("http://"), NilID, NilInReplyTo, SameID("http://example.com/1"))},
			want:   Checks{Any(IDLike("http://"), NilID, SameID("http://example.com/1"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IDChecks(tt.checks...); !cmp.Equal(got, tt.want) {
				t.Errorf("IDChecks() = %s", cmp.Diff(tt.want, got))
			}
		})
	}
}

//func Test___Checks(t *testing.T) {
//	tests := []struct {
//		name string
//		checks Checks
//		want Checks
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := ItemChecks(tt.checks...); !cmp.Equal(got, tt.want) {
//				t.Errorf("ItemChecks() = %s", cmp.Diff( tt.want, got))
//			}
//		})
//	}
//}

func TestItemChecks(t *testing.T) {
	tests := []struct {
		name string
		args Checks
		want Checks
	}{
		{
			name: "nil",
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "just activity checks",
			args: Checks{Actor(SameID("http://example.com/~jdoe")), Object(SameID("http://example.com/1")), Target(SameID("http://example.com/tgt"))},
			want: Checks{},
		},
		{
			name: "just tags checks",
			args: Checks{Tag(NameIs("Test"))},
			want: Checks{},
		},
		{
			name: "same id check",
			args: Checks{SameID("http://example.com")},
			want: Checks{SameID("http://example.com")},
		},
		{
			name: "type check",
			args: Checks{HasType("t1")},
			want: Checks{HasType("t1")},
		},
		{
			name: "any types has one check",
			args: Checks{Any(HasType("t1"), Actor(NilItem))},
			want: Checks{HasType("t1")},
		},
		{
			name: "all types has one check",
			args: Checks{All(HasType("t1"), Actor(NilItem))},
			want: Checks{HasType("t1")},
		},
		{
			name: "any types has multiple checks",
			args: Checks{Any(HasType("t1"), Actor(NilItem), SameID("http://example.com"))},
			want: Checks{Any(HasType("t1"), SameID("http://example.com"))},
		},
		{
			name: "all types has multiple checks",
			args: Checks{All(HasType("t1"), Actor(NilItem), SameID("http://example.com"))},
			want: Checks{All(HasType("t1"), SameID("http://example.com"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ItemChecks(tt.args...); !cmp.Equal(got, tt.want) {
				t.Errorf("ItemChecks() = %s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestCursorChecks(t *testing.T) {
	tests := []struct {
		name   string
		checks Checks
		want   Checks
	}{
		{
			name: "empty",
			want: Checks{},
		},
		{
			name:   "has After",
			checks: Checks{After(NilID), NilItem},
			want:   Checks{After(NilID)},
		},
		{
			name:   "has Before",
			checks: Checks{SummaryEmpty, Before(SameID("http://example.com")), NilItem},
			want:   Checks{Before(SameID("http://example.com"))},
		},
		{
			name:   "has both",
			checks: Checks{After(SameID("http://example.com/1")), SummaryEmpty, Before(SameID("http://example.com/666")), NilItem},
			want:   Checks{After(SameID("http://example.com/1")), Before(SameID("http://example.com/666"))},
		},
		{
			name:   "has none",
			checks: Checks{SummaryEmpty, NilItem, WithMaxCount(100)},
			want:   Checks{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CursorChecks(tt.checks...); !cmp.Equal(got, tt.want, cmp.Comparer(checksEq)) {
				t.Errorf("CursorChecks() = %s", cmp.Diff(tt.want, got, cmp.Comparer(checksEq)))
			}
		})
	}
}

func checksEq(c1, c2 Check) bool {
	u1 := urlValue(c1)
	u2 := urlValue(c2)
	return u1.Encode() == u2.Encode()
}

func TestCounted(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want int
	}{
		{
			name: "nil",
			want: -1,
		},
		{
			name: "no hits",
			args: Checks{&counter{max: 10, cnt: 0}},
			want: 0,
		},
		{
			name: "1 hit",
			args: Checks{&counter{max: 10, cnt: 1}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Counted(tt.args...); got != tt.want {
				t.Errorf("Counted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no type checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "with type check",
			args: Checks{SameID("http://example.com"), HasType("t1")},
			want: Checks{HasType("t1")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeChecks(tt.args...); !cmp.Equal(got, tt.want) {
				t.Errorf("TypeChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no tag checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "with tag check",
			args: Checks{SameID("http://example.com"), Tag(NilItem)},
			want: Checks{Tag(NilItem)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TagChecks(tt.args...); !cmp.Equal(got, tt.want) {
				t.Errorf("TagChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntransitiveActivityChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no intransitive checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "with actor check",
			args: Checks{SameID("http://example.com"), Actor(NilItem)},
			want: Checks{NilItem},
		},
		{
			name: "with target check",
			args: Checks{SameID("http://example.com"), Target(NameIs("test"))},
			want: Checks{NameIs("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntransitiveActivityChecks(tt.args...); !cmp.Equal(got, tt.want, cmp.Comparer(NaturalLanguageValuesComparer)) {
				t.Errorf("IntransitiveActivityChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActivityChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no  checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "with actor check",
			args: Checks{SameID("http://example.com"), Actor(NilItem)},
			want: Checks{NilItem},
		},
		{
			name: "with multiple actor checks",
			args: Checks{SameID("http://example.com"), Actor(NilItem, SameID("http://example.com/~jdoe"))},
			want: Checks{NilItem, SameID("http://example.com/~jdoe")},
		},
		{
			name: "with target check",
			args: Checks{SameID("http://example.com"), Target(NameIs("test"))},
			want: Checks{NameIs("test")},
		},
		{
			name: "with multiple target check",
			args: Checks{SameID("http://example.com"), Target(NameIs("test"), SummaryLike("lipsum"))},
			want: Checks{NameIs("test"), SummaryLike("lipsum")},
		},
		{
			name: "with object check",
			args: Checks{SameID("http://example.com"), Object(NilItem)},
			want: Checks{NilItem},
		},
		{
			name: "with multiple object check",
			args: Checks{SameID("http://example.com"), Object(HasType("t1"), SameInReplyTo("http://example.com/~jdoe"))},
			want: Checks{HasType("t1"), SameInReplyTo("http://example.com/~jdoe")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ActivityChecks(tt.args...); !cmp.Equal(got, tt.want, cmp.Comparer(NaturalLanguageValuesComparer)) {
				t.Errorf("ActivityChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no object checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "with object check",
			args: Checks{SameID("http://example.com"), Object(NilItem)},
			want: Checks{NilItem},
		},
		{
			name: "with multiple object check",
			args: Checks{SameID("http://example.com"), Object(HasType("t1"), SameInReplyTo("http://example.com/~jdoe"))},
			want: Checks{HasType("t1"), SameInReplyTo("http://example.com/~jdoe")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ObjectChecks(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActorChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no actor checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "with actor check",
			args: Checks{SameID("http://example.com"), Actor(NilItem)},
			want: Checks{NilItem},
		},
		{
			name: "with multiple actor check",
			args: Checks{SameID("http://example.com"), Actor(HasType("t1"), SameInReplyTo("http://example.com/~jdoe"))},
			want: Checks{HasType("t1"), SameInReplyTo("http://example.com/~jdoe")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ActorChecks(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ActorChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTargetChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no target checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "with target check",
			args: Checks{SameID("http://example.com"), Target(NilItem)},
			want: Checks{NilItem},
		},
		{
			name: "with multiple target check",
			args: Checks{SameID("http://example.com"), Target(HasType("t1"), SameInReplyTo("http://example.com/~jdoe"))},
			want: Checks{HasType("t1"), SameInReplyTo("http://example.com/~jdoe")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TargetChecks(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TargetChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}
