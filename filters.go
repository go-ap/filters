// Package filters contains helper functions to be used by the storage implementations for filtering out elements
// at load time.
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

var orderedCollectionTypes = vocab.ActivityVocabularyTypes{
	vocab.OrderedCollectionPageType,
	vocab.OrderedCollectionType,
}
var collectionTypes = vocab.ActivityVocabularyTypes{
	vocab.CollectionPageType,
	vocab.CollectionType,
}

func (ff Checks) Run(item vocab.Item) vocab.Item {
	if len(ff) == 0 || vocab.IsNil(item) {
		return item
	}

	if item.IsCollection() {
		_ = vocab.OnItemCollection(item, func(col *vocab.ItemCollection) error {
			if vocab.IsItemCollection(item) {
				item = FilterChecks(ff...).runOnItems(*col)
			} else {
				*col = FilterChecks(ff...).runOnItems(*col)
			}
			return nil
		})

		switch item.GetType() {
		case vocab.OrderedCollectionType:
			_ = vocab.OnOrderedCollection(item, func(c *vocab.OrderedCollection) error {
				c.TotalItems = c.Count()
				return nil
			})
		case vocab.OrderedCollectionPageType:
			_ = vocab.OnOrderedCollectionPage(item, func(c *vocab.OrderedCollectionPage) error {
				c.TotalItems = c.Count()
				return nil
			})
		case vocab.CollectionType:
			_ = vocab.OnCollection(item, func(c *vocab.Collection) error {
				c.TotalItems = c.Count()
				return nil
			})
		case vocab.CollectionPageType:
			_ = vocab.OnCollectionPage(item, func(c *vocab.CollectionPage) error {
				c.TotalItems = c.Count()
				return nil
			})
		}
		return PaginateCollection(item, ff...)
	}
	return FilterChecks(ff...).runOnItem(item)
}

func (ff Checks) runOnItem(it vocab.Item) vocab.Item {
	if checkFn(ff)(it) {
		return it
	}
	return nil
}

func checkFn(ff Checks) func(vocab.Item) bool {
	if len(ff) == 0 {
		return func(_ vocab.Item) bool {
			return true
		}
	}
	if len(ff) == 1 && ff[0] != nil {
		return Check(ff[0]).Apply
	}
	return All(ff...).Apply
}

func (ff Checks) runOnItems(col vocab.ItemCollection) vocab.ItemCollection {
	if len(ff) == 0 {
		return col
	}
	result := make(vocab.ItemCollection, 0)
	for _, it := range col {
		if !checkFn(ff)(it) {
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

	keyURL = "url"

	keyActor  = "actor"
	keyObject = "object"
	keyTarget = "target"

	keyAfter  = "after"
	keyBefore = "before"

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
	q := u.Query()
	return append(fromValues(q), paginationFromValues(q)...)
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
	return append(fromValues(q), paginationFromValues(q)...)
}

func paginationFromValues(q url.Values) Checks {
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
	return f
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
				} else {
					f = append(f, Any(fns...))
				}
			}
		case keySummary:
			fns := make(Checks, 0)
			for _, n := range vv {
				if n == "" {
					fns = append(fns, SummaryEmpty)
				} else if n == "!" || n == "!-" {
					fns = append(fns, Not(SummaryEmpty))
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
				} else {
					f = append(f, Any(fns...))
				}
			}
		case keyContent:
			fns := make(Checks, 0)
			for _, n := range vv {
				if n == "" {
					fns = append(fns, ContentEmpty)
				} else if n == "!" || n == "!-" {
					fns = append(fns, Not(ContentEmpty))
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
				} else {
					f = append(f, Any(fns...))
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
		case keyURL:
			fns := make(Checks, 0)
			for _, n := range vv {
				if n == "" {
					f = append(f, NilURL)
				} else if n == "!" || n == "!-" {
					f = append(f, Not(NilURL))
				} else if strings.HasPrefix(n, "!") {
					f = append(f, Not(SameURL(vocab.IRI(n[1:]))))
				} else if strings.HasPrefix(n, "~") {
					f = append(f, URLLike(n[1:]))
				} else {
					f = append(f, SameURL(vocab.IRI(n)))
				}
				f = append(f)
			}
			if len(fns) > 0 {
				if len(fns) == 1 {
					f = append(f, fns...)
				} else {
					f = append(f, Any(fns...))
				}
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
