package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func _mockTrue(_ vocab.Item) bool {
	return true
}

func _mockFalse(_ vocab.Item) bool {
	return false
}

func TestAny(t *testing.T) {
	tests := []struct {
		name string
		fns  []Fn
		want bool
	}{
		{
			name: "empty is false",
			want: false,
		},
		{
			name: "one true, one false",
			fns:  []Fn{_mockTrue, _mockFalse},
			want: true,
		},
		{
			name: "all true",
			fns:  []Fn{_mockTrue, _mockTrue},
			want: true,
		},
		{
			name: "all false",
			fns:  []Fn{_mockFalse, _mockFalse},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{}
			if got := Any(tt.fns...)(ob); got != tt.want {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name string
		fns  []Fn
		want bool
	}{
		{
			name: "empty is false",
			want: false,
		},
		{
			name: "one true, one false",
			fns:  []Fn{_mockTrue, _mockFalse},
			want: false,
		},
		{
			name: "all true",
			fns:  []Fn{_mockTrue, _mockTrue},
			want: true,
		},
		{
			name: "all false",
			fns:  []Fn{_mockFalse, _mockFalse},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{}
			if got := All(tt.fns...)(ob); got != tt.want {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}
