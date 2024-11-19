package index

import (
	"bytes"
	"encoding/gob"

	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
	"github.com/spaolacci/murmur3"
)

type (
	tokener interface{ ~string | uint32 }

	Indexable interface {
		Add(vocab.LinkOrIRI) (uint32, error)
	}

	bitmaps[T tokener] interface {
		get(key T) (*roaring.Bitmap, bool)
	}

	HashFnType               func(vocab.LinkOrIRI) uint32
	ExtractFnType[T tokener] func(vocab.LinkOrIRI) []T
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

type tokenMap[T tokener] struct {
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

func (i *tokenMap[T]) Add(li vocab.LinkOrIRI) (uint32, error) {
	ref := HashFn(li)
	if ref == 0 {
		return 0, nil
	}

	tokens := i.extractFn(li)
	if len(tokens) == 0 {
		return ref, nil
	}

	for _, tok := range tokens {
		if _, ok := i.m[tok]; !ok {
			i.m[tok] = roaring.New()
		}
		i.m[tok].Add(ref)
	}
	return ref, nil
}

func (i *tokenMap[T]) get(key T) (*roaring.Bitmap, bool) {
	b, ok := i.m[key]
	return b, ok
}

func TokenBitmap[T tokener](extractFn ExtractFnType[T]) Indexable {
	return &tokenMap[T]{
		m:         make(map[T]*roaring.Bitmap),
		extractFn: extractFn,
	}
}
