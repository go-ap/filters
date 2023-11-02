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

type Fn func(vocab.Item) bool

func Authorized(iri vocab.IRI) Fn {
	return func(it vocab.Item) bool {
		return fullAudience(it).Contains(iri)
	}
}

func (f Fn) Run(item vocab.Item) vocab.Item {
	if f != nil && !f(item) {
		return nil
	}
	return item
}

type Fns []Fn

func (ff Fns) Run(item vocab.Item) vocab.Item {
	if len(ff) == 0 {
		return item
	}
	if vocab.IsItemCollection(item) {
		_ = vocab.OnItemCollection(item, func(col *vocab.ItemCollection) error {
			item = ff.runOnItems(*col)
			return nil
		})
		return item
	}
	if !Any(ff...)(item) {
		return nil
	}
	return item
}

func (ff Fns) runOnItems(col vocab.ItemCollection) vocab.ItemCollection {
	result := make(vocab.ItemCollection, 0)
	for _, it := range col {
		if Any(ff...)(it) {
			result = append(result, it)
		}
	}
	return result
}

func VocabularyTypesFilter(types ...string) vocab.ActivityVocabularyTypes {
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

	return append(f, fromValues(u.Query())...)
}

func FromIRI(i vocab.IRI) (Fns, error) {
	if vocab.IsNil(i) {
		return nil, nil
	}
	u, err := i.URL()
	if err != nil {
		return nil, err
	}
	return FromURL(*u), nil
}

func FromValues(q url.Values) Fns {
	return fromValues(q)
}

func PaginationFromURL(u url.URL) Fns {
	q := u.Query()

	f := make(Fns, 0)
	if q.Has(keyBefore) {
		vv := q[keyBefore]
		if len(vv) > 0 {
			if _, err := url.ParseRequestURI(vv[0]); err == nil {
				f = append(f, Before(ID(vocab.IRI(vv[0]))))
			} else {
				f = append(f, Before(IDLike(vv[0])))
			}
		}
	}
	if q.Has(keyAfter) {
		vv := q[keyAfter]
		if len(vv) > 0 {
			if _, err := url.ParseRequestURI(vv[0]); err == nil {
				f = append(f, After(ID(vocab.IRI(vv[0]))))
			} else {
				f = append(f, After(IDLike(vv[0])))
			}
		}
	}
	if q.Has(keyMaxItems) {
		vv := q[keyMaxItems]
		if len(vv) > 0 {
			if maxItems, err := strconv.ParseInt(vv[0], 10, 32); err == nil {
				f = append(f, WithMaxCount(int(maxItems)))
			}
		}
	}
	return Fns{All(f...)}
}

func fromValues(q url.Values) Fns {

	actorQ := make(url.Values)
	objectQ := make(url.Values)
	targetQ := make(url.Values)

	f := make(Fns, 0)
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
		case keyType:
			f = append(f, HasType(VocabularyTypesFilter(vv...)...))
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
		f = append(f, Actor(fromValues(actorQ)...))
	}
	if len(objectQ) > 0 {
		f = append(f, Object(fromValues(objectQ)...))
	}
	if len(targetQ) > 0 {
		f = append(f, Target(fromValues(targetQ)...))
	}
	return f
}
