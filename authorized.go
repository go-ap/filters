package filters

import vocab "github.com/go-ap/activitypub"

type authorized vocab.IRI

func itemIsBlock(it vocab.Item) bool {
	isBlock := false
	_ = vocab.OnActivity(it, func(a *vocab.Activity) error {
		isBlock = vocab.BlockType.Match(a.Type)
		return nil
	})
	return isBlock
}

func (a authorized) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	i := vocab.IRI(a)
	ff := Checks{
		Actor(SameID(i)),
		SameAttributedTo(i),
		Recipients(i),
		// NOTE(marius): check also if the item is public.
		IsPublic(),
	}

	if itemIsBlock(it) {
		// NOTE(marius): if the authorized actor filter matches a Block activity's object, we return false.
		// This translates into that specific actor not having access to the Block activities operated against them.
		ff = append(ff, Object(Not(SameID(i))))
	} else {
		ff = append(ff, Object(SameID(i)))
	}
	return Any(ff...).Match(it)
}

// Authorized creates a filter that checks the [vocab.IRI] against the recipients list of the item it gets applied on.
// The ActivityStreams Public Namespace IRI gets special treatment, because servers use it to signify that the audience of
// an object is public.
// The rules for matching this filter are like follows:
//   - for Objects we check their attributedTo property, and their recipients (to, bto, cc, bcc and audience)
//   - for Activities and Intransitive Activities we also check the actor property.
func Authorized(iri vocab.IRI) Check {
	if vocab.PublicNS.Equal(iri) {
		return public{}
	}
	return authorized(iri)
}

// AuthorizedChecks returns all the Authorized checks in the fns slice.
// It recurses if there are Any or All checks, which is not always what you'd want, so take care.
func AuthorizedChecks(fns ...Check) Checks {
	if len(fns) == 0 {
		return fns
	}
	validCheck := func(c Check) bool {
		_, ok := c.(authorized)
		return ok
	}
	return filterCheckFns(validCheck, fns...)
}

// IsPublic creates a filter that matches if items have the [vocab.PublicNS] collection as a recipient.
// This is the general convention for publicly addressed items.
func IsPublic() Check {
	return public{}
}

type public struct{}

// Match checks all the "it" Item's recipients against the public namespace.
func (p public) Match(it vocab.Item) bool {
	if vocab.IsNil(it) {
		return false
	}
	return accumRecipients(it).Contains(vocab.PublicNS)
}
