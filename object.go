package filters

import (
	"net/url"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"golang.org/x/text/unicode/norm"
)

// NilID checks if the activitypub.Object's ID property matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
var NilID = idNil{}

// NotNilID checks if the activitypub.Object's ID property is not nil
var NotNilID = Not(NilID)

type idNil iriNil

func (n idNil) Apply(it vocab.Item) bool {
	return vocab.IsNil(it) || Any(SameIRI(vocab.NilIRI), SameIRI(vocab.EmptyIRI)).Apply(it.GetID())
}

// SameID checks an activitypub.Object's ID property against the received iri.
func SameID(i vocab.IRI) Check {
	return idEquals(i)
}

type idEquals iriEquals

func (i idEquals) Apply(item vocab.Item) bool {
	if vocab.IsNil(item) {
		return len(i) == 0
	}
	return item.GetID().Equals(vocab.IRI(i), true)
}

// IDLike
func IDLike(frag string) Check {
	return idLike(frag)
}

type idLike iriLike

func (l idLike) Apply(item vocab.Item) bool {
	if vocab.IsNil(item) {
		return false
	}
	nfc := norm.NFC.String
	fragStr, _ := url.QueryUnescape(string(l))
	return strings.Contains(nfc(item.GetID().String()), nfc(fragStr))
}

// HasType checks an activitypub.Object's Type property against a series of values.
// If any of the ty values matches, the function returns true.
func HasType(ty ...vocab.ActivityVocabularyType) Check {
	return withTypes(ty)
}

type withTypes vocab.ActivityVocabularyTypes

func (types withTypes) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(types) == 0
	}
	return vocab.ActivityVocabularyTypes(types).Contains(it.GetType())
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

// SameURL checks an activitypub.Object's IRI
func SameURL(iri vocab.IRI) Check {
	return urlEquals(iri)
}

type urlEquals iriEquals

func (i urlEquals) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(i) == 0
	}
	return accumURLs(it).Contains(vocab.IRI(i))
}

type urlLike iriLike

func (frag urlLike) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(frag) == 0
	}
	nfc := norm.NFC.String
	fragStr, _ := url.QueryUnescape(string(frag))
	for _, u := range accumURLs(it) {
		if strings.Contains(nfc(u.String()), nfc(fragStr)) {
			return true
		}
	}
	return false
}

func URLLike(frag string) Check {
	return urlLike(frag)
}

func SameContext(iri vocab.IRI) Check {
	return iriEquals(iri)
}

type contextEquals iriEquals

func (c contextEquals) Apply(item vocab.Item) bool {
	return iriEquals(c).Apply(item)
}

func ContextLike(frag string) Check {
	return iriLike(frag)
}

type contextLike iriLike

func (c contextLike) Apply(item vocab.Item) bool {
	return iriLike(c).Apply(item)
}

var NilContext = idNil{}

type contextNil idNil

func (c contextNil) Apply(item vocab.Item) bool {
	return idNil(c).Apply(item)
}
