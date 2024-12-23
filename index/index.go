package index

import (
	"bytes"
	"encoding/gob"
	"sync"

	vocab "github.com/go-ap/activitypub"
)

const (
	ByID Type = iota
	ByType
	ByName
	ByPreferredUsername
	BySummary
	ByContent
	ByActor
	ByObject
	ByRecipients
	ByAttributedTo
	ByInReplyTo
	ByPublished
	ByUpdated
)

// Index represents a full index
// It contains the fast tokenized bitmaps, together with a cross-reference map that provides the corresponding
// [vocab.IRI] list that results after resolving the bitmap searches.
type Index struct {
	w       sync.RWMutex
	Ref     map[uint64]vocab.IRI
	Indexes map[Type]Indexable
}

var objectIndexTypes = []Type{
	ByID, ByType,
	ByRecipients, ByAttributedTo, ByInReplyTo,
	ByPublished, ByUpdated,
	ByName, BySummary, ByContent,
}

var actorIndexTypes = append(objectIndexTypes, ByPreferredUsername)

var activityIndexTypes = append(objectIndexTypes, ByActor, ByObject)

var allIndexTypes = append(append(objectIndexTypes, actorIndexTypes...), activityIndexTypes...)

// Full returns a full index data type.
// The complete list of types can be found in the "ByXX" constants.
func Full() *Index {
	return Partial(allIndexTypes...)
}

// Partial returns a partial index. It will create tokenized bitmaps only for the types it receives as parameters.
// The types can be found in the "ByXX" constants.
func Partial(types ...Type) *Index {
	i := Index{
		w:       sync.RWMutex{},
		Ref:     make(map[uint64]vocab.IRI),
		Indexes: make(map[Type]Indexable),
	}
	for _, typ := range types {
		switch typ {
		case ByID:
			i.Indexes[typ] = All()
		case ByType:
			i.Indexes[typ] = NewTokenIndex(ExtractType)
		case ByName:
			i.Indexes[typ] = NewTokenIndex(ExtractName)
		case ByPreferredUsername:
			i.Indexes[typ] = NewTokenIndex(ExtractPreferredUsername)
		case BySummary:
			i.Indexes[typ] = NewTokenIndex(ExtractSummary)
		case ByContent:
			i.Indexes[typ] = NewTokenIndex(ExtractContent)
		case ByActor:
			i.Indexes[typ] = NewTokenIndex(ExtractActor)
		case ByObject:
			i.Indexes[typ] = NewTokenIndex(ExtractObject)
		case ByRecipients:
			i.Indexes[typ] = NewTokenIndex(ExtractRecipients)
		case ByAttributedTo:
			i.Indexes[typ] = NewTokenIndex(ExtractAttributedTo)
		case ByInReplyTo:
			i.Indexes[typ] = NewTokenIndex(ExtractInReplyTo)
		case ByPublished:
			i.Indexes[typ] = NewIndex(ExtractPublished, ExtractID)
		case ByUpdated:
			i.Indexes[typ] = NewIndex(ExtractUpdated, ExtractID)
		}
	}
	return &i
}

// Add adds a [vocab.LinkOrIRI] object to the index.
func (i *Index) Add(items ...vocab.LinkOrIRI) {
	i.w.Lock()
	defer i.w.Unlock()

	for _, li := range items {
		ref := HashFn(li)
		if ref == 0 {
			continue
		}

		i.Ref[ref] = li.GetLink()

		for _, bmp := range i.Indexes {
			_ = bmp.Add(li)
		}
	}
}

type bareIndex struct {
	Ref     map[uint64]vocab.IRI
	Indexes map[Type]Indexable
}

func (i *Index) MarshalBinary() ([]byte, error) {
	buff := bytes.Buffer{}
	b := bareIndex{Ref: i.Ref, Indexes: i.Indexes}
	err := gob.NewEncoder(&buff).Encode(b)
	return buff.Bytes(), err
}

func (i *Index) UnmarshalBinary(data []byte) error {
	b := bareIndex{
		Ref:     make(map[uint64]vocab.IRI),
		Indexes: make(map[Type]Indexable),
	}
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&b)
	if err != nil {
		return err
	}
	i.Ref = b.Ref
	i.Indexes = b.Indexes
	return nil
}
