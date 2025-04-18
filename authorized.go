package filters

import vocab "github.com/go-ap/activitypub"

type authorized vocab.IRI

func (a authorized) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	i := vocab.IRI(a)
	return Any(
		Any(Actor(SameID(i)), SameAttributedTo(i)),
		Any(Recipients(vocab.PublicNS), Recipients(i)),
	).Match(it)
}

// Authorized creates a filter that checks the [vocab.IRI] against the recipients list of the item it gets applied on.
// The ActivityStreams Public Namespace IRI gets special treatment, because servers use it to signify that the audience of
// an object is public.
// The rules for matching this filter are like follows:
//   - for Objects we check their attributedTo property, and their recipients (to, bto, cc, bcc and audience)
//   - for Activities and Intransitive Activities we also check the actor property.
func Authorized(iri vocab.IRI) Check {
	return authorized(iri)
}

// AuthorizedChecks returns all the Authorized checks in the fns slice.
// It recurses if there are Any or All checks, which is not always what you'd want, so take care.
func AuthorizedChecks(fns ...Check) Checks {
	validCheck := func(c Check) bool {
		_, ok := c.(authorized)
		return ok
	}
	return filterCheckFns(validCheck, fns...)
}
