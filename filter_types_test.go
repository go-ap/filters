package filters

import (
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxCount(tt.fns...); got != tt.want {
				t.Errorf("MaxCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
