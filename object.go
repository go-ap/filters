package filters

import (
	"net/url"
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
	fragStr, _ := url.QueryUnescape(string(frag))
	return strings.Contains(nfc(item.GetID().String()), nfc(fragStr))
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

func accumURLs(item vocab.Item) vocab.IRIs {
	urls := make(vocab.IRIs, 0)
	if vocab.LinkTypes.Contains(item.GetType()) {
		_ = vocab.OnLink(item, func(lnk *vocab.Link) error {
			urls = append(urls, lnk.Href)
			return nil
		})
	} else {
		_ = vocab.OnObject(item, func(ob *vocab.Object) error {
			if vocab.IsNil(ob.URL) {
				return nil
			}
			if ob.URL.IsObject() {
				_ = vocab.OnObject(ob.URL, func(url *vocab.Object) error {
					urls = append(urls, url.GetLink())
					return nil
				})
			} else {
				urls = append(urls, ob.URL.GetLink())
			}
			return nil
		})
	}
	return urls
}

type urlEquals vocab.IRI

func (i urlEquals) Apply(item vocab.Item) bool {
	if vocab.IsNil(item) {
		return len(i) == 0
	}
	return accumURLs(item).Contains(vocab.IRI(i))
}

// SameURL checks an activitypub.Object's IRI
func SameURL(iri vocab.IRI) Check {
	return urlEquals(iri)
}

type urlLike string

func (frag urlLike) Apply(item vocab.Item) bool {
	nfc := norm.NFC.String
	fragStr, _ := url.QueryUnescape(string(frag))
	for _, u := range accumURLs(item) {
		if strings.Contains(nfc(u.String()), nfc(fragStr)) {
			return true
		}
	}
	return false
}

func URLLike(frag string) Check {
	return urlLike(frag)
}

type nilURL struct{}

func (n nilURL) Apply(it vocab.Item) bool {
	return vocab.IsNil(it) || Any(SameURL(vocab.NilIRI), SameURL(vocab.EmptyIRI)).Apply(it)
}

// NilURL checks if the activitypub.Object's URL property matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
var NilURL = nilURL{}
