// Package filters contains helper functions to be used by the storage implementations for filtering out elements
// at load time.
package filters

import (
	"net/url"
	"strconv"
	"strings"

	vocab "github.com/go-ap/activitypub"
)

// Check represents an interface for a filter that can be applied on a [vocab.Item]
// and it returns true if it matches and false if it does not.
type Check interface {
	Match(vocab.Item) bool
}

// Checks aggregates a list of Check functions to be tested on a [vocab.Item].
type Checks []Check

func (ff Checks) Filter(item vocab.Item) vocab.Item {
	return FilterChecks(ff...).runOnItem(item)
}

func (ff Checks) Paginate(item vocab.Item) vocab.Item {
	return PaginateCollection(item, ff...)
}

func (ff Checks) Run(item vocab.Item) vocab.Item {
	if len(ff) == 0 || vocab.IsNil(item) {
		return item
	}

	if !item.IsCollection() {
		return FilterChecks(ff...).runOnItem(item)
	}

	_ = vocab.OnItemCollection(item, func(col *vocab.ItemCollection) error {
		if vocab.IsItemCollection(item) {
			item = FilterChecks(ff...).runOnItems(*col)
		} else {
			*col = FilterChecks(ff...).runOnItems(*col)
		}
		return nil
	})

	return PaginateCollection(item, ff...)
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
		return Check(ff[0]).Match
	}
	return All(ff...).Match
}

func (ff Checks) runOnItems(col vocab.ItemCollection) vocab.ItemCollection {
	if len(ff) == 0 {
		return col
	}
	result := make(vocab.ItemCollection, 0)
	for _, it := range col {
		if vocab.IsNil(it) || !checkFn(ff)(it) {
			continue
		}
		result = append(result, it)
	}

	return result
}

// VocabularyTypesFilter converts the received list of strings to a list of ActivityVocabularyTypes
// that can be used with the HasType filter function.
// The individual strings are not validated against the known vocabulary types.
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

	keyURL          = "url"
	keyAttributedTo = "attributedTo"
	keyInReplyTo    = "inReplyTo"
	keyContext      = "context"

	keyActor  = "actor"
	keyObject = "object"
	keyTarget = "target"

	keyTag = "tag"

	keyAfter  = "after"
	keyBefore = "before"

	keyMaxItems = "maxItems"
)

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
				f = append(f, Before(SameID(vocab.IRI(vv[0]))))
			} else {
				f = append(f, Before(IDLike(vv[0])))
			}
		}
	}
	if q.Has(keyAfter) {
		vv := q[keyAfter]
		if len(vv) > 0 {
			if _, err := url.ParseRequestURI(vv[0]); err == nil {
				f = append(f, After(SameID(vocab.IRI(vv[0]))))
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

type buildFilterFn func(string) Check

type checkGroup struct {
	nilFn  Check
	likeFn buildFilterFn
	sameFn buildFilterFn
}

func parseURLValue(v string) (string, string) {
	if len(v) < 1 {
		return opNone, v
	}
	op := opNone
	if len(v) > 2 && v[0:2] == opNotLike {
		op = v[0:2]
		v = v[2:]
	}
	if op == opNone && (v[0:1] == opNot || v[0:1] == opLike) {
		op = v[0:1]
		v = v[1:]
	}
	return op, v
}

const (
	opNot     = "!"
	opLike    = "~"
	opNotLike = "!~"
	opNone    = ""
	sNilIRI   = string(vocab.NilIRI)
	sEmptyIRI = string(vocab.EmptyIRI)
)

func (cg checkGroup) build(vv ...string) Check {
	f := make(Checks, 0)
	for _, n := range vv {
		switch op, v := parseURLValue(n); op {
		case opNone:
			if v == sNilIRI || v == sEmptyIRI {
				f = append(f, cg.nilFn)
			} else {
				f = append(f, cg.sameFn(n))
			}
		case opNot:
			if v == sNilIRI || v == sEmptyIRI {
				f = append(f, Not(cg.nilFn))
			} else {
				f = append(f, Not(cg.sameFn(v)))
			}
		case opLike:
			f = append(f, cg.likeFn(v))
		case opNotLike:
			f = append(f, Not(cg.likeFn(v)))
		}
	}
	if len(f) == 0 {
		return nil
	}
	if len(f) == 1 {
		return f[0]
	}
	return Any(f...)
}

var idFilters = checkGroup{
	nilFn:  NilID,
	likeFn: IDLike,
	sameFn: func(s string) Check {
		return SameID(vocab.IRI(s))
	},
}

var iriFilters = checkGroup{
	nilFn:  NilIRI,
	likeFn: IRILike,
	sameFn: func(s string) Check {
		return SameIRI(vocab.IRI(s))
	},
}

var nameFilters = checkGroup{
	nilFn:  NameEmpty,
	likeFn: NameLike,
	sameFn: NameIs,
}

var summaryFilters = checkGroup{
	nilFn:  SummaryEmpty,
	likeFn: SummaryLike,
	sameFn: SummaryIs,
}

var contentFilters = checkGroup{
	nilFn:  ContentEmpty,
	likeFn: ContentLike,
	sameFn: ContentIs,
}

var urlFilters = checkGroup{
	nilFn:  NilIRI,
	likeFn: URLLike,
	sameFn: func(s string) Check {
		return SameURL(vocab.IRI(s))
	},
}

var attributedToFilters = checkGroup{
	nilFn:  NilAttributedTo,
	likeFn: AttributedToLike,
	sameFn: func(s string) Check {
		return SameAttributedTo(vocab.IRI(s))
	},
}

var contextFilters = checkGroup{
	nilFn:  NilContext,
	likeFn: ContextLike,
	sameFn: func(s string) Check {
		return SameContext(vocab.IRI(s))
	},
}

var inReplyToFilters = checkGroup{
	nilFn:  NilInReplyTo,
	likeFn: InReplyToLike,
	sameFn: func(s string) Check {
		return SameInReplyTo(vocab.IRI(s))
	},
}

func fromValues(q url.Values) Checks {
	actorQ := make(url.Values)
	objectQ := make(url.Values)
	targetQ := make(url.Values)
	tagQ := make(url.Values)

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
			f = append(f, idFilters.build(vv...))
		case keyType:
			f = append(f, HasType(VocabularyTypesFilter(vv...)...))
		case keyName:
			f = append(f, nameFilters.build(vv...))
		case keySummary:
			f = append(f, summaryFilters.build(vv...))
		case keyContent:
			f = append(f, contentFilters.build(vv...))
		case keyActor:
			if len(remainder) == 0 {
				remainder = keyID
			}
			actorQ[remainder] = vv
		case keyObject:
			if len(remainder) == 0 {
				remainder = keyID
			}
			objectQ[remainder] = vv
		case keyTarget:
			if len(remainder) == 0 {
				remainder = keyID
			}
			targetQ[remainder] = vv
		case keyTag:
			if len(remainder) == 0 {
				remainder = keyID
			}
			tagQ[remainder] = vv
		case keyURL:
			f = append(f, urlFilters.build(vv...))
		case keyAttributedTo:
			f = append(f, attributedToFilters.build(vv...))
		case keyInReplyTo:
			f = append(f, inReplyToFilters.build(vv...))
		case keyContext:
			f = append(f, contextFilters.build(vv...))
		}
	}
	if len(actorQ) > 0 {
		if af := fromValues(actorQ); len(af) > 0 {
			f = append(f, Actor(af...))
		}
	}
	if len(objectQ) > 0 {
		if of := fromValues(objectQ); len(of) > 0 {
			f = append(f, Object(of...))
		}
	}
	if len(targetQ) > 0 {
		if tf := fromValues(targetQ); len(tf) > 0 {
			f = append(f, Target(tf...))
		}
	}
	if len(tagQ) > 0 {
		if tf := fromValues(tagQ); len(tf) > 0 {
			f = append(f, Tag(tf...))
		}
	}
	if len(f) == 0 {
		return nil
	}
	if len(f) == 1 {
		return f
	}
	return Checks{All(f...)}
}

func urlValue(f Check) url.Values {
	if f == nil {
		return nil
	}

	q := url.Values{}
	switch check := f.(type) {
	case iriEquals:
		q[keyIRI] = []string{string(check)}
	case idEquals:
		q[keyID] = []string{string(check)}
	case withTypes:
		for _, vv := range check {
			q.Add(keyType, string(vv))
		}
	case objectChecks:
		p := keyObject
		for kk, vv := range ToValues(check...) {
			q[p+"."+kk] = vv
		}
	case actorChecks:
		p := keyActor
		for kk, vv := range ToValues(check...) {
			q[p+"."+kk] = vv
		}
	case targetChecks:
		p := keyTarget
		for kk, vv := range ToValues(check...) {
			q[p+"."+kk] = vv
		}
	case tagChecks:
		p := keyTag
		for kk, vv := range ToValues(check...) {
			q[p+"."+kk] = vv
		}
	case *beforeCrit:
		if len(check.fns) >= 1 {
			if sameId, ok := check.fns[0].(iriEquals); ok {
				q.Set(keyBefore, string(sameId))
			}
			if sameId, ok := check.fns[0].(idEquals); ok {
				q.Set(keyBefore, string(sameId))
			}
		}
	case *afterCrit:
		if len(check.fns) >= 1 {
			if sameId, ok := check.fns[0].(iriEquals); ok {
				q.Set(keyAfter, string(sameId))
			}
			if sameId, ok := check.fns[0].(idEquals); ok {
				q.Set(keyAfter, string(sameId))
			}
		}
	case *counter:
		q.Set(keyMaxItems, strconv.FormatInt(int64(check.max), 10))
	}
	return q
}

func urlValues(ff ...Check) url.Values {
	q := url.Values{}
	for _, f := range ff {
		if qq := urlValue(f); len(qq) > 0 {
			for k, v := range qq {
				q[k] = v
			}
		}
	}
	return q
}

func ToValues(ff ...Check) url.Values {
	return urlValues(ff...)
}
