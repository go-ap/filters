package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func Test_withTypes_GoString(t *testing.T) {
	tests := []struct {
		name  string
		types withTypes
		want  string
	}{
		{
			name:  "empty",
			types: nil,
			want:  "type=[]",
		},
		{
			name:  "one type",
			types: withTypes{"t1"},
			want:  "type=[t1]",
		},
		{
			name:  "multiple types",
			types: withTypes{"t1", "t2"},
			want:  "type=[t1,t2]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.types.GoString(); got != tt.want {
				t.Errorf("GoString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withTypes_Match(t *testing.T) {
	tests := []struct {
		name  string
		types withTypes
		it    vocab.Item
		want  bool
	}{
		{
			name: "nil types match nil item",
			want: true,
		},
		{
			name:  "empty types match nil item",
			types: withTypes{},
			want:  true,
		},
		{
			name:  "empty types do not match item with type",
			types: withTypes{},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
			want:  false,
		},
		{
			name:  "empty types do not match item with multiple types",
			types: withTypes{},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1", "t2"}},
			want:  false,
		},
		{
			name:  "empty types do not match item with multiple types, including nil type",
			types: withTypes{},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1", "t2", vocab.NilType}},
			want:  false,
		},
		{
			name:  "empty types do not match item with type",
			types: withTypes{},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
			want:  false,
		},
		{
			name:  "t1 matches item type t1",
			types: withTypes{"t1"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
			want:  true,
		},
		{
			name:  "t1,t2 matches item type t1",
			types: withTypes{"t1", "t2"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
			want:  true,
		},
		{
			name:  "t1 matches item type {t1}",
			types: withTypes{"t1"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1"}},
			want:  true,
		},
		{
			name:  "t1 matches item type {t1, t2}",
			types: withTypes{"t1"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1", "t2"}},
			want:  true,
		},
		{
			name:  "t1,t2 matches item type {t1}",
			types: withTypes{"t1", "t2"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1"}},
			want:  true,
		},
		{
			name:  "t1,t2 matches item type {t1,t2}",
			types: withTypes{"t1", "t2"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1", "t2"}},
			want:  true,
		},
		{
			name:  "t1,t2 matches item type {t1,t2,t3}",
			types: withTypes{"t1", "t2"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1", "t2", "t3"}},
			want:  true,
		},
		{
			name:  "t1 matches item type {t1, t2}",
			types: withTypes{"t1"},
			it:    &vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1", "t2"}},
			want:  true,
		},
		{
			name:  "t1 matches item collection with item with type {t1}",
			types: withTypes{"t1"},
			it:    vocab.ItemCollection{&vocab.Object{Type: vocab.ActivityVocabularyTypes{"t1"}}},
			want:  false,
		},
		{
			name:  "t1 does not match item collection with item with type {t2}",
			types: withTypes{"t1"},
			it:    vocab.ItemCollection{&vocab.Object{Type: vocab.ActivityVocabularyTypes{"t2"}}},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.types.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_noType_Match(t *testing.T) {
	tests := []struct {
		name string
		it   vocab.Item
		want bool
	}{
		{
			name: "empty",
		},
		{
			name: "nil type does not match on item collection",
			it:   vocab.ItemCollection{},
			want: false,
		},
		{
			name: "nil type does not match on iris",
			it:   vocab.IRIs{},
			want: false,
		},
		{
			name: "nil type does not match on iri",
			it:   vocab.IRI("http://example.com"),
			want: false,
		},
		{
			name: "Create type doesn't match",
			it:   vocab.Activity{Type: vocab.CreateType},
			want: false,
		},
		{
			name: "not empty type slice does not match",
			it:   vocab.Object{Type: vocab.ActivityVocabularyTypes{vocab.NilType}},
			want: true,
		},
		{
			name: "multiple types don't match",
			it:   vocab.Activity{Type: vocab.ActivityVocabularyTypes{vocab.UpdateType, "CustomUpdate"}},
			want: false,
		},
		{
			name: "empty type matches",
			it:   vocab.Object{Type: vocab.NilType},
			want: true,
		},
		{
			name: "empty type slice matches",
			it:   vocab.Object{Type: vocab.ActivityVocabularyTypes{}},
			want: true,
		},
		{
			name: "nil type matches",
			it:   vocab.Object{},
			want: true,
		},
		{
			name: "nil type matches on ptr",
			it:   &vocab.Object{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := noType{}
			if got := n.Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
