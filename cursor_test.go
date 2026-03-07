package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
	"github.com/google/go-cmp/cmp"
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
		{
			name: "no pagination",
			fns:  Checks{HasType("Note", "Article", "Document")},
			item: vocab.ItemCollection{
				vocab.Object{ID: "https://example.com/0", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/1", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/2", Type: vocab.DocumentType},
				vocab.Object{ID: "https://example.com/3", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/4", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/5", Type: vocab.DocumentType},
				vocab.Object{ID: "https://example.com/6", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/7", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/8", Type: vocab.DocumentType},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "https://example.com/0", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/1", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/2", Type: vocab.DocumentType},
				vocab.Object{ID: "https://example.com/3", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/4", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/5", Type: vocab.DocumentType},
				vocab.Object{ID: "https://example.com/6", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/7", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/8", Type: vocab.DocumentType},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChecks := PaginationChecks(tt.fns...)
			if got := gotChecks.Run(tt.item); !cmp.Equal(got, tt.want) {
				t.Errorf("CursorChecks() = %s", cmp.Diff(tt.want, got))
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

func TestPaginateCollection(t *testing.T) {
	type args struct {
		it      vocab.Item
		filters []Check
	}
	tests := []struct {
		name string
		args args
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "just after",
			args: args{
				filters: Checks{After(SameID("https://example.com"))},
				it: vocab.ItemCollection{
					vocab.Activity{ID: "https://example.com/1"},
					vocab.Activity{ID: "https://example.com"},
					vocab.Activity{ID: "https://example.com/2"},
				},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
		{
			name: "after with filters",
			args: args{
				filters: Checks{After(SameID("example.com"), HasType("Activity"), After(SameID("example.com")))},
			},
			want: nil,
		},
		{
			name: "maxItems=2 of 2",
			args: args{
				filters: Checks{WithMaxCount(2)},
				it: vocab.ItemCollection{
					vocab.Activity{ID: "https://example.com/1"},
					vocab.Activity{ID: "https://example.com/2"},
				},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
		{
			name: "maxItems=2 of 3",
			args: args{
				filters: Checks{WithMaxCount(2)},
				it: vocab.ItemCollection{
					vocab.Activity{ID: "https://example.com/1"},
					vocab.Activity{ID: "https://example.com/2"},
					vocab.Activity{ID: "https://example.com/3"},
				},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
		{
			name: "before=https://example.com/1 single item",
			args: args{
				filters: Checks{Before(SameID("https://example.com/1"))},
				it: vocab.ItemCollection{
					vocab.Activity{ID: "https://example.com/1"},
				},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "before=https://example.com/1 second item",
			args: args{
				filters: Checks{Before(SameID("https://example.com/1"))},
				it: vocab.ItemCollection{
					vocab.Activity{ID: "https://example.com/0"},
					vocab.Activity{ID: "https://example.com/1"},
					vocab.Activity{ID: "https://example.com/2"},
				},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/0"},
			},
		},
		{
			name: "after=https://example.com/1 first item",
			args: args{
				filters: Checks{After(SameID("https://example.com/1"))},
				it: vocab.ItemCollection{
					vocab.Activity{ID: "https://example.com/1"},
				},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "after=https://example.com/1 second item",
			args: args{
				filters: Checks{After(SameID("https://example.com/1"))},
				it: vocab.ItemCollection{
					vocab.Activity{ID: "https://example.com/0"},
					vocab.Activity{ID: "https://example.com/1"},
					vocab.Activity{ID: "https://example.com/2"},
				},
			},
			want: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/2"},
			},
		},
		{
			name: "no pagination",
			args: args{
				it: vocab.ItemCollection{
					vocab.Object{ID: "https://example.com/0", Type: vocab.NoteType},
					vocab.Object{ID: "https://example.com/1", Type: vocab.ArticleType},
					vocab.Object{ID: "https://example.com/2", Type: vocab.DocumentType},
					vocab.Object{ID: "https://example.com/3", Type: vocab.NoteType},
					vocab.Object{ID: "https://example.com/4", Type: vocab.ArticleType},
					vocab.Object{ID: "https://example.com/5", Type: vocab.DocumentType},
					vocab.Object{ID: "https://example.com/6", Type: vocab.NoteType},
					vocab.Object{ID: "https://example.com/7", Type: vocab.ArticleType},
					vocab.Object{ID: "https://example.com/8", Type: vocab.DocumentType},
				},
				filters: Checks{HasType("Note", "Article", "Document")},
			},
			want: vocab.ItemCollection{
				vocab.Object{ID: "https://example.com/0", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/1", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/2", Type: vocab.DocumentType},
				vocab.Object{ID: "https://example.com/3", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/4", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/5", Type: vocab.DocumentType},
				vocab.Object{ID: "https://example.com/6", Type: vocab.NoteType},
				vocab.Object{ID: "https://example.com/7", Type: vocab.ArticleType},
				vocab.Object{ID: "https://example.com/8", Type: vocab.DocumentType},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PaginateCollection(tt.args.it, tt.args.filters...); !cmp.Equal(got, tt.want) {
				t.Errorf("PaginateCollection() = %s", cmp.Diff(tt.want, got))
			}
		})
	}
}
