package filters

import (
	"fmt"

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
