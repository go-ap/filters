package filters

import vocab "github.com/go-ap/activitypub"

type recipients vocab.IRI

func (r recipients) Match(it vocab.Item) bool {
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
