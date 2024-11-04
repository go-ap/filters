package filters

import vocab "github.com/go-ap/activitypub"

type authorized vocab.IRI

func (a authorized) Apply(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	i := vocab.IRI(a)
	return Any(
		Actor(SameID(i)),
		Any(Recipients(vocab.PublicNS), Recipients(i)),
	).Apply(it)
}

// Authorized creates a filter that checks the [vocab.IRI] against the recipients list of the item it gets applied on.
// The ActivityStreams Public Namespace IRI gets special treatment, because servers use it to signify that the audience of
// an object is public.
func Authorized(iri vocab.IRI) Check {
	return authorized(iri)
}
