package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func Test_Recipients_Match(t *testing.T) {
	tests := []struct {
		name string
		a    vocab.IRI
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
		},
		{
			name: "example.com negative",
			a:    "https://example.com",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: false,
		},
		{
			name: "example.com negative with public To",
			a:    "https://example.com",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.PublicNS}},
			want: false,
		},
		{
			name: "example.com in To",
			a:    "https://example.com",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "PublicNS negative for empty recipients",
			a:    vocab.PublicNS,
			it:   &vocab.Object{Type: vocab.TombstoneType},
			want: false,
		},
		{
			name: "example.com in CC",
			a:    "https://example.com",
			it:   &vocab.Object{CC: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "example.com in Bto",
			a:    "https://example.com",
			it:   &vocab.Object{Bto: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
		{
			name: "example.com in BCC",
			a:    "https://example.com",
			it:   &vocab.Object{BCC: vocab.ItemCollection{vocab.IRI("https://example.com")}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Recipients(tt.a).Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
