package filters

import (
	"net/url"
	"sort"
	"strconv"
	"time"

	vocab "github.com/go-ap/activitypub"
)

// PaginateCollection is a function that populates the received collection
func PaginateCollection(it vocab.Item, filters ...Check) vocab.Item {
	if vocab.IsNil(it) {
		return it
	}

	col, prevIRI, nextIRI := collectionPageFromItem(it, CursorChecks(filters...)...)
	if vocab.IsNil(col) {
		return it
	}
	if vocab.IsItemCollection(col) {
		return col
	}

	partOfIRI := it.GetID()
	firstIRI := partOfIRI
	if u, err := it.GetLink().URL(); err == nil {
		u.RawQuery = ""
		partOfIRI = vocab.IRI(u.String())
		q := make(url.Values)
		q.Set(keyMaxItems, strconv.Itoa(MaxItems))
		u.RawQuery = q.Encode()
		firstIRI = vocab.IRI(u.String())
	}

	switch col.GetType() {
	case vocab.OrderedCollectionType:
		_ = vocab.OnOrderedCollection(col, func(c *vocab.OrderedCollection) error {
			c.First = firstIRI
			return nil
		})
	case vocab.OrderedCollectionPageType:
		_ = vocab.OnOrderedCollectionPage(col, func(c *vocab.OrderedCollectionPage) error {
			c.PartOf = partOfIRI
			c.First = firstIRI
			if !nextIRI.GetLink().Equals(firstIRI, true) {
				c.Next = nextIRI
			}
			if !prevIRI.GetLink().Equals(firstIRI, true) {
				c.Prev = prevIRI
			}
			return nil
		})
	case vocab.CollectionType:
		_ = vocab.OnCollection(col, func(c *vocab.Collection) error {
			c.First = firstIRI
			return nil
		})
	case vocab.CollectionPageType:
		_ = vocab.OnCollectionPage(col, func(c *vocab.CollectionPage) error {
			c.PartOf = partOfIRI
			c.First = firstIRI
			if !nextIRI.GetLink().Equals(firstIRI, true) {
				c.Next = nextIRI
			}
			if !prevIRI.GetLink().Equals(firstIRI, true) {
				c.Prev = prevIRI
			}
			return nil
		})
	}

	return col
}

func getURL(i vocab.IRI, f url.Values) vocab.IRI {
	if f == nil {
		return i
	}
	if u, err := i.URL(); err == nil {
		u.RawQuery = f.Encode()
		i = vocab.IRI(u.String())
	}
	return i
}

func collectionPageFromItem(it vocab.Item, filters ...Check) (vocab.Item, vocab.Item, vocab.Item) {
	typ := it.GetType()

	if !vocab.CollectionTypes.Contains(typ) {
		return it, nil, nil
	}

	var prev url.Values
	var next url.Values

	var prevIRI vocab.IRI
	var nextIRI vocab.IRI

	shouldBePage := len(CursorChecks(filters...)) > 0

	switch typ {
	case vocab.OrderedCollectionPageType:
		_ = vocab.OnOrderedCollectionPage(it, func(new *vocab.OrderedCollectionPage) error {
			new.OrderedItems, prev, next = filterCollection(sortItemsByPublishedUpdated(new.Collection()), filters...)
			if len(prev) > 0 {
				prevIRI = getURL(it.GetLink(), prev)
			}
			if len(next) > 0 {
				nextIRI = getURL(it.GetLink(), next)
			}
			return nil
		})
	case vocab.CollectionPageType:
		_ = vocab.OnCollectionPage(it, func(new *vocab.CollectionPage) error {
			new.Items, prev, next = filterCollection(new.Collection(), filters...)
			if len(prev) > 0 {
				prevIRI = getURL(it.GetLink(), prev)
			}
			if len(next) > 0 {
				nextIRI = getURL(it.GetLink(), next)
			}
			return nil
		})
	case vocab.OrderedCollectionType:
		if shouldBePage {
			result := new(vocab.OrderedCollectionPage)
			old, _ := it.(*vocab.OrderedCollection)
			err := vocab.OnOrderedCollection(result, func(new *vocab.OrderedCollection) error {
				_, err := vocab.CopyOrderedCollectionProperties(new, old)
				new.Type = vocab.OrderedCollectionPageType
				new.OrderedItems, prev, next = filterCollection(sortItemsByPublishedUpdated(new.Collection()), filters...)
				if len(prev) > 0 {
					prevIRI = getURL(it.GetLink(), prev)
				}
				if len(next) > 0 {
					nextIRI = getURL(it.GetLink(), next)
				}
				return err
			})
			if err == nil {
				it = result
			}
		} else {
			_ = vocab.OnOrderedCollection(it, func(new *vocab.OrderedCollection) error {
				new.OrderedItems, prev, next = filterCollection(new.Collection(), filters...)
				if len(next) > 0 {
					new.First = getURL(it.GetLink(), next)
				}
				return nil
			})
		}
	case vocab.CollectionType:
		if shouldBePage {
			result := new(vocab.CollectionPage)
			old, _ := it.(*vocab.Collection)
			err := vocab.OnCollection(result, func(new *vocab.Collection) error {
				_, err := vocab.CopyCollectionProperties(new, old)
				new.Type = vocab.CollectionPageType
				new.Items, prev, next = filterCollection(new.Collection(), filters...)
				if len(prev) > 0 {
					prevIRI = getURL(it.GetLink(), prev)
				}
				if len(next) > 0 {
					nextIRI = getURL(it.GetLink(), next)
				}
				return err
			})
			if err == nil {
				it = result
			}
		} else {
			_ = vocab.OnCollection(it, func(new *vocab.Collection) error {
				new.Items, prev, next = filterCollection(new.Collection(), filters...)
				if len(next) > 0 {
					new.First = getURL(it.GetLink(), next)
				}
				return nil
			})
		}
	case vocab.CollectionOfItems:
		_ = vocab.OnItemCollection(it, func(col *vocab.ItemCollection) error {
			it, _, _ = filterCollection(sortItemsByPublishedUpdated(*col), filters...)
			return nil
		})
		return it, nil, nil
	}

	return it, prevIRI, nextIRI
}

func filterCollection(col vocab.ItemCollection, fns ...Check) (vocab.ItemCollection, url.Values, url.Values) {
	if len(col) == 0 {
		return col, nil, nil
	}

	pp := url.Values{}
	np := url.Values{}

	fpEnd := len(col) - 1
	if fpEnd > MaxItems {
		fpEnd = MaxItems
	}
	bpEnd := 0
	if len(col) > MaxItems {
		bpEnd = len(col) - MaxItems
	}

	firstPage := col[0:fpEnd]
	lastPage := col[len(col)-bpEnd:]

	result := Checks(fns).runOnItems(col)
	if len(result) == 0 {
		return result, pp, np
	}
	first := result.First()
	if len(col) > MaxItems {
		pp.Add(keyMaxItems, strconv.Itoa(MaxItems))
		np.Add(keyMaxItems, strconv.Itoa(MaxItems))

		onFirstPage := false
		for _, top := range firstPage {
			if onFirstPage = first.GetLink().Equals(top.GetLink(), true); onFirstPage {
				break
			}
		}
		if !onFirstPage {
			pp.Add(keyBefore, first.GetLink().String())
		} else {
			pp = nil
		}
		if len(result) > 1 && len(col) > MaxItems+1 {
			last := result[len(result)-1]
			onLastPage := false
			for _, bottom := range lastPage {
				if onLastPage = last.GetLink().Equals(bottom.GetLink(), true); onLastPage {
					break
				}
			}
			if !onLastPage {
				np.Add(keyAfter, last.GetLink().String())
			} else {
				np = nil
			}
		}
	}
	return result, pp, np
}

func sortItemsByPublishedUpdated(col vocab.ItemCollection) vocab.ItemCollection {
	sort.SliceStable(col, func(i int, j int) bool {
		it1 := col.Collection()[i]
		it2 := col.Collection()[j]
		var (
			p1 time.Time
			p2 time.Time
			u1 time.Time
			u2 time.Time
		)
		_ = vocab.OnObject(it1, func(ob *vocab.Object) error {
			p1 = ob.Published
			u1 = ob.Updated
			return nil
		})
		_ = vocab.OnObject(it2, func(ob *vocab.Object) error {
			p2 = ob.Published
			u2 = ob.Updated
			return nil
		})

		if d1 := u1.Sub(p1); d1 < 0 {
			u1 = p1
		}
		if d2 := u2.Sub(p2); d2 < 0 {
			u2 = p2
		}
		return u2.Sub(u1) < 0
	})
	return col
}

func isCursorFn(fn Check) bool {
	ok := false
	switch fn.(type) {
	case *afterCrit:
		ok = true
	case *beforeCrit:
		ok = true
	case *counter:
		ok = true
	}
	return ok
}

func isFilterFn(fn Check) bool {
	return !isCursorFn(fn)
}

func FilterChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		if isFilterFn(fn) {
			c = append(c, fn)
		}
	}
	return c
}

func CursorChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		if isCursorFn(fn) {
			c = append(c, fn)
		}
	}
	return c
}

func ObjectChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case authorized:
			c = append(c, fn)
		case objectChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func ActorChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case authorized:
			c = append(c, fn)
		case actorChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func TargetChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case authorized:
			c = append(c, fn)
		case targetChecks:
			c = append(c, fns...)
		}
	}
	return c
}

func ActivityChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case authorized:
			c = append(c, fn)
		case targetChecks:
			c = append(c, fns...)
		case objectChecks:
			c = append(c, fns...)
		case actorChecks:
			c = append(c, fns...)
		}
	}
	return c
}
func IntransitiveActivityChecks(fns ...Check) Checks {
	c := make([]Check, 0)
	for _, fn := range fns {
		switch fns := fn.(type) {
		case authorized:
			c = append(c, fn)
		case targetChecks:
			c = append(c, fns...)
		case actorChecks:
			c = append(c, fns...)
		}
	}
	return c
}
