package filters

import (
	"reflect"
	"testing"

	vocab "github.com/go-ap/activitypub"
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
			if got := FilterChecks(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}
