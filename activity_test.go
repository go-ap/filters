package filters

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_objectChecks_GoString(t *testing.T) {
	tests := []struct {
		name string
		o    objectChecks
		want string
	}{
		{
			name: "empty",
			o:    nil,
			want: "",
		},
		{
			name: "same id",
			o:    objectChecks{SameID("http://example.com")},
			want: "object={id=http://example.com}",
		},
		{
			name: "same id and type t1",
			o:    objectChecks{SameID("http://example.com"), HasType("t1")},
			want: "object={id=http://example.com,type=[t1]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GoString(); got != tt.want {
				t.Errorf("GoString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_targetChecks_GoString(t *testing.T) {
	tests := []struct {
		name string
		t    targetChecks
		want string
	}{
		{
			name: "empty",
			t:    nil,
			want: "",
		},
		{
			name: "same id",
			t:    targetChecks{SameID("http://example.com")},
			want: "target={id=http://example.com}",
		},
		{
			name: "same id and type t1",
			t:    targetChecks{SameID("http://example.com"), HasType("t1")},
			want: "target={id=http://example.com,type=[t1]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.GoString(); got != tt.want {
				t.Errorf("GoString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_actorChecks_GoString(t *testing.T) {
	tests := []struct {
		name string
		a    actorChecks
		want string
	}{
		{
			name: "empty",
			a:    nil,
			want: "",
		},
		{
			name: "same id",
			a:    actorChecks{SameID("http://example.com")},
			want: "actor={id=http://example.com}",
		},
		{
			name: "same id and type t1",
			a:    actorChecks{SameID("http://example.com"), HasType("t1")},
			want: "actor={id=http://example.com,type=[t1]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.GoString(); got != tt.want {
				t.Errorf("GoString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTarget(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Check
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: []Check{},
			want: nil,
		},
		{
			name: "empty",
			args: []Check{SameID("http://example.com")},
			want: targetChecks{SameID("http://example.com")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Target(tt.args...); !cmp.Equal(got, tt.want) {
				t.Errorf("Target() = %v, want %v", got, tt.want)
			}
		})
	}
}
