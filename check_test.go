package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
	"github.com/google/go-cmp/cmp"
)

func TestChecks_Filter(t *testing.T) {
	tests := []struct {
		name string
		ff   Checks
		item vocab.Item
		want vocab.Item
	}{
		{
			name: "empty",
			want: nil,
		},
		{
			name: "type that doesn't match",
			ff:   Checks{HasType("t1")},
			item: &vocab.Object{Type: vocab.ActivityVocabularyType("t2")},
			want: nil,
		},
		{
			name: "type that matches",
			ff:   Checks{HasType("t1")},
			item: &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
			want: &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ff.Filter(tt.item); !cmp.Equal(got, tt.want) {
				t.Errorf("Filter() = %s", cmp.Diff(tt.want, got, cmp.Comparer(vocab.ItemsEqual)))
			}
		})
	}
}

func TestChecks_Paginate(t *testing.T) {
	tests := []struct {
		name string
		ff   Checks
		item vocab.Item
		want vocab.Item
	}{
		{
			name: "empty",
			want: nil,
		},
		{
			name: "no pagination type that doesn't match",
			ff:   Checks{HasType("t1")},
			item: &vocab.Object{Type: vocab.ActivityVocabularyType("t2")},
			want: &vocab.Object{Type: vocab.ActivityVocabularyType("t2")},
		},
		{
			name: "no pagination type that matches",
			ff:   Checks{HasType("t1")},
			item: &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
			want: &vocab.Object{Type: vocab.ActivityVocabularyType("t1")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ff.Paginate(tt.item); !cmp.Equal(got, tt.want, cmp.Comparer(vocab.ItemsEqual)) {
				t.Errorf("Paginate() = %s", cmp.Diff(tt.want, got, cmp.Comparer(vocab.ItemsEqual)))
			}
		})
	}
}

func TestChecks_Run(t *testing.T) {
	tests := []struct {
		name string
		ff   Checks
		item vocab.Item
		want vocab.Item
	}{
		{
			name: "empty",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ff.Run(tt.item); !cmp.Equal(got, tt.want, cmp.Comparer(vocab.ItemsEqual)) {
				t.Errorf("Run() = %s", cmp.Diff(tt.want, got, cmp.Comparer(vocab.ItemsEqual)))
			}
		})
	}
}

func TestChecks_runOnItem(t *testing.T) {
	tests := []struct {
		name string
		ff   Checks
		it   vocab.Item
		want vocab.Item
	}{
		{
			name: "empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ff.runOnItem(tt.it); !cmp.Equal(got, tt.want, cmp.Comparer(vocab.ItemsEqual)) {
				t.Errorf("runOnItem() = %s", cmp.Diff(tt.want, got, cmp.Comparer(vocab.ItemsEqual)))
			}
		})
	}
}

func TestChecks_runOnItems(t *testing.T) {
	tests := []struct {
		name string
		ff   Checks
		col  vocab.ItemCollection
		want vocab.ItemCollection
	}{
		{
			name: "empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ff.runOnItems(tt.col); !cmp.Equal(got, tt.want, cmp.Comparer(vocab.ItemsEqual)) {
				t.Errorf("runOnItems() =  %s", cmp.Diff(tt.want, got, cmp.Comparer(vocab.ItemsEqual)))
			}
		})
	}
}

func Test_checkFn(t *testing.T) {
	t.Skipf("can't compare functions")
	tests := []struct {
		name string
		ff   Checks
		want func(vocab.Item) bool
	}{
		{
			name: "empty",
			ff:   nil,
			want: nilCheck,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkFn(tt.ff); !cmp.Equal(got, tt.want, cmp.Comparer(sameFns)) {
				t.Errorf("checkFn() = %s", cmp.Diff(tt.want, got, cmp.Comparer(sameFns)))
			}
		})
	}
}
