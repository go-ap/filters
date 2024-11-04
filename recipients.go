package filters

import vocab "github.com/go-ap/activitypub"

type recipients vocab.IRI

func (r recipients) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	aud := accumRecipients(it)
	return aud.Contains(vocab.IRI(r))
}

// Recipients creates a filter that checks the [vocab.IRI] against the recipients list of the item it gets applied on.
func Recipients(iri vocab.IRI) Check {
	return recipients(iri)
}

func accumRecipients(it vocab.Item) vocab.ItemCollection {
	if withRec, ok := it.(vocab.HasRecipients); ok {
		return withRec.Recipients()
	}
	return nil
}

func derefObject(it vocab.Item) vocab.ItemCollection {
	if vocab.IsNil(it) {
		return nil
	}
	if vocab.IsIRI(it) {
		return vocab.ItemCollection{it.GetLink()}
	}

	iris := make(vocab.ItemCollection, 0)
	if vocab.IsIRIs(it) {
		_ = vocab.OnIRIs(it, func(col *vocab.IRIs) error {
			for _, r := range *col {
				iris = append(iris, r)
			}
			return nil
		})
	} else if vocab.IsItemCollection(it) {
		_ = vocab.OnItemCollection(it, func(col *vocab.ItemCollection) error {
			for _, a := range *col {
				iris = append(iris, a)
			}
			return nil
		})
	} else {
		_ = vocab.OnObject(it, func(ob *vocab.Object) error {
			iris = append(iris, ob)
			return nil
		})
	}
	return iris
}
