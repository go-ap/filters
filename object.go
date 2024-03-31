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
			if vocab.IsIRI(ob.URL) {
				urls = append(urls, ob.URL.GetLink())
			} else if vocab.IsIRIs(ob.URL) {
				_ = vocab.OnIRIs(ob.URL, func(replTos *vocab.IRIs) error {
					for _, r := range *replTos {
						urls = append(urls, r.GetLink())
					}
					return nil
				})
			} else if vocab.IsItemCollection(ob.URL) {
				_ = vocab.OnItemCollection(ob.URL, func(uc *vocab.ItemCollection) error {
					for _, u := range *uc {
						urls = append(urls, u.GetLink())
					}
					return nil
				})
			} else {
				_ = vocab.OnObject(ob.URL, func(url *vocab.Object) error {
					urls = append(urls, url.GetLink())
					return nil
				})
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
		if vocab.IsIRI(ob.Context) {
			iris = append(iris, ob.Context.GetLink())
		} else if vocab.IsIRIs(ob.Context) {
			_ = vocab.OnIRIs(ob.Context, func(col *vocab.IRIs) error {
				for _, r := range *col {
					iris = append(iris, r.GetLink())
				}
				return nil
			})
		} else if vocab.IsItemCollection(ob.Context) {
			_ = vocab.OnItemCollection(ob.Context, func(col *vocab.ItemCollection) error {
				for _, c := range *col {
					iris = append(iris, c.GetLink())
				}
				return nil
			})
		} else {
			_ = vocab.OnObject(ob.Context, func(c *vocab.Object) error {
				iris = append(iris, c.GetLink())
				return nil
			})
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

type contextNil iriNil

func (c contextNil) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return len(accumContexts(it)) == 0
}

func accumAttributedTos(item vocab.Item) vocab.IRIs {
	iris := make(vocab.IRIs, 0)
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		if vocab.IsNil(ob.AttributedTo) {
			return nil
		}
		if vocab.IsIRI(ob.AttributedTo) {
			iris = append(iris, ob.AttributedTo.GetLink())
		} else if vocab.IsIRIs(ob.AttributedTo) {
			_ = vocab.OnIRIs(ob.AttributedTo, func(col *vocab.IRIs) error {
				for _, r := range *col {
					iris = append(iris, r.GetLink())
				}
				return nil
			})
		} else if vocab.IsItemCollection(ob.AttributedTo) {
			_ = vocab.OnItemCollection(ob.AttributedTo, func(attrTos *vocab.ItemCollection) error {
				for _, a := range *attrTos {
					iris = append(iris, a.GetLink())
				}
				return nil
			})
		} else {
			_ = vocab.OnObject(ob.AttributedTo, func(attrTo *vocab.Object) error {
				iris = append(iris, attrTo.GetLink())
				return nil
			})
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

type attributedToNil iriNil

func (a attributedToNil) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return len(accumAttributedTos(it)) == 0
}

func accumInReplyTos(item vocab.Item) vocab.IRIs {
	iris := make(vocab.IRIs, 0)
	_ = vocab.OnObject(item, func(ob *vocab.Object) error {
		if vocab.IsNil(ob.InReplyTo) {
			return nil
		}
		if vocab.IsIRI(ob.InReplyTo) {
			iris = append(iris, ob.InReplyTo.GetLink())
		} else if vocab.IsIRIs(ob.InReplyTo) {
			_ = vocab.OnIRIs(ob.InReplyTo, func(replTos *vocab.IRIs) error {
				for _, r := range *replTos {
					iris = append(iris, r.GetLink())
				}
				return nil
			})
		} else if vocab.IsItemCollection(ob.InReplyTo) {
			_ = vocab.OnItemCollection(ob.InReplyTo, func(replTos *vocab.ItemCollection) error {
				for _, r := range *replTos {
					iris = append(iris, r.GetLink())
				}
				return nil
			})
		} else {
			_ = vocab.OnObject(ob.InReplyTo, func(inReplyTo *vocab.Object) error {
				iris = append(iris, inReplyTo.GetLink())
				return nil
			})
		}
		return nil
	})
	return iris
}

var NilInReplyTo = inReplyToNil{}

type inReplyToNil iriNil

func (c inReplyToNil) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return true
	}
	return len(accumInReplyTos(it)) == 0
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
