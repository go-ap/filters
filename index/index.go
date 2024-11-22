package index

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
)

type Type int8

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
	ByCollection
)

// Index represents a full index
// It contains the fast tokenized bitmaps, together with a cross-reference map that provides the corresponding
// [vocab.IRI] list that results after resolving the bitmap searches.
type Index struct {
	w       sync.RWMutex
	Ref     map[uint32]vocab.IRI
	Indexes map[Type]Indexable
}

var objectIndexTypes = []Type{
	ByID, ByType,
	ByRecipients, ByAttributedTo,
	ByName, BySummary, ByContent,
}

var actorIndexTypes = append(objectIndexTypes, ByPreferredUsername)

var activityIndexTypes = append(objectIndexTypes, ByActor, ByObject)

var collectionIndexTypes = append(objectIndexTypes /*, ByCollection*/)

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
		Ref:     make(map[uint32]vocab.IRI),
		Indexes: make(map[Type]Indexable),
	}
	for _, typ := range types {
		switch typ {
		case ByID:
			i.Indexes[typ] = TokenBitmap(ExtractID)
		case ByType:
			i.Indexes[typ] = TokenBitmap(ExtractType)
		case ByName:
			i.Indexes[typ] = TokenBitmap(ExtractName)
		case ByPreferredUsername:
			i.Indexes[typ] = TokenBitmap(ExtractPreferredUsername)
		case BySummary:
			i.Indexes[typ] = TokenBitmap(ExtractSummary)
		case ByContent:
			i.Indexes[typ] = TokenBitmap(ExtractContent)
		case ByActor:
			i.Indexes[typ] = TokenBitmap(ExtractActor)
		case ByObject:
			i.Indexes[typ] = TokenBitmap(ExtractObject)
		case ByRecipients:
			i.Indexes[typ] = TokenBitmap(ExtractRecipients)
		case ByAttributedTo:
			i.Indexes[typ] = TokenBitmap(ExtractAttributedTo)
		case ByCollection:
			i.Indexes[typ] = CollectionBitmap()
		}
	}
	return &i
}

// Add adds a [vocab.LinkOrIRI] object to the index.
func (i *Index) Add(items ...vocab.LinkOrIRI) error {
	i.w.Lock()
	defer i.w.Unlock()

	errs := make([]error, 0)
	for _, li := range items {
		ref := HashFn(li)
		if ref == 0 {
			errs = append(errs, errors.Newf("invalid hash"))
			continue
		}

		i.Ref[ref] = li.GetLink()

		for _, bmp := range i.Indexes {
			if _, err := bmp.Add(li); err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}

	return errors.Join(errs...)
}

type bareIndex struct {
	Ref     map[uint32]vocab.IRI
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
		Ref:     make(map[uint32]vocab.IRI),
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

// GetBitmaps returns the ORing of the underlying search bitmaps corresponding to the received tokens.
func GetBitmaps[T Tokenizable](in Indexable, tokens ...T) []*roaring.Bitmap {
	bmp, ok := in.(bitmaps[T])
	if !ok {
		return nil
	}
	ors := make([]*roaring.Bitmap, 0, len(tokens))
	for _, typ := range tokens {
		ti := bmp.get(typ)
		ors = append(ors, ti)
	}
	return ors
}
