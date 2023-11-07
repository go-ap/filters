package filters

import (
	"strings"

	vocab "github.com/go-ap/activitypub"
	"golang.org/x/text/unicode/norm"
)

type idEquals vocab.IRI

func (i idEquals) Apply(item vocab.Item) bool {
	if vocab.IsNil(item) {
		return len(i) == 0
	}
	return item.GetID().Equals(vocab.IRI(i), true)
}

type nilId struct{}

func (n nilId) Apply(it vocab.Item) bool {
	return vocab.IsNil(it) || Any(SameIRI(vocab.NilIRI), SameIRI(vocab.EmptyIRI)).Apply(it.GetLink())
}

// NilID checks if the activitypub.Object's ID property matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
var NilID = nilId{}

// ID checks an activitypub.Object's ID property against the received iri.
func ID(i vocab.IRI) Check {
	return idEquals(i)
}

type iriLike string

func (frag iriLike) Apply(item vocab.Item) bool {
	nfc := norm.NFC.String
	return strings.Contains(nfc(item.GetID().String()), nfc(string(frag)))
}

func IDLike(frag string) Check {
	return iriLike(frag)
}

type iriEquals vocab.IRI

func (i iriEquals) Apply(item vocab.Item) bool {
	if vocab.IsNil(item) {
		return len(i) == 0
	}
	return item.GetLink().Equals(vocab.IRI(i), true)
}

// SameIRI checks an activitypub.Object's IRI
func SameIRI(iri vocab.IRI) Check {
	return iriEquals(iri)
}

type withTypes vocab.ActivityVocabularyTypes

func (types withTypes) Apply(it vocab.Item) bool {
	return vocab.ActivityVocabularyTypes(types).Contains(it.GetType())
}

// HasType checks an activitypub.Object's Type property against a series of values.
// If any of the ty values matches, the function returns true.
func HasType(ty ...vocab.ActivityVocabularyType) Check {
	return withTypes(ty)
}
