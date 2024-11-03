package index

import (
	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
	"github.com/spaolacci/murmur3"
)

type (
	tokenTypes interface{ ~string }

	Indexer interface {
		Add(vocab.LinkOrIRI) error
	}

	bitmaps[T tokenTypes] interface {
		get(key T) (*roaring.Bitmap, bool)
	}

	hashFnType                  func(iri vocab.LinkOrIRI) uint32
	extractFnType[T tokenTypes] func(vocab.LinkOrIRI) []T
)

func murmurHash(it vocab.LinkOrIRI) uint32 {
	if it == nil {
		return 0
	}
	h := murmur3.New32()
	_, _ = h.Write([]byte(it.GetLink()))
	return h.Sum32()
}

var hashFn hashFnType = murmurHash

type index[T tokenTypes] struct {
	tokenMap  map[T]*roaring.Bitmap
	extractFn extractFnType[T]
}

func (i *index[T]) Add(li vocab.LinkOrIRI) error {
	ref := hashFn(li)
	if ref == 0 {
		return nil
	}
	tokens := i.extractFn(li)
	for _, tok := range tokens {
		bmp, ok := i.tokenMap[tok]
		if !ok {
			i.tokenMap[tok] = roaring.New()
		}
		bmp.Add(ref)
	}
	return nil
}

func (i *index[T]) get(key T) (*roaring.Bitmap, bool) {
	b, ok := i.tokenMap[key]
	return b, ok
}

func TokenBitmap[T tokenTypes](extractFn extractFnType[T]) Indexer {
	return &index[T]{
		tokenMap:  make(map[T]*roaring.Bitmap),
		extractFn: extractFn,
	}
}
