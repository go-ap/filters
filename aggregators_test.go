package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func TestAny(t *testing.T) {
	type args struct {
		fns []Fn
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := vocab.Object{}
			if got := Any(tt.args.fns...)(ob); got != tt.want {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}
