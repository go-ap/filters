package filters

import (
	"strings"

	vocab "github.com/go-ap/activitypub"
	"golang.org/x/text/unicode/norm"
)

// NilID checks if the activitypub.Object's ID property matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
func NilID(it vocab.Item) bool {
	// NOTE(marius): I'm not sure that a nil Item returning true, is entirely sane/safe logic
	return vocab.IsNil(it) || Any(SameIRI(vocab.NilIRI), SameIRI(vocab.EmptyIRI))(it.GetLink())
}

// ID checks an activitypub.Object's ID property against the received iri.
func ID(iri vocab.IRI) Fn {
	return func(item vocab.Item) bool {
		return item.GetID().Equals(iri, true)
	}
}

func IDLike(frag string) Fn {
	return func(item vocab.Item) bool {
		nfc := norm.NFC.String
		return strings.Contains(nfc(item.GetID().String()), nfc(frag))
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
		_ = vocab.OnObject(it, func(object *vocab.Object) error {
			result = types.Contains(it.GetType())
			return nil
		})
		return result
	}
}
