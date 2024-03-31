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

func accumContexts(item vocab.Item) vocab.IRIs {
	iris := make(vocab.IRIs, 0)
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		if vocab.IsNil(ob.Context) {
			return nil
		}
		if ob.AttributedTo.IsObject() {
			_ = vocab.OnObject(ob.Context, func(c *vocab.Object) error {
				iris = append(iris, c.GetLink())
				return nil
			})
		} else {
			iris = append(iris, ob.Context.GetLink())
		}
		return nil
	})
	return iris
}

func SameContext(iri vocab.IRI) Check {
	return contextEquals(iri)
}

type contextEquals iriEquals

func (c contextEquals) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(c) == 0
	}

	return accumContexts(it).Contains(vocab.IRI(c))
}

func ContextLike(frag string) Check {
	return contextLike(frag)
}

type contextLike iriLike

func (c contextLike) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(c) == 0
	}

	return accumContexts(it).Contains(vocab.IRI(c))
}

var NilContext = contextNil{}

type contextNil idNil

func (c contextNil) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return accumContexts(it) == nil
}

func accumAttributedTos(item vocab.Item) vocab.IRIs {
	iris := make(vocab.IRIs, 0)
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		if vocab.IsNil(ob.AttributedTo) {
			return nil
		}
		if ob.AttributedTo.IsObject() {
			_ = vocab.OnObject(ob.AttributedTo, func(attrTo *vocab.Object) error {
				iris = append(iris, attrTo.GetLink())
				return nil
			})
		} else {
			iris = append(iris, ob.AttributedTo.GetLink())
		}
		return nil
	})
	return iris
}

func SameAttributedTo(iri vocab.IRI) Check {
	return attributedToEquals(iri)
}

type attributedToEquals iriEquals

func (a attributedToEquals) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(a) == 0
	}

	return accumAttributedTos(it).Contains(vocab.IRI(a))
}

func AttributedToLike(frag string) Check {
	return attributedToLike(frag)
}

type attributedToLike iriLike

func (a attributedToLike) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(a) == 0
	}
	nfc := norm.NFC.String
	fragStr, _ := url.QueryUnescape(string(a))
	for _, u := range accumAttributedTos(it) {
		if strings.Contains(nfc(u.String()), nfc(fragStr)) {
			return true
		}
	}
	return false
}

var NilAttributedTo = attributedToNil{}

type attributedToNil idNil

func (a attributedToNil) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return accumAttributedTos(it) == nil
}

func accumInReplyTos(item vocab.Item) vocab.IRIs {
	iris := make(vocab.IRIs, 0)
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		if vocab.IsNil(ob.InReplyTo) {
			return nil
		}
		if ob.AttributedTo.IsObject() {
			_ = vocab.OnObject(ob.InReplyTo, func(inReplyTo *vocab.Object) error {
				iris = append(iris, inReplyTo.GetLink())
				return nil
			})
		} else {
			iris = append(iris, ob.InReplyTo.GetLink())
		}
		return nil
	})
	return iris
}

var NilInReplyTo = inReplyToNil{}

type inReplyToNil idNil

func (c inReplyToNil) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return accumInReplyTos(it) == nil
}

func InReplyToLike(frag string) Check {
	return inReplyToLike(frag)
}

type inReplyToLike iriLike

func (a inReplyToLike) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(a) == 0
	}
	nfc := norm.NFC.String
	fragStr, _ := url.QueryUnescape(string(a))
	for _, u := range accumInReplyTos(it) {
		if strings.Contains(nfc(u.String()), nfc(fragStr)) {
			return true
		}
	}
	return false
}

// SameInReplyTo checks an activitypub.Object's InReplyTo
func SameInReplyTo(iri vocab.IRI) Check {
	return inReplyToEquals(iri)
}

type inReplyToEquals iriEquals

func (i inReplyToEquals) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(i) == 0
	}
	return accumInReplyTos(it).Contains(vocab.IRI(i))
}
