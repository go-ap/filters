package filters

import vocab "github.com/go-ap/activitypub"

type authorized vocab.IRI

func (a authorized) Apply(it vocab.Item) bool {
	aud := fullAudience(it)
	return aud.Contains(vocab.PublicNS) || aud.Contains(vocab.IRI(a))
}

func Authorized(iri vocab.IRI) Check {
	return authorized(iri)
}
