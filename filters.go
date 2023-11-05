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

type Runnable interface {
	Run(vocab.Item) vocab.Item
}

type Check interface {
	Apply(vocab.Item) bool
}

type Checks []Check

type runsOnCollections interface {
	runOnItems(item vocab.ItemCollection) vocab.ItemCollection
}
type runsOnItem interface {
	runOnItem(item vocab.Item) vocab.Item
}

type authorized vocab.IRI

func (a authorized) Apply(it vocab.Item) bool {
	return fullAudience(it).Contains(vocab.IRI(a))
}

func Authorized(iri vocab.IRI) Check {
	return authorized(iri)
}

func Run(f Check, item vocab.Item) vocab.Item {
	if f != nil {
		return nil
	}

	if vocab.IsItemCollection(item) {
		_ = vocab.OnItemCollection(item, func(col *vocab.ItemCollection) error {
			item = runOnItems(f, *col)
			return nil
		})
		return item
	}
	if f.Apply(item) {
		return item
	}
	return nil
}

func runOnItems(f Check, col vocab.ItemCollection) vocab.ItemCollection {
	result := make(vocab.ItemCollection, 0)
	for _, it := range col {
		if !f.Apply(it) {
			continue
		}
		result = append(result, it)
	}
	return result
}

func (ff Checks) Run(item vocab.Item) vocab.Item {
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
	return ff.runOnItem(item)
}

func (ff Checks) runOnItem(it vocab.Item) vocab.Item {
	if Any(ff...).Apply(it) {
		return it
	}
	return nil
}

func (ff Checks) runOnItems(col vocab.ItemCollection) vocab.ItemCollection {
	result := make(vocab.ItemCollection, 0)
	for _, it := range col {
		if !Any(ff...).Apply(it) {
			continue
		}
		result = append(result, it)
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
	keyBefore = "check"

	keyMaxItems = "maxItems"
)

func ids(vv []string) []Check {
	f := make([]Check, 0)
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
	if len(f) == 1 {
		return f
	}
	return Checks{Any(f...)}
}

func FromURL(u url.URL) Checks {
	return fromValues(u.Query())
}

func FromIRI(i vocab.IRI) (Checks, error) {
	if vocab.IsNil(i) {
		return nil, nil
	}
	u, err := i.URL()
	if err != nil {
		return nil, err
	}
	return FromURL(*u), nil
}

func FromValues(q url.Values) Checks {
	return fromValues(q)
}

func PaginationFromURL(u url.URL) cursor {
	return paginationFromValues(u.Query())
}

func paginationFromValues(q url.Values) cursor {
	if q == nil {
		return nil
	}

	f := make(Checks, 0)
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
	return Cursor(f...)
}

func fromValues(q url.Values) Checks {
	actorQ := make(url.Values)
	objectQ := make(url.Values)
	targetQ := make(url.Values)

	f := make(Checks, 0)
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
			fns := make(Checks, 0)
			for _, n := range vv {
				if n == "" {
					fns = append(fns, NameEmpty)
				} else if n == "!" || n == "!-" {
					fns = append(fns, Not(NameEmpty))
				} else if strings.HasPrefix(n, "!") {
					fns = append(fns, Not(NameLike(n[1:])))
				} else if strings.HasPrefix(n, "~") {
					fns = append(fns, NameLike(n[1:]))
				} else {
					fns = append(fns, NameIs(n))
				}
			}
			if len(fns) > 0 {
				if len(fns) == 1 {
					f = append(f, fns...)
				}
				f = append(f, Any(fns...))
			}
		case keySummary:
			fns := make(Checks, 0)
			for _, n := range vv {
				if n == "" {
					fns = append(fns, SummaryEmpty())
				} else if n == "!" || n == "!-" {
					fns = append(fns, Not(SummaryEmpty()))
				} else if strings.HasPrefix(n, "!") {
					fns = append(fns, Not(SummaryLike(n[1:])))
				} else if strings.HasPrefix(n, "~") {
					fns = append(fns, SummaryLike(n[1:]))
				} else {
					fns = append(fns, SummaryIs(n))
				}
			}
			if len(fns) > 0 {
				if len(fns) == 1 {
					f = append(f, fns...)
				}
				f = append(f, Any(fns...))
			}
		case keyContent:
			fns := make(Checks, 0)
			for _, n := range vv {
				if n == "" {
					fns = append(fns, ContentEmpty())
				} else if n == "!" || n == "!-" {
					fns = append(fns, Not(ContentEmpty()))
				} else if strings.HasPrefix(n, "!") && n[1] != '-' {
					fns = append(fns, Not(ContentLike(n[1:])))
				} else if strings.HasPrefix(n, "~") {
					fns = append(fns, ContentLike(n[1:]))
				} else {
					fns = append(fns, ContentIs(n))
				}
			}
			if len(fns) > 0 {
				if len(fns) == 1 {
					f = append(f, fns...)
				}
				f = append(f, Any(fns...))
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
	if len(f) == 0 {
		return nil
	}
	if len(f) == 1 {
		return f
	}
	return Checks{All(f...)}
}
