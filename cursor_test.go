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
