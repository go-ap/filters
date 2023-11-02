package filters

import vocab "github.com/go-ap/activitypub"

type cursor Fns

// Cursor is an alias for running an All() aggregate filter function on the incoming fns functions
func Cursor(fns ...Fn) cursor {
	return cursor(fns)
}

var collectionTypes = vocab.ActivityVocabularyTypes{
	vocab.CollectionType,
	vocab.CollectionPageType,
}

var orderedCollectionTypes = vocab.ActivityVocabularyTypes{
	vocab.OrderedCollectionType,
	vocab.OrderedCollectionPageType,
}

func (c cursor) Run(item vocab.Item) vocab.Item {
	if len(c) == 0 {
		return item
	}

	// NOTE(marius): here is the place where we should add pagination iris to the
	// collections, we can compute first/before/after
	if item.IsCollection() {
		var items vocab.ItemCollection
		var cnt uint
		_ = vocab.OnCollectionIntf(item, func(col vocab.CollectionInterface) error {
			cnt = col.Count()
			items = c.runOnItems(col.Collection())
			return nil
		})
		switch {
		case collectionTypes.Contains(item.GetType()):
			_ = vocab.OnCollection(item, func(col *vocab.Collection) error {
				col.Items = items
				col.TotalItems = cnt
				return nil
			})
		case orderedCollectionTypes.Contains(item.GetType()):
			_ = vocab.OnOrderedCollection(item, func(col *vocab.OrderedCollection) error {
				col.OrderedItems = items
				col.TotalItems = cnt
				return nil
			})
		case vocab.CollectionOfItems == item.GetType():
			item = &items
		}
		return item
	}

	return c.runOnItem(item)
}

func (c cursor) runOnItem(it vocab.Item) vocab.Item {
	for _, fn := range c {
		if fn == nil {
			continue
		}
		if !fn(it) {
			return nil
		}
	}
	return it
}

func (c cursor) runOnItems(col vocab.ItemCollection) vocab.ItemCollection {
	result := make(vocab.ItemCollection, 0)
	for _, it := range col {
		if it = c.runOnItem(it); vocab.IsNil(it) {
			continue
		}
		result = append(result, it)
	}
	return result
}
