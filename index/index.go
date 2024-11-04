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
	ByType Type = iota
	ByName
	ByPreferredUsername
	BySummary
	ByContent
	ByActor
	ByObject
	ByRecipients
	ByAttributedTo
)

// Index represents a full index
// It contains the fast tokenized bitmaps, together with a cross-reference map that provides the corresponding
// [vocab.IRI] list that results after resolving the bitmap searches.
type Index struct {
	w       sync.RWMutex
	Ref     map[uint32]vocab.IRI
	Indexes map[Type]Indexable
}

// Full returns a full index data type.
// The indexable fields can be found in the "ByXX" constants.
func Full() *Index {
	return &Index{
		w:   sync.RWMutex{},
		Ref: make(map[uint32]vocab.IRI),
		Indexes: map[Type]Indexable{
			ByType:              TokenBitmap(extractType),
			ByName:              TokenBitmap(extractName),
			ByPreferredUsername: TokenBitmap(extractPreferredUsername),
			BySummary:           TokenBitmap(extractSummary),
			ByContent:           TokenBitmap(extractContent),
			ByActor:             TokenBitmap(extractActor),
			ByObject:            TokenBitmap(extractObject),
			ByRecipients:        TokenBitmap(extractRecipients),
			ByAttributedTo:      TokenBitmap(extractAttributedTo),
		},
	}
}

// Add adds a [vocab.LinkOrIRI] object to the index.
func (i *Index) Add(li vocab.LinkOrIRI) error {
	i.w.Lock()
	defer i.w.Unlock()

	ref := hashFn(li)
	if ref == 0 {
		return errors.Newf("invalid hash")
	}
	i.Ref[ref] = li.GetLink()

	errs := make([]error, 0)
	for _, bmp := range i.Indexes {
		if err := bmp.Add(li); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errors.Join(errs...)
}

// Find does a fast index search for the received filters.
func (i *Index) Find(filters ...BasicFilter) ([]vocab.IRI, error) {
	if len(filters) == 0 {
		return nil, errors.Errorf("nil filters for index search")
	}
	i.w.RLock()
	defer i.w.RUnlock()

	ands := make([]*roaring.Bitmap, 0, len(filters))
	for _, f := range filters {
		switch f.Type {
		case ByType, ByName, ByPreferredUsername, BySummary, ByContent:
			ands = append(ands, getStringBitmaps(i.Indexes[f.Type], f.Values...))
		case ByActor, ByObject, ByRecipients, ByAttributedTo:
			ands = append(ands, getIRIBitmaps(i.Indexes[f.Type], f.Values...))
		}
	}

	bmp := roaring.FastAnd(ands...)
	result := make([]vocab.IRI, 0, bmp.GetCardinality())
	bmp.Iterate(func(x uint32) bool {
		if iri, ok := i.Ref[x]; ok {
			result = append(result, iri)
		}
		return true
	})
	return result, nil
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

// getIRIBitmaps returns the union of the underlying search bitmaps corresponding to the received values.
func getIRIBitmaps(in Indexable, iris ...string) *roaring.Bitmap {
	bmp, ok := in.(bitmaps[vocab.IRI])
	if !ok {
		return nil
	}
	ors := make([]*roaring.Bitmap, 0, len(iris))
	for _, typ := range iris {
		ti, ok := bmp.get(vocab.IRI(typ))
		if !ok {
			continue
		}
		ors = append(ors, ti)
	}
	return roaring.FastOr(ors...)
}

// getIRIBitmaps returns the union of the underlying search bitmaps to the received values.
func getStringBitmaps(in Indexable, types ...string) *roaring.Bitmap {
	bmp, ok := in.(bitmaps[string])
	if !ok {
		return nil
	}
	ors := make([]*roaring.Bitmap, 0, len(types))
	for _, typ := range types {
		ti, ok := bmp.get(typ)
		if !ok {
			continue
		}
		ors = append(ors, ti)
	}
	return roaring.FastOr(ors...)
}
