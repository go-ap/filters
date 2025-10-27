package filters

import (
	"net/url"
	"sort"
	"strconv"

	vocab "github.com/go-ap/activitypub"
)

// PaginateCollection is a function that populates the received collection
func PaginateCollection(it vocab.Item, filters ...Check) vocab.Item {
	if vocab.IsNil(it) || !it.IsCollection() {
		return it
	}

	total := uint(0)
	_ = vocab.OnCollectionIntf(it, func(col vocab.CollectionInterface) error {
		total = col.Count()
		return nil
	})

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
			if !nextIRI.GetLink().Equal(firstIRI) {
				c.Next = nextIRI
			}
			if !prevIRI.GetLink().Equal(firstIRI) {
				c.Prev = prevIRI
			}
			return nil
		})
	case vocab.CollectionType:
		_ = vocab.OnCollection(col, func(c *vocab.Collection) error {
			c.TotalItems = total
			c.First = firstIRI
			return nil
		})
	case vocab.CollectionPageType:
		_ = vocab.OnCollectionPage(col, func(c *vocab.CollectionPage) error {
			c.TotalItems = total
			c.PartOf = partOfIRI
			c.First = firstIRI
			if !nextIRI.GetLink().Equal(firstIRI) {
				c.Next = nextIRI
			}
			if !prevIRI.GetLink().Equal(firstIRI) {
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

func CursorFromItem(it vocab.Item, filters ...Check) (vocab.Item, vocab.Item, vocab.Item) {
	typ := it.GetType()

	if !vocab.CollectionTypes.Contains(typ) {
		return it, nil, nil
	}

	var prev url.Values
	var next url.Values

	var prevIRI vocab.IRI
	var nextIRI vocab.IRI

	shouldBePage := len(PaginationChecks(filters...)) > 0

	maxCount := MaxCount(filters...)
	switch typ {
	case vocab.OrderedCollectionPageType:
		_ = vocab.OnOrderedCollectionPage(it, func(new *vocab.OrderedCollectionPage) error {
			items := new.OrderedItems
			if maxCount < 0 && len(items) > MaxItems {
				filters = append(filters, WithMaxCount(MaxItems))
			}
			new.OrderedItems, prev, next = filterCollection(sortItemsByPublishedUpdated(items), filters...)
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
			items := new.Items
			if maxCount < 0 && len(items) > MaxItems {
				filters = append(filters, WithMaxCount(MaxItems))
			}
			new.Items, prev, next = filterCollection(items, filters...)
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
				items := new.OrderedItems
				if maxCount < 0 && len(items) > MaxItems {
					filters = append(filters, WithMaxCount(MaxItems))
				}
				new.OrderedItems, prev, next = filterCollection(sortItemsByPublishedUpdated(items), filters...)
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
				items := new.OrderedItems
				if maxCount < 0 && len(items) > MaxItems {
					filters = append(filters, WithMaxCount(MaxItems))
				}
				new.OrderedItems, prev, next = filterCollection(sortItemsByPublishedUpdated(items), filters...)
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
				items := new.Items
				if maxCount < 0 && len(items) > MaxItems {
					filters = append(filters, WithMaxCount(MaxItems))
				}
				new.Items, prev, next = filterCollection(items, filters...)
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
				items := new.Items
				if maxCount < 0 && len(items) > MaxItems {
					filters = append(filters, WithMaxCount(MaxItems))
				}
				new.Items, prev, next = filterCollection(items, filters...)
				if len(next) > 0 {
					new.First = getURL(it.GetLink(), next)
				}
				return nil
			})
		}
	case vocab.CollectionOfItems:
		_ = vocab.OnItemCollection(it, func(col *vocab.ItemCollection) error {
			items := *col
			if maxCount < 0 && len(items) > MaxItems {
				filters = append(filters, WithMaxCount(MaxItems))
			}
			it, prev, next = filterCollection(sortItemsByPublishedUpdated(items), filters...)
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

	filteredNotPaginated := FilterChecks(fns...).runOnItems(col)
	if len(filteredNotPaginated) == 0 {
		return filteredNotPaginated, pp, np
	}

	maxItems := MaxCount(fns...)
	if maxItems < 0 {
		maxItems = MaxItems
	}
	if maxItems == 0 {
		// NOTE(marius): this is a shortcut. We're assuming that if the calling code wants max 0 items in the
		// list, they're ok with circumventing the rest of filtering and receiving a hard 0 items collection.
		return vocab.ItemCollection{}, nil, nil
	}
	resetCounter(MaxCountCheck(fns...))
	resetAfter(fns...)
	resetBefore(fns...)

	result = PaginationChecks(fns...).runOnItems(filteredNotPaginated)
	if len(result) == 0 {
		return result, pp, np
	}
	onLastPage := len(AfterChecks(fns...)) > 0 && len(filteredNotPaginated) < maxItems
	onFirstPage := len(AfterChecks(fns...)) == 0 && filteredNotPaginated.First().GetLink().Equal(result.First().GetLink())

	var firstPage vocab.ItemCollection
	first := filteredNotPaginated.First()
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
		pp.Add(keyBefore, first.GetLink().String())
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

func sortItemsByPublishedUpdated(col vocab.ItemCollection) vocab.ItemCollection {
	sort.SliceStable(col, func(i, j int) bool {
		return vocab.ItemOrderTimestamp(col[i], col[j])
	})
	return col
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
