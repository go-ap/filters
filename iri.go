package filters

import (
	"net/url"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"golang.org/x/text/unicode/norm"
)

type iriEquals vocab.IRI

func (i iriEquals) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(i) == 0
	}
	return it.GetLink().Equals(vocab.IRI(i), true)
}

// SameIRI checks an activitypub.Object's IRI
func SameIRI(iri vocab.IRI) Check {
	return iriEquals(iri)
}

type iriLike string

func (frag iriLike) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	nfc := norm.NFC.String
	fragStr, _ := url.QueryUnescape(string(frag))
	return strings.Contains(nfc(it.GetLink().String()), nfc(fragStr))
}

// IRILike
func IRILike(frag string) Check {
	return iriLike(frag)
}

// NilIRI checks if the activitypub.Item's IRI that matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
var NilIRI = iriNil{}

type iriNil struct{}

func (n iriNil) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	if vocab.IsItemCollection(it) {
		result := false
		_ = vocab.OnItemCollection(it, func(col *vocab.ItemCollection) error {
			result = len(*col) == 0
			return nil
		})
		return result
	}
	return Any(SameIRI(vocab.NilIRI), SameIRI(vocab.EmptyIRI)).Apply(it.GetLink())
}

// NotNilIRI checks if the activitypub.Object's URL property matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
var NotNilIRI = Not(iriNil{})
