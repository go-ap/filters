package filters

import vocab "github.com/go-ap/activitypub"

type cursor Fns

// Cursor is an alias for running an All() aggregate filter function on the incoming fns functions
func Cursor(fns ...Fn) cursor {
	return cursor(fns)
}

func (c cursor) Run(item vocab.Item) vocab.Item {
	if len(c) == 0 {
		return item
	}

	if vocab.IsItemCollection(item) {
		_ = vocab.OnItemCollection(item, func(col *vocab.ItemCollection) error {
			item = c.runOnItems(*col)
			return nil
		})
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
