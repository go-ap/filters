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

func FromIRI(i vocab.IRI) (Fns, error) {
	f := make(Fns, 0)
	u, err := i.URL()
	if err != nil {
		return nil, err
	}

	if u.User != nil {
		if us, err := url.Parse(u.User.Username()); err == nil {
			if id := vocab.IRI(us.String()); id != vocab.PublicNS {
				f = append(f, Authorized(id))
			}
		}
	}
	q := u.Query()

	if iri := q.Get("iri"); len(iri) > 0 {
		f = append(f, ID(vocab.IRI(iri)))
	}
	if iri := q.Get("id"); len(iri) > 0 {
		f = append(f, ID(vocab.IRI(iri)))
	}
	if maxItems, _ := strconv.ParseInt(q.Get("maxItems"), 10, 32); maxItems > 0 {
		f = append(f, WithMaxItems(int(maxItems)))
	}
	if typ, ok := q["type"]; ok && len(typ) > 0 {
		f = append(f, HasType(ActivityTypesFilter(typ...)...))
	}
	if names, ok := q["name"]; ok && len(names) > 0 {
		f = append(f, NameIn(names...))
	}

	return f, nil
}
