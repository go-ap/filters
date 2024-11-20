package filters

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters/index"
)

func ExampleSearchIndex() {
	activities := []vocab.LinkOrIRI{
		&vocab.Activity{
			ID:    "https://federated.local/1",
			Type:  vocab.CreateType,
			To:    vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
			Actor: vocab.IRI("https://federated.local/~jdoe"),
			Object: &vocab.Object{
				ID:   "https://federated.local/objects/1",
				Type: vocab.PageType,
				Name: vocab.NaturalLanguageValues{{Ref: "-", Value: vocab.Content("Link to example.com")}},
				URL:  vocab.IRI("https://example.com"),
			},
		},
		&vocab.Activity{
			ID:     "https://federated.local/2",
			Type:   vocab.LikeType,
			To:     vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
			Actor:  vocab.IRI("https://federated.local/~jdoe"),
			Object: vocab.IRI("https://federated.local/objects/1"),
		},
		&vocab.Activity{
			ID:     "https://federated.local/3",
			Type:   vocab.DislikeType,
			To:     vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
			Actor:  vocab.IRI("https://federated.local/~jdoe"),
			Object: vocab.IRI("https://federated.local/objects/1"),
		},
		&vocab.Activity{
			ID:     "https://federated.local/4",
			Type:   vocab.FlagType,
			To:     vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
			Actor:  vocab.IRI("https://federated.local/~jdoe"),
			Object: vocab.IRI("https://federated.local/objects/1"),
		},
	}

	in := index.Full()
	// Add the activities to the index
	_ = in.Add(activities...)

	findCreate := Checks{
		HasType(vocab.CreateType),
		Object(SameID("https://federated.local/objects/1")),
	}
	iris, err := SearchIndex(in, findCreate...)
	fmt.Printf("Find Create:\n")
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("IRIs: %#v\n", iris)

	findBlock := Checks{
		HasType(vocab.FlagType),
		AttributedToLike("https://federated.local/~jdoe"),
	}
	iris, err = SearchIndex(in, findBlock...)
	fmt.Printf("Find Flag:\n")
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("IRIs: %#v\n", iris)

	// Output:
	// Find Create:
	// Error: <nil>
	// IRIs: []activitypub.IRI{https://federated.local/1}
	// Find Flag:
	// Error: <nil>
	// IRIs: []activitypub.IRI{https://federated.local/4}
}

var indexableActivities = []vocab.LinkOrIRI{
	&vocab.Actor{
		ID:                "https://federated.local/~jdoe",
		Type:              vocab.PersonType,
		AttributedTo:      vocab.IRI("https://federated.local/~jdoe"),
		PreferredUsername: vocab.DefaultNaturalLanguageValue("jDoe"),
		Summary:           vocab.DefaultNaturalLanguageValue("An anonymous dude"),
	},
	&vocab.Actor{
		ID:                "https://federated.local/~alice",
		Type:              vocab.PersonType,
		AttributedTo:      vocab.IRI("https://federated.local"),
		Name:              vocab.DefaultNaturalLanguageValue("Alice in Wonderland"),
		PreferredUsername: vocab.DefaultNaturalLanguageValue("alice"),
	},
	&vocab.Object{
		ID:           "https://federated.local/objects/1",
		Type:         vocab.PageType,
		AttributedTo: vocab.IRI("https://federated.local/~alice"),
		To:           vocab.ItemCollection{vocab.PublicNS},
		Name:         vocab.DefaultNaturalLanguageValue("Link to example.com"),
		Summary:      vocab.DefaultNaturalLanguageValue("An example for a link to example.com"),
		URL:          vocab.IRI("https://example.com"),
	},
	&vocab.Activity{
		ID:     "https://federated.local/1",
		Type:   vocab.CreateType,
		To:     vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
		CC:     vocab.ItemCollection{vocab.PublicNS},
		Actor:  vocab.IRI("https://federated.local/~jdoe"),
		Object: vocab.IRI("https://federated.local/objects/1"),
	},
	&vocab.Activity{
		ID:     "https://federated.local/2",
		Type:   vocab.LikeType,
		To:     vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
		Actor:  vocab.IRI("https://federated.local/~jdoe"),
		Object: vocab.IRI("https://federated.local/objects/1"),
	},
	&vocab.Activity{
		ID:     "https://federated.local/3",
		Type:   vocab.DislikeType,
		To:     vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
		Actor:  vocab.IRI("https://federated.local/~jdoe"),
		Object: vocab.IRI("https://federated.local/objects/1"),
	},
	&vocab.Activity{
		ID:     "https://federated.local/4",
		Type:   vocab.FlagType,
		To:     vocab.ItemCollection{vocab.IRI("https://federated.local/~jdoe")},
		Actor:  vocab.IRI("https://federated.local/~jdoe"),
		Object: vocab.IRI("https://federated.local/objects/1"),
	},
	&vocab.Activity{
		ID:      "https://federated.local/5",
		Type:    vocab.FlagType,
		Content: vocab.DefaultNaturalLanguageValue("flagged object"),
		To:      vocab.ItemCollection{vocab.IRI("https://federated.local/~alice")},
		Actor:   vocab.IRI("https://federated.local/~alice"),
		Object:  vocab.IRI("https://federated.local/objects/1"),
	},
}

func buildIndex() map[index.Type]index.Indexable {
	f := index.Full()
	_ = f.Add(indexableActivities...)
	return f.Indexes
}

func wantedBmp[T ~string](x ...T) *roaring.Bitmap {
	dat := make([]uint32, len(x))
	for i, tt := range x {
		dat[i] = index.HashFn(vocab.IRI(tt))
	}
	return roaring.BitmapOf(dat...)
}

func TestChecks_IndexMatch(t *testing.T) {
	idx := buildIndex()

	tests := []struct {
		name    string
		ff      Checks
		indexes map[index.Type]index.Indexable
		want    *roaring.Bitmap
	}{
		{
			name: "empty",
			ff:   nil,
		},
		{
			name: "byType Flag",
			ff: Checks{
				HasType(vocab.FlagType),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/4", "https://federated.local/5"),
		},
		{
			name: "byType Flag with actor",
			ff: Checks{
				HasType(vocab.FlagType),
				Actor(SameID("https://federated.local/~jdoe")),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/4"),
		},
		{
			name: "byType Flag with object",
			ff: Checks{
				HasType(vocab.FlagType),
				Object(SameID("https://federated.local/objects/1")),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/4", "https://federated.local/5"),
		},
		{
			name: "byType Page",
			ff: Checks{
				HasType(vocab.PageType),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1"),
		},
		{
			name: "by content",
			ff: Checks{
				ContentLike("flagged"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/5"),
		},
		{
			name: "by name",
			ff: Checks{
				NameIs("Link"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1"),
		},
		{
			name: "by summary",
			ff: Checks{
				SummaryIs("example"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1"),
		},
		{
			name:    "by ID",
			ff:      Checks{SameID("https://federated.local/objects/1")},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1"),
		},
		{
			name: "by summary",
			ff: Checks{
				SummaryIs("example"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1"),
		},
		{
			name: "authorized:public",
			ff: Checks{
				Authorized("https://www.w3.org/ns/activitystreams#Public"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1", "https://federated.local/1"),
		},
		{
			name:    "authorized:~alice",
			ff:      Checks{Authorized("https://federated.local/~alice")},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1", "https://federated.local/1", "https://federated.local/5"),
		},
		{
			name: "by recipients",
			ff: Checks{
				Recipients("https://federated.local/~alice"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/5"),
		},
		{
			name: "by attributedTo",
			ff: Checks{
				Recipients("https://federated.local/~alice"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/5"),
		},
		{
			name: "by attributedTo",
			ff: Checks{
				SameAttributedTo("https://federated.local/~alice"),
			},
			indexes: idx,
			want:    wantedBmp("https://federated.local/objects/1"),
		},
		{
			name: "anyOf(type:Flag,attributedTo:~alice)",
			ff: Checks{
				Any(
					HasType(vocab.FlagType),
					SameAttributedTo("https://federated.local/~alice"),
				),
			},
			indexes: idx,
			want: wantedBmp(
				"https://federated.local/objects/1",
				"https://federated.local/4",
				"https://federated.local/5",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ff.IndexMatch(tt.indexes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
