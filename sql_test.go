package filters

import (
	"reflect"
	"testing"
)

func TestGetLimit(t *testing.T) {
	tests := []struct {
		name string
		f    []Check
		want int
	}{
		{
			name: "empty",
			f:    nil,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLimit(tt.f...); got != tt.want {
				t.Errorf("GetLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWhereClauses(t *testing.T) {
	tests := []struct {
		name  string
		args  []Check
		want  []string
		want1 []any
	}{
		{
			name:  "empty",
			args:  nil,
			want:  nil,
			want1: nil,
		},
		{
			name:  "not empty",
			args:  Checks{HasType("test", "test1")},
			want:  []string{"type IN (?,?)"},
			want1: []any{"test", "test1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetWhereClauses(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWhereClauses() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetWhereClauses() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
