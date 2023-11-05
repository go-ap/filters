package filters

import (
	"reflect"
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func TestCursor(t *testing.T) {
	type args struct {
		filters Checks
		after   Check
		before  Check
		limit   Check
	}
	tests := []struct {
		name string
		args args
		it   vocab.Item
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "maxItems=2 of 2",
			args: args{
				limit: WithMaxCount(2),
			},
			it: &vocab.ItemCollection{
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
			args: args{
				limit: WithMaxCount(2),
			},
			it: &vocab.ItemCollection{
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
			name: "check=https://example.com/1 single item",
			args: args{
				before: ID("https://example.com/1"),
			},
			it: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "check=https://example.com/1 second item",
			args: args{
				before: ID("https://example.com/1"),
			},
			it: vocab.ItemCollection{
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
			args: args{
				after: ID("https://example.com/1"),
			},
			it: vocab.ItemCollection{
				vocab.Activity{ID: "https://example.com/1"},
			},
			want: vocab.ItemCollection{},
		},
		{
			name: "after=https://example.com/1 second item",
			args: args{
				after: ID("https://example.com/1"),
			},
			it: vocab.ItemCollection{
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
			rgs := tt.args
			curFns := make(Checks, 0)
			if rgs.before != nil {
				curFns = append(curFns, Before(rgs.before))
			}
			if rgs.after != nil {
				curFns = append(curFns, After(rgs.after))
			}
			if rgs.limit != nil {
				curFns = append(curFns, rgs.limit)
			}
			c := Cursor(curFns...)
			if got := c.Run(tt.it); !vocab.ItemsEqual(tt.want, got) {
				t.Errorf("Cursor() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestCursorFns(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		fns  []Check
		item vocab.Item
		want vocab.Item
	}{
		{
			name: "just after",
			fns:  Checks{After(ID("example.com"))},
			item: nil,
			want: nil,
		},
		{
			name: "after with filters",
			fns:  Checks{ID("example.com"), HasType("Activity"), After(ID("example.com"))},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChecks := CursorFns(tt.fns...)
			if got := gotChecks.Run(tt.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CursorFns() = %v, want %v", got, tt.want)
			}
		})
	}
}
