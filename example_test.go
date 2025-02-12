package filters

import (
	"fmt"

	vocab "github.com/go-ap/activitypub"
)

func ExampleFilters() {
	collection := vocab.ItemCollection{
		// doesn't match due to Actor ID
		vocab.Create{
			Type:   "Create",
			Actor:  vocab.IRI("https://example.com/bob"),
			Object: vocab.IRI("https//example.com/test"),
		},
		// doesn't match due to nil Object
		vocab.Create{
			Type:  "Create",
			Actor: vocab.IRI("https://example.com/jdoe"),
		},
		// match
		vocab.Create{
			Type:   "Create",
			Actor:  vocab.IRI("https://example.com/jdoe"),
			Object: vocab.IRI("https//example.com/test"),
		},
		// match
		vocab.Create{
			Type: "Create",
			Actor: vocab.Person{
				ID:   "https://example.com/jdoe1",
				Name: vocab.DefaultNaturalLanguageValue("JohnDoe"),
			},
			Object: vocab.IRI("https//example.com/test"),
		},
		// doesn't match due to the activity Type
		vocab.Follow{Type: "Follow"},
		// doesn't match due to Arrive being an intransitive activity
		vocab.Arrive{Type: "Arrive"},
		// doesn't match due to Question being an intransitive activity
		vocab.Question{Type: "Question"},
	}
	// This filters all activities that are not:
	// Create activities,
	// published by an Actor with the ID https://example.com/authors/jdoe, or with the name "JohnDoe"
	// and, which have an object with a non nil ID.
	filterFn := All(
		HasType("Create"),
		Actor(
			Any(
				SameID("https://example.com/jdoe"),
				NameIs("JohnDoe"),
			),
		),
		Object(Not(NilID)),
	)

	result := make(vocab.ItemCollection, 0)
	for _, it := range collection {
		if filterFn.Match(it) {
			result = append(result, it)
		}
	}
	output, _ := vocab.MarshalJSON(result)
	fmt.Printf("Result[%d]: %s", len(result), output)
	// Output: Result[2]: [{"type":"Create","actor":"https://example.com/jdoe","object":"https//example.com/test"},{"type":"Create","actor":{"id":"https://example.com/jdoe1","name":"JohnDoe"},"object":"https//example.com/test"}]
}
