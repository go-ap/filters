// Package filters contains helper functions to be used by the storage implementations for filtering out elements
// at load time.
package filters

import vocab "github.com/go-ap/activitypub"

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
