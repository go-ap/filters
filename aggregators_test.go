package filters

import (
	"reflect"
	"testing"
)

func TestAny(t *testing.T) {
	type args struct {
		fns []Fn
	}
	tests := []struct {
		name string
		args args
		want Fn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Any(tt.args.fns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}
