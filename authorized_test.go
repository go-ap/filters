package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func Test_authorized_Apply(t *testing.T) {
	tests := []struct {
		name string
		a    vocab.IRI
		it   vocab.Item
		want bool
	}{
		{
			name: "empty is not authorized",
		},
		{
			name: "example.com is not allowed",
			a:    "https://example.com",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: false,
		},
		{
			name: "example.com is allowed if object has public audience",
			a:    "https://example.com",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.PublicNS}},
			want: true,
		},
		{
			name: "PublicNS should be authorized for object with empty recipients list",
			a:    vocab.PublicNS,
			it:   &vocab.Object{Type: vocab.TombstoneType},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := authorized(tt.a).Apply(tt.it); got != tt.want {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
