package filters

import (
	"net/url"
	"slices"
	"strconv"

	vocab "github.com/go-ap/activitypub"
)

func ResetPagination(fns ...Check) {
	resetCounter(MaxCountCheck(fns...))
	resetAfter(fns...)
	resetBefore(fns...)
}

// PaginateCollection is a function that populates the received collection
func PaginateCollection(it vocab.Item, filters ...Check) vocab.Item {
	if vocab.IsNil(it) || !vocab.IsCollection(it) {
		return it
	}

	col, prevIRI, nextIRI := CursorFromItem(it, filters...)
	if vocab.IsNil(col) {
		return it
	}
	if vocab.IsItemCollection(col) {
		return col
	}

	maxItems := MaxCount(filters...)
	if maxItems < 0 {
		maxItems = MaxItems
	}
	partOfIRI := it.GetID()
	firstIRI := partOfIRI
	if u, err := it.GetLink().URL(); err == nil {
		q := u.Query()
		for k := range q {
			if k == keyMaxItems || k == keyAfter || k == keyBefore {
				q.Del(k)
			}
		}
		u.RawQuery = q.Encode()
		partOfIRI = vocab.IRI(u.String())
		if !q.Has(keyMaxItems) {
			q.Set(keyMaxItems, strconv.Itoa(maxItems))
		}
		u.RawQuery = q.Encode()
		firstIRI = vocab.IRI(u.String())
	}

	typ := col.GetType()
	switch {
	case vocab.OrderedCollectionType.Match(typ), vocab.CollectionType.Match(typ):
		_ = vocab.OnOrderedCollection(col, func(c *vocab.OrderedCollection) error {
			c.First = firstIRI
			return nil
		})
	case vocab.OrderedCollectionPageType.Match(typ), vocab.CollectionPageType.Match(typ):
		_ = vocab.OnOrderedCollectionPage(col, func(c *vocab.OrderedCollectionPage) error {
			c.PartOf = partOfIRI
			c.First = firstIRI
			if !nextIRI.GetLink().Equal(vocab.EmptyIRI) && !nextIRI.GetLink().Equal(firstIRI) {
				c.Next = nextIRI
			}
			if !prevIRI.GetLink().Equal(vocab.EmptyIRI) && !prevIRI.GetLink().Equal(firstIRI) {
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
	_, hasAfter := f[keyAfter]
	_, hasBefore := f[keyBefore]
	if u, err := i.URL(); err == nil {
		q := u.Query()
		if hasAfter || hasBefore {
			q.Del(keyAfter)
			q.Del(keyBefore)
		}
		for k, v := range f {
			q[k] = v
		}
		u.RawQuery = q.Encode()
		i = vocab.IRI(u.String())
	}
	return i
}

func getCollectionProperty(it vocab.CollectionInterface, colFn func(*vocab.IRI, *vocab.Collection) error, pageFn func(*vocab.IRI, *vocab.CollectionPage) error) vocab.IRI {
	iri := vocab.EmptyIRI
	if vocab.IsNil(it) {
		return iri
	}

	// NOTE(marius): we don't need to mess with the item's type
	//  Additionally the OrderedCollection is compatible with the memory layout of the Collection
	//  so we can use a single branch here.
	_ = vocab.OnCollection(it, func(c *vocab.Collection) error {
		return colFn(&iri, c)
	})
	_ = vocab.OnCollectionPage(it, func(p *vocab.CollectionPage) error {
		return pageFn(&iri, p)
	})
	return iri
}

func NextPageFromCollection(it vocab.CollectionInterface) vocab.IRI {
	nextColFn := func(iri *vocab.IRI, c *vocab.Collection) error {
		if !vocab.IsNil(c.First) {
			*iri = c.First.GetLink()
		}
		return nil
	}
	nextPageFn := func(iri *vocab.IRI, c *vocab.CollectionPage) error {
		if !vocab.IsNil(c.Next) {
			*iri = c.Next.GetLink()
		}
		return nil
	}
	return getCollectionProperty(it, nextColFn, nextPageFn)
}

func PrevPageFromCollection(it vocab.CollectionInterface) vocab.IRI {
	prefColFn := func(iri *vocab.IRI, c *vocab.Collection) error {
		if !vocab.IsNil(c.First) {
			*iri = c.First.GetLink()
		}
		return nil
	}
	prevPageFn := func(iri *vocab.IRI, c *vocab.CollectionPage) error {
		if !vocab.IsNil(c.Prev) {
			*iri = c.Prev.GetLink()
		}
		return nil
	}
	return getCollectionProperty(it, prefColFn, prevPageFn)
}

func CursorFromItem(it vocab.Item, filters ...Check) (vocab.Item, vocab.Item, vocab.Item) {
	typ := it.GetType()

	if !vocab.CollectionTypes.Match(typ) {
		return it, nil, nil
	}

	var prev url.Values
	var next url.Values

	var prevIRI vocab.IRI
	var nextIRI vocab.IRI

	shouldBePage := len(PaginationChecks(filters...)) > 0

	switch {
	case vocab.OrderedCollectionPageType.Match(typ):
		_ = vocab.OnOrderedCollectionPage(it, func(new *vocab.OrderedCollectionPage) error {
			new.ID = IRIf(new.ID, filters...)
			items := new.OrderedItems
			slices.SortStableFunc(items, vocab.TimestampSortFunc)
			new.OrderedItems, prev, next = filterCollection(items, filters...)
			if len(prev) > 0 {
				prevIRI = getURL(it.GetLink(), prev)
			}
			if len(next) > 0 {
				nextIRI = getURL(it.GetLink(), next)
			}
			return nil
		})
	case vocab.CollectionPageType.Match(typ):
		_ = vocab.OnCollectionPage(it, func(new *vocab.CollectionPage) error {
			new.ID = IRIf(new.ID, filters...)
			items := new.Items
			new.Items, prev, next = filterCollection(items, filters...)
			if len(prev) > 0 {
				prevIRI = getURL(it.GetLink(), prev)
			}
			if len(next) > 0 {
				nextIRI = getURL(it.GetLink(), next)
			}
			return nil
		})
	case vocab.OrderedCollectionType.Match(typ):
		if shouldBePage {
			result := new(vocab.OrderedCollectionPage)
			err := vocab.OnOrderedCollection(it, func(old *vocab.OrderedCollection) error {
				return vocab.OnOrderedCollection(result, func(new *vocab.OrderedCollection) error {
					_, err := vocab.CopyOrderedCollectionProperties(new, old)
					new.ID = IRIf(new.ID, filters...)
					new.Type = vocab.OrderedCollectionPageType
					items := new.OrderedItems
					slices.SortStableFunc(items, vocab.TimestampSortFunc)
					new.OrderedItems, prev, next = filterCollection(items, filters...)
					if len(prev) > 0 {
						prevIRI = getURL(it.GetLink(), prev)
					}
					if len(next) > 0 {
						nextIRI = getURL(it.GetLink(), next)
					}
					return err
				})
			})
			if err == nil {
				it = result
			}
		} else {
			_ = vocab.OnOrderedCollection(it, func(new *vocab.OrderedCollection) error {
				new.ID = IRIf(new.ID, filters...)
				items := new.OrderedItems
				slices.SortStableFunc(items, vocab.TimestampSortFunc)
				new.OrderedItems, prev, next = filterCollection(items, filters...)
				if len(next) > 0 {
					new.First = getURL(it.GetLink(), next)
				}
				return nil
			})
		}
	case vocab.CollectionType.Match(typ):
		if shouldBePage {
			result := new(vocab.CollectionPage)
			err := vocab.OnCollection(it, func(old *vocab.Collection) error {
				return vocab.OnCollection(result, func(new *vocab.Collection) error {
					_, err := vocab.CopyCollectionProperties(new, old)
					new.ID = IRIf(new.ID, filters...)
					new.Type = vocab.CollectionPageType
					items := new.Items
					new.Items, prev, next = filterCollection(items, filters...)
					if len(prev) > 0 {
						prevIRI = getURL(it.GetLink(), prev)
					}
					if len(next) > 0 {
						nextIRI = getURL(it.GetLink(), next)
					}
					return err
				})
			})
			if err == nil {
				it = result
			}
		} else {
			_ = vocab.OnCollection(it, func(new *vocab.Collection) error {
				new.ID = IRIf(new.ID, filters...)
				items := new.Items
				new.Items, prev, next = filterCollection(items, filters...)
				if len(next) > 0 {
					new.First = getURL(it.GetLink(), next)
				}
				return nil
			})
		}
	case vocab.CollectionOfItems.Match(typ):
		_ = vocab.OnItemCollection(it, func(col *vocab.ItemCollection) error {
			items := *col
			slices.SortStableFunc(items, vocab.TimestampSortFunc)
			it, prev, next = filterCollection(items, filters...)
			if len(prev) > 0 {
				prevIRI = getURL(it.GetLink(), prev)
			}
			if len(next) > 0 {
				nextIRI = getURL(it.GetLink(), next)
			}
			return nil
		})
	}

	return it, prevIRI, nextIRI
}

func resetAfter(fns ...Check) {
	for _, fn := range fns {
		if af, ok := fn.(*afterCrit); ok {
			af.check = false
		}
	}
}

func resetBefore(fns ...Check) {
	for _, fn := range fns {
		if bf, ok := fn.(*beforeCrit); ok {
			bf.check = true
		}
	}
}

func resetCounter(fn Check) {
	if mit, ok := fn.(*counter); ok {
		mit.cnt = 0
	}
}

func filterCollection(col vocab.ItemCollection, fns ...Check) (vocab.ItemCollection, url.Values, url.Values) {
	if len(col) == 0 {
		return col, nil, nil
	}

	pp := url.Values{}
	np := url.Values{}

	var lastPage vocab.ItemCollection
	var result vocab.ItemCollection

	if ff := FilterChecks(fns...); len(ff) > 0 {
		if col = ff.runOnItems(col); len(col) == 0 {
			return col, pp, np
		}
	}

	maxItems := MaxCount(fns...)
	if maxItems < 0 {
		maxItems = MaxItems
	}
	if maxItems == 0 {
		// NOTE(marius): this is a shortcut. We're assuming that if the calling code wants max 0 items in the
		//  list, they're ok with circumventing the rest of filtering and receiving a hard 0 items collection.
		return vocab.ItemCollection{}, nil, nil
	}
	resetCounter(MaxCountCheck(fns...))
	resetAfter(fns...)
	resetBefore(fns...)

	result = col
	after := AfterChecks(fns...)
	if len(after) > 0 {
		result = Checks{After(after...)}.runOnItems(result)
	}
	// NOTE(marius): checking a before filter is equivalent to reversing the result and using after
	before := BeforeChecks(fns...)
	if len(before) > 0 {
		slices.Reverse(result)
		result = Checks{After(before...)}.runOnItems(result)
	}
	if maxCountCheck := MaxCountCheck(fns...); maxCountCheck != nil {
		result = Checks{maxCountCheck}.runOnItems(result)
	}
	if len(result) == 0 {
		return result, pp, np
	}
	if len(before) > 0 {
		slices.Reverse(result)
	}
	onLastPage := len(after) > 0 && len(col) < maxItems
	onFirstPage := len(after) == 0 && col.First().GetLink().Equal(result.First().GetLink())

	var firstPage vocab.ItemCollection
	first := col.First()
	if len(col) <= maxItems {
		return result, pp, np
	}

	pp.Add(keyMaxItems, strconv.Itoa(maxItems))
	np.Add(keyMaxItems, strconv.Itoa(maxItems))

	for _, top := range firstPage {
		if onFirstPage = first.GetLink().Equal(top.GetLink()); onFirstPage {
			break
		}
	}
	if !onFirstPage {
		prev := result[0]
		pp.Add(keyBefore, prev.GetLink().String())
	} else {
		pp = nil
	}
	if len(result) >= 1 && len(col) > maxItems+1 {
		last := result[len(result)-1]
		for _, bottom := range lastPage {
			if onLastPage = last.GetLink().Equal(bottom.GetLink()); onLastPage {
				break
			}
		}
		if !onLastPage {
			np.Add(keyAfter, last.GetLink().String())
		} else {
			np = nil
		}
	}
	return result, pp, np
}

func isCounterFn(fn Check) bool {
	_, ok := fn.(*counter)
	return ok
}

func isCursorFn(fn Check) bool {
	ok := false
	switch fn.(type) {
	case *afterCrit:
		ok = true
	case *beforeCrit:
		ok = true
	}
	return ok
}

func isFilterFn(fn Check) bool {
	return !(isCursorFn(fn) || isCounterFn(fn))
}
