// Package filters contains helper functions to be used by the storage implementations for filtering out elements
// at load time.
/*
Example:

	s, err := fs.New()
	if err != nil {
		// error handling
	}
	// This searches for all Create activities published by the Actor with the
	// ID https://example.com/authors/jdoe, or with the name "JohnDoe" and, which
	// have an object with a non nil ID.
	collectionItem, err := s.Load("https://example.com/outbox", All(
		Type("Create"),
		Actor(Any(ID("https://example.com/authors/jdoe"), NameIs("JohnDoe")),
		Object(NotNilID),
	)))
*/
package filters

import (
	"net/url"
	"strconv"
	"strings"

	vocab "github.com/go-ap/activitypub"
)

type Fn = func(vocab.Item) bool

func Authorized(iri vocab.IRI) Fn {
	return func(it vocab.Item) bool {
		if r, ok := it.(vocab.HasRecipients); ok {
			return r.Recipients().Contains(iri)
		}
		return false
	}
}

type Fns []Fn

func ActivityTypesFilter(types ...string) vocab.ActivityVocabularyTypes {
	r := make(vocab.ActivityVocabularyTypes, 0, len(types))
	for _, t := range types {
		typ := vocab.ActivityVocabularyType(t)
		if vocab.Types.Contains(typ) {
			r = append(r, typ)
		}
	}
	return r
}

const (
	keyID   = "id"
	keyIRI  = "iri"
	keyType = "type"

	keyName    = "name"
	keySummary = "summary"
	keyContent = "content"

	keyActor  = "actor"
	keyObject = "object"
	keyTarget = "target"

	keyAfter  = "after"
	keyBefore = "before"

	keyMaxItems = "maxItems"
)

func ids(vv []string) []Fn {
	f := make([]Fn, 0)
	for _, v := range vv {
		if v == "" {
			f = append(f, NilID)
		} else if v == "!" || v == "!-" {
			f = append(f, Not(NilID))
		} else if strings.HasPrefix(v, "!") {
			f = append(f, Not(ID(vocab.IRI(v[1:]))))
		} else if strings.HasPrefix(v, "~") {
			f = append(f, IDLike(v[1:]))
		} else {
			f = append(f, ID(vocab.IRI(v)))
		}
	}
	return f
}

func FromURL(u url.URL) Fns {
	f := make(Fns, 0)

	if u.User != nil {
		if us, err := url.ParseRequestURI(u.User.Username()); err == nil {
			if id := vocab.IRI(us.String()); id != vocab.PublicNS {
				f = append(f, Authorized(id))
			}
		}
	}

	return append(f, FromValues(u.Query())...)
}

func FromIRI(i vocab.IRI) (Fns, error) {
	u, err := i.URL()
	if err != nil {
		return nil, err
	}
	return FromURL(*u), nil
}

func FromValues(q url.Values) Fns {
	f := make(Fns, 0)
	actorQ := make(url.Values)
	objectQ := make(url.Values)
	targetQ := make(url.Values)
	for k, vv := range q {
		pieces := strings.SplitN(k, ".", 2)
		piece := k
		remainder := ""
		if len(pieces) > 1 {
			piece = pieces[0]
			remainder = pieces[1]
		}
		switch piece {
		case keyID, keyIRI:
			f = append(f, ids(vv)...)
		case keyMaxItems:
			if maxItems, _ := strconv.ParseInt(q.Get(keyMaxItems), 10, 32); maxItems > 0 {
				f = append(f, WithMaxItems(int(maxItems)))
			}
		case keyType:
			f = append(f, HasType(ActivityTypesFilter(vv...)...))
		case keyName:
			for _, n := range vv {
				if n == "" {
					f = append(f, NameEmpty())
				} else if n == "!" || n == "!-" {
					f = append(f, Not(NameEmpty()))
				} else if strings.HasPrefix(n, "!") {
					f = append(f, Not(NameLike(n[1:])))
				} else if strings.HasPrefix(n, "~") {
					f = append(f, NameLike(n[1:]))
				} else {
					f = append(f, NameIs(n))
				}
			}
		case keySummary:
			for _, n := range vv {
				if n == "" {
					f = append(f, SummaryEmpty())
				} else if n == "!" || n == "!-" {
					f = append(f, Not(SummaryEmpty()))
				} else if strings.HasPrefix(n, "!") {
					f = append(f, Not(SummaryLike(n[1:])))
				} else if strings.HasPrefix(n, "~") {
					f = append(f, SummaryLike(n[1:]))
				} else {
					f = append(f, SummaryIs(n))
				}
			}
		case keyContent:
			for _, n := range vv {
				if n == "" {
					f = append(f, ContentEmpty())
				} else if n == "!" || n == "!-" {
					f = append(f, Not(ContentEmpty()))
				} else if strings.HasPrefix(n, "!") && n[1] != '-' {
					f = append(f, Not(ContentLike(n[1:])))
				} else if strings.HasPrefix(n, "~") {
					f = append(f, ContentLike(n[1:]))
				} else {
					f = append(f, ContentIs(n))
				}
			}
		case keyAfter:
			if len(vv) > 0 {
				if _, err := url.ParseRequestURI(vv[0]); err == nil {
					f = append(f, After(ID(vocab.IRI(vv[0]))))
				} else {
					f = append(f, After(IDLike(vv[0])))
				}
			}
		case keyBefore:
			if len(vv) > 0 {
				if _, err := url.ParseRequestURI(vv[0]); err == nil {
					f = append(f, Before(ID(vocab.IRI(vv[0]))))
				} else {
					f = append(f, Before(IDLike(vv[0])))
				}
			}
		case keyActor:
			if len(remainder) > 0 {
				actorQ[remainder] = vv
			}
		case keyObject:
			if len(remainder) > 0 {
				objectQ[remainder] = vv
			}
		case keyTarget:
			if len(remainder) > 0 {
				targetQ[remainder] = vv
			}
		}
	}
	if len(actorQ) > 0 {
		f = append(f, Actor(FromValues(actorQ)...))
	}
	if len(objectQ) > 0 {
		f = append(f, Object(FromValues(objectQ)...))
	}
	if len(targetQ) > 0 {
		f = append(f, Target(FromValues(targetQ)...))
	}
	return f
}
