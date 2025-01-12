package filters

import (
	"net/url"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"golang.org/x/text/unicode/norm"
)

// NilID checks if the [vocab.Object]'s ID property matches any of the two magic values
// that denote an empty value: [vocab.NilID] = "-", or [vocab.EmptyID] = ""
var NilID = idNil{}

// NotNilID checks if the [vocab.Object]'s ID property is not nil
var NotNilID = Not(NilID)

type idNil iriNil

func (n idNil) Match(it vocab.Item) bool {
	return vocab.IsNil(it) || Any(SameIRI(vocab.NilIRI), SameIRI(vocab.EmptyIRI)).Match(it.GetID())
}

// SameID checks a [vocab.Object]'s ID property against the received iri.
func SameID(i vocab.IRI) Check {
	return idEquals(i)
}

type idEquals iriEquals

func (i idEquals) Match(item vocab.Item) bool {
	if vocab.IsNil(item) {
		return len(i) == 0
	}
	return item.GetID().Equals(vocab.IRI(i), false)
}

// IDLike
func IDLike(frag string) Check {
	return idLike(frag)
}

type idLike iriLike

func (l idLike) Match(item vocab.Item) bool {
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

func (types withTypes) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(types) == 0
	}
	return vocab.ActivityVocabularyTypes(types).Contains(it.GetType())
}

func accumURLs(item vocab.Item) vocab.IRIs {
	var urls vocab.ItemCollection
	switch it := item.(type) {
	case vocab.Item:
		_ = vocab.OnObject(it, func(ob *vocab.Object) error {
			urls = vocab.DerefItem(ob.URL)
			return nil
		})
	case vocab.Link:
		_ = vocab.OnLink(item, func(lnk *vocab.Link) error {
			urls = vocab.ItemCollection{lnk.Href}
			return nil
		})
	}
	return urls.IRIs()
}

// SameURL checks an activitypub.Object's IRI
func SameURL(iri vocab.IRI) Check {
	return urlEquals(iri)
}

type urlEquals iriEquals

func (i urlEquals) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(i) == 0
	}
	return accumURLs(it).Contains(vocab.IRI(i))
}

type urlLike iriLike

func (frag urlLike) Match(it vocab.Item) bool {
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
	var items vocab.ItemCollection
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		items = vocab.DerefItem(ob.Context)
		return nil
	})
	return items.IRIs()
}

type urlNil iriNil

func URLNil() Check {
	return iriNil{}
}

func (frag urlNil) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return len(accumURLs(it)) > 0
}

func SameContext(iri vocab.IRI) Check {
	return contextEquals(iri)
}

type contextEquals iriEquals

func (c contextEquals) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(c) == 0
	}

	return accumContexts(it).Contains(vocab.IRI(c))
}

func ContextLike(frag string) Check {
	return contextLike(frag)
}

type contextLike iriLike

func (c contextLike) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(c) == 0
	}

	return accumContexts(it).Contains(vocab.IRI(c))
}

var NilContext = contextNil{}

type contextNil iriNil

func (c contextNil) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return len(accumContexts(it)) == 0
}

func accumAttributedTos(item vocab.Item) vocab.IRIs {
	var items vocab.ItemCollection
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		items = vocab.DerefItem(ob.AttributedTo)
		return nil
	})
	return items.IRIs()
}

// SameAttributedTo creates a filter that checks the [vocab.IRI] against the attributedTo property of the item
// it gets applied on.
func SameAttributedTo(iri vocab.IRI) Check {
	return attributedToEquals(iri)
}

type attributedToEquals iriEquals

func (a attributedToEquals) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(a) == 0
	}

	return accumAttributedTos(it).Contains(vocab.IRI(a))
}

// AttributedToLike creates a filter that checks the [vocab.IRI] against the attributedTo property of the item
// it gets applied on using a similarity match.
func AttributedToLike(frag string) Check {
	return attributedToLike(frag)
}

type attributedToLike iriLike

func (a attributedToLike) Match(it vocab.Item) bool {
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

type attributedToNil iriNil

func (a attributedToNil) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return len(accumAttributedTos(it)) == 0
}

func accumInReplyTos(item vocab.Item) vocab.IRIs {
	var iris vocab.ItemCollection
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		iris = vocab.DerefItem(ob.InReplyTo)
		return nil
	})
	return iris.IRIs()
}

var NilInReplyTo = inReplyToNil{}

type inReplyToNil iriNil

func (c inReplyToNil) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return len(accumInReplyTos(it)) == 0
}

func InReplyToLike(frag string) Check {
	return inReplyToLike(frag)
}

type inReplyToLike iriLike

func (a inReplyToLike) Match(it vocab.Item) bool {
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

func (i inReplyToEquals) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return len(i) == 0
	}
	return accumInReplyTos(it).Contains(vocab.IRI(i))
}
