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

type attrTo vocab.IRI

func (a attrTo) Apply(it vocab.Item) bool {
	author := attributedTo(it)
	return author.Contains(vocab.IRI(a))
}

// AttributedTo creates a filter that checks the [vocab.IRI] against the attributedTo property of the item it gets applied on.
func AttributedTo(iri vocab.IRI) Check {
	return attrTo(iri)
}

func attributedTo(it vocab.Item) vocab.ItemCollection {
	var rec vocab.ItemCollection
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		if ob.AttributedTo != nil {
			rec = derefObject(ob.AttributedTo)
		}
		return nil
	})
	return rec
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
	var iris vocab.ItemCollection
	if it.IsCollection() {
		_ = vocab.OnCollectionIntf(it, func(c vocab.CollectionInterface) error {
			for _, ob := range c.Collection() {
				iris = append(iris, ob)
			}
			return nil
		})
	} else {
		iris = vocab.ItemCollection{it}
	}
	return iris
}
