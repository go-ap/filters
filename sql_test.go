package filters

import "testing"

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
