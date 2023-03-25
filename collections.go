package filters

import vocab "github.com/go-ap/activitypub"

func WithMaxCount(max int) Fn {
	count := 0
	return func(item vocab.Item) bool {
		if count >= max {
			return false
		}
		count += 1
		return true
	}
}

func WithMaxItems(max int) Fn {
	var OrderedCollectionTypes = vocab.ActivityVocabularyTypes{vocab.OrderedCollectionType, vocab.OrderedCollectionPageType}
	var CollectionTypes = vocab.ActivityVocabularyTypes{vocab.CollectionType, vocab.CollectionPageType}

	return func(it vocab.Item) bool {
		if vocab.IsItemCollection(it) {
			vocab.OnItemCollection(it, func(col *vocab.ItemCollection) error {
				if max < len(*col) {
					*col = (*col)[0:max]
				}
				return nil
			})
		}
		if OrderedCollectionTypes.Contains(it.GetType()) {
			vocab.OnOrderedCollection(it, func(col *vocab.OrderedCollection) error {
				if max < len(col.OrderedItems) {
					col.OrderedItems = col.OrderedItems[0:max]
				}
				return nil
			})
		}
		if CollectionTypes.Contains(it.GetType()) {
			vocab.OnCollection(it, func(col *vocab.Collection) error {
				if max < len(col.Items) {
					col.Items = col.Items[0:max]
				}
				return nil
			})
		}
		return true
	}
}

func After(iri vocab.IRI) Fn {
	isAfter := false
	return func(it vocab.Item) bool {
		if vocab.IsNil(it) {
			return isAfter
		}
		if it.GetLink().Equals(iri, true) {
			isAfter = true
			return false
		}
		return isAfter
	}
}

func Before(iri vocab.IRI) Fn {
	isBefore := true
	return func(it vocab.Item) bool {
		if vocab.IsNil(it) {
			return isBefore
		}
		if it.GetLink().Equals(iri, true) {
			isBefore = false
		}
		return isBefore
	}
}
