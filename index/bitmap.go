package index

import (
	"bytes"
	"encoding/gob"

	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
	"github.com/spaolacci/murmur3"
)

type (
	Tokenizable interface{ ~string | uint32 }

	Indexable interface {
		Add(vocab.LinkOrIRI) uint32
	}

	bitmaps[T Tokenizable] interface {
		get(key T) *roaring.Bitmap
		not(key T) *roaring.Bitmap
	}

	HashFnType                   func(vocab.LinkOrIRI) uint32
	ExtractFnType[T Tokenizable] func(vocab.LinkOrIRI) []T
)

var HashSeed uint32 = 666

func murmurHash(it vocab.LinkOrIRI) uint32 {
	if it == nil {
		return 0
	}
	h := murmur3.New32WithSeed(HashSeed)
	_, _ = h.Write([]byte(it.GetLink()))
	return h.Sum32()
}

var HashFn HashFnType = murmurHash

type tokenMap[T Tokenizable] struct {
	m         map[T]*roaring.Bitmap
	extractFn ExtractFnType[T]
}

func (i *tokenMap[T]) MarshalBinary() ([]byte, error) {
	buff := bytes.Buffer{}
	err := gob.NewEncoder(&buff).Encode(i.m)
	return buff.Bytes(), err
}

func (i *tokenMap[T]) UnmarshalBinary(data []byte) error {
	if i.m == nil {
		i.m = make(map[T]*roaring.Bitmap)
	}
	return gob.NewDecoder(bytes.NewReader(data)).Decode(&i.m)
}

func (i *tokenMap[T]) Add(li vocab.LinkOrIRI) uint32 {
	ref := HashFn(li)
	if ref == 0 {
		return 0
	}

	tokens := i.extractFn(li)
	if len(tokens) == 0 {
		return ref
	}

	for _, tok := range tokens {
		if _, ok := i.m[tok]; !ok {
			i.m[tok] = roaring.New()
		}
		i.m[tok].Add(ref)
	}
	return ref
}

// get returns the bitmap values corresponding to the key.
func (i *tokenMap[T]) get(key T) *roaring.Bitmap {
	b, ok := i.m[key]
	if !ok {
		return roaring.New()
	}
	return b
}

// not returns the OR'ed bitmap values for all token maps not corresponding
// to the key.
func (i *tokenMap[T]) not(key T) *roaring.Bitmap {
	b := roaring.New()
	for k, v := range i.m {
		if k == key {
			continue
		}
		b.Or(v)
	}
	return b
}

func TokenBitmap[T Tokenizable](extractFn ExtractFnType[T]) Indexable {
	return &tokenMap[T]{
		m:         make(map[T]*roaring.Bitmap),
		extractFn: extractFn,
	}
}
