package filters

import (
	vocab "github.com/go-ap/activitypub"
)

type counter int32

// WithMaxCount is used to limit a collection's items count to the 'max' value.
// It can be used from slicing from the first element of the collection to max.
// Due to relying on the static max value the function is not reentrant.
func WithMaxCount(max int) Fn {
	cnt := counter(0)
	return cnt.onReachMax(max)
}

func (cnt *counter) onReachMax(max int) Fn {
	return func(item vocab.Item) bool {
		if int32(max) <= int32(*cnt) {
			return false
		}
		*cnt = *cnt + 1
		return true
	}
}

func WithMaxItems(max int) Fn {
	var OrderedCollectionTypes = vocab.ActivityVocabularyTypes{vocab.OrderedCollectionType, vocab.OrderedCollectionPageType}
	var CollectionTypes = vocab.ActivityVocabularyTypes{vocab.CollectionType, vocab.CollectionPageType}

	return func(it vocab.Item) bool {
		if vocab.IsItemCollection(it) {
			_ = vocab.OnItemCollection(it, func(col *vocab.ItemCollection) error {
				if max < len(*col) {
					*col = (*col)[0:max]
				}
				return nil
			})
		}
		if OrderedCollectionTypes.Contains(it.GetType()) {
			_ = vocab.OnOrderedCollection(it, func(col *vocab.OrderedCollection) error {
				if max < len(col.OrderedItems) {
					col.OrderedItems = col.OrderedItems[0:max]
				}
				return nil
			})
		}
		if CollectionTypes.Contains(it.GetType()) {
			_ = vocab.OnCollection(it, func(col *vocab.Collection) error {
				if max < len(col.Items) {
					col.Items = col.Items[0:max]
				}
				return nil
			})
		}
		return true
	}
}

// After checks the activitypub.Item against a specified "fn" filter function.
// This should be used when iterating over a collection, and it resolves to true
// after fn returns true and to false before.
//
// Due to relying on the static check function return value the After is not reentrant.
func After(fn Fn) Fn {
	p := pagination(false)
	return p.after(fn)
}

func (isAfter *pagination) after(fn Fn) Fn {
	return func(it vocab.Item) bool {
		if vocab.IsNil(it) {
			return bool(*isAfter)
		}
		if fn(it) {
			*isAfter = true
			return false
		}
		return bool(*isAfter)
	}
}

type pagination bool

// Before checks the activitypub.Item against a specified "fn" filter function.
// This should be used when iterating over a collection, and it resolves to true before
// the fn has returned true and to false after.
//
// Due to relying on the static check function return value the function is not reentrant.
func Before(fn Fn) Fn {
	p := pagination(true)
	return p.before(fn)
}

func (isBefore *pagination) before(fn Fn) Fn {
	return func(it vocab.Item) bool {
		if vocab.IsNil(it) {
			return true
		}
		if fn(it) {
			*isBefore = false
		}
		return bool(*isBefore)
	}
}
