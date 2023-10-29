package filters

import (
	"strings"

	vocab "github.com/go-ap/activitypub"
)

// NilID checks if the activitypub.Object's ID property matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
func NilID(it vocab.Item) bool {
	return Any(SameIRI(vocab.NilIRI), SameIRI(vocab.EmptyIRI))(it.GetLink())
}

// NotNilID checks if the activitypub.Object's ID property doesn't match any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
func NotNilID(it vocab.Item) bool {
	return !NilID(it)
}

// ID checks an activitypub.Object's ID property against the received iri.
func ID(iri vocab.IRI) Fn {
	return func(item vocab.Item) bool {
		return item.GetID().Equals(iri, true)
	}
}

// SameIRI checks an activitypub.Object's IRI
func SameIRI(iri vocab.IRI) Fn {
	return func(item vocab.Item) bool {
		return item.GetLink().Equals(iri, true)
	}
}

// HasType checks an activitypub.Object's Type property against a series of values.
// If any of the ty values matches, the function returns true.
func HasType(ty ...vocab.ActivityVocabularyType) Fn {
	types := vocab.ActivityVocabularyTypes(ty)
	return func(it vocab.Item) bool {
		result := false
		vocab.OnObject(it, func(object *vocab.Object) error {
			result = types.Contains(it.GetType())
			return nil
		})
		return result
	}
}
