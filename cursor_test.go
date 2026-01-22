package filters

import (
	"reflect"
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func TestPaginationChecks(t *testing.T) {
	tests := []struct {
		name string
		fns  []Check
		item vocab.Item
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "just after",
			fns:  Checks{After(SameID("https://example.com"))},
			item: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/2"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
		{
			name: "after with filters",
			fns:  Checks{After(SameID("example.com"), HasType("Activity"), After(SameID("example.com")))},
			want: nil,
		},
		{
			name: "maxItems=2 of 2",
			fns:  Checks{WithMaxCount(2)},
			item: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
		{
			name: "maxItems=2 of 3",
			fns:  Checks{WithMaxCount(2)},
			item: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
				vocab.Activity{ID: "https://example.com/3"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
		{
			name: "before=https://example.com/1 single item",
			fns:  Checks{Before(SameID("https://example.com/1"))},
			item: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "before=https://example.com/1 second item",
			fns:  Checks{Before(SameID("https://example.com/1"))},
			item: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/0"},
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/0"},
			},
		},
		{
			name: "after=https://example.com/1 first item",
			fns:  Checks{After(SameID("https://example.com/1"))},
			item: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "after=https://example.com/1 second item",
			fns:  Checks{After(SameID("https://example.com/1"))},
			item: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/0"},
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChecks := PaginationChecks(tt.fns...)
			if got := gotChecks.Run(tt.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CursorChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextPageFromCollection(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.CollectionInterface
		want vocab.IRI
	}{
		{
			name: "empty",
			arg:  nil,
			want: vocab.EmptyIRI,
		},
		{
			name: "Collection with empty First",
			arg:  &vocab.Collection{Type: vocab.CollectionType},
			want: vocab.EmptyIRI,
		},
		{
			name: "Collection with First",
			arg: &vocab.Collection{
				Type:  vocab.CollectionType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		{
			name: "OrderedCollection with empty First",
			arg:  &vocab.OrderedCollection{Type: vocab.OrderedCollectionType},
			want: vocab.EmptyIRI,
		},
		{
			name: "OrderedCollection with First",
			arg: &vocab.OrderedCollection{
				Type:  vocab.OrderedCollectionType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		// CollectionPages

		{
			name: "CollectionPage with empty First",
			arg:  &vocab.CollectionPage{Type: vocab.CollectionPageType},
			want: vocab.EmptyIRI,
		},
		{
			name: "CollectionPage with First",
			arg: &vocab.CollectionPage{
				Type:  vocab.CollectionPageType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		{
			name: "CollectionPage with First and Next",
			arg: &vocab.CollectionPage{
				Type:  vocab.CollectionPageType,
				First: vocab.IRI("https://example.com?first"),
				Next:  vocab.IRI("https://example.com?next"),
			},
			want: "https://example.com?next",
		},
		{
			name: "OrderedCollectionPage with empty First",
			arg:  &vocab.OrderedCollectionPage{Type: vocab.OrderedCollectionPageType},
			want: vocab.EmptyIRI,
		},
		{
			name: "OrderedCollectionPage with First",
			arg: &vocab.OrderedCollectionPage{
				Type:  vocab.OrderedCollectionPageType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		{
			name: "OrderedCollectionPage with First and Next",
			arg: &vocab.OrderedCollectionPage{
				Type:  vocab.OrderedCollectionPageType,
				First: vocab.IRI("https://example.com?first"),
				Next:  vocab.IRI("https://example.com?next"),
			},
			want: "https://example.com?next",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextPageFromCollection(tt.arg); got != tt.want {
				t.Errorf("NextPageFromCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrevPageFromCollection(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.CollectionInterface
		want vocab.IRI
	}{
		{
			name: "empty",
			arg:  nil,
			want: vocab.EmptyIRI,
		},
		{
			name: "Collection with empty First",
			arg:  &vocab.Collection{Type: vocab.CollectionType},
			want: vocab.EmptyIRI,
		},
		{
			name: "Collection with First",
			arg: &vocab.Collection{
				Type:  vocab.CollectionType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		{
			name: "OrderedCollection with empty First",
			arg:  &vocab.OrderedCollection{Type: vocab.OrderedCollectionType},
			want: vocab.EmptyIRI,
		},
		{
			name: "OrderedCollection with First",
			arg: &vocab.OrderedCollection{
				Type:  vocab.OrderedCollectionType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		// CollectionPages

		{
			name: "CollectionPage with empty First",
			arg:  &vocab.CollectionPage{Type: vocab.CollectionPageType},
			want: vocab.EmptyIRI,
		},
		{
			name: "CollectionPage with First",
			arg: &vocab.CollectionPage{
				Type:  vocab.CollectionPageType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		{
			name: "CollectionPage with First and Prev",
			arg: &vocab.CollectionPage{
				Type:  vocab.CollectionPageType,
				First: vocab.IRI("https://example.com?first"),
				Prev:  vocab.IRI("https://example.com?previous"),
			},
			want: "https://example.com?previous",
		},
		{
			name: "OrderedCollectionPage with empty First",
			arg:  &vocab.OrderedCollectionPage{Type: vocab.OrderedCollectionPageType},
			want: vocab.EmptyIRI,
		},
		{
			name: "OrderedCollectionPage with First",
			arg: &vocab.OrderedCollectionPage{
				Type:  vocab.OrderedCollectionPageType,
				First: vocab.IRI("https://example.com?first"),
			},
			want: "https://example.com?first",
		},
		{
			name: "OrderedCollectionPage with First and Prev",
			arg: &vocab.OrderedCollectionPage{
				Type:  vocab.OrderedCollectionPageType,
				First: vocab.IRI("https://example.com?first"),
				Prev:  vocab.IRI("https://example.com?previous"),
			},
			want: "https://example.com?previous",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrevPageFromCollection(tt.arg); got != tt.want {
				t.Errorf("PrevPageFromCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
