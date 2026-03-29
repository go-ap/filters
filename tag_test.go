package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func Test_tagChecks_Match(t *testing.T) {
	tests := []struct {
		name string
		a    tagChecks
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name: "item does not match empty tag filter",
			a:    tagChecks{},
			it:   &vocab.Object{ID: "http://example.com"},
			want: true,
		},
		{
			name: "item matches empty nil tag filter",
			a:    tagChecks{NilItem},
			it:   &vocab.Object{ID: "http://example.com"},
			want: true,
		},
		{
			name: "nil item does not match empty tag filter",
			a:    tagChecks{},
			want: false,
		},
		{
			name: "nil item does not match not empty tag filter",
			a:    tagChecks{NameIs("#tag")},
			want: false,
		},
		{
			name: "item match tag IRI filter",
			a:    tagChecks{SameIRI("http://example.com/tag")},
			it:   &vocab.Object{ID: "http://example.com", Tag: vocab.ItemCollection{vocab.IRI("http://example.com/tag")}},
			want: true,
		},
		{
			name: "item matches name tag filter",
			a:    tagChecks{NameIs("#tag")},
			it:   &vocab.Object{ID: "http://example.com", Tag: vocab.ItemCollection{&vocab.Object{Name: vocab.DefaultNaturalLanguage("#tag")}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
