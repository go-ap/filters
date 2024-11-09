package index

import vocab "github.com/go-ap/activitypub"

type opType uint8

const (
	OPEq opType = iota
	OPLike
	OPNot
)

type BasicFilter struct {
	Values []string
	Op     opType
	Type   Type
}

func InCollection(iri vocab.IRI) BasicFilter {
	return BasicFilter{Type: ByCollection, Op: OPEq, Values: []string{iri.String()}}
}
