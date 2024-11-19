package index

import (
	"bytes"
	"encoding"
	"encoding/gob"

	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
)

// CollectionBitmap uses a slightly different logic than a regular tokenMap.
// For collections, instead of storing the item's extracted tokens to the reference of the object's IRI
// we use the collection IRI as a token, and we store the references to the collection's items in the bitmap.
func CollectionBitmap() Indexable {
	return &colIndex{
		m:         make(map[vocab.IRI]*roaring.Bitmap),
		extractFn: ExtractCollectionItems,
	}
}

type colIndex tokenMap[vocab.IRI]

func (i *colIndex) MarshalBinary() ([]byte, error) {
	buff := bytes.Buffer{}
	err := gob.NewEncoder(&buff).Encode(i.m)
	return buff.Bytes(), err
}

func (i *colIndex) UnmarshalBinary(data []byte) error {
	if i.m == nil {
		i.m = make(map[vocab.IRI]*roaring.Bitmap)
	}
	return gob.NewDecoder(bytes.NewReader(data)).Decode(&i.m)
}

var _ encoding.BinaryMarshaler = new(colIndex)
var _ encoding.BinaryUnmarshaler = new(colIndex)

func (i *colIndex) get(key vocab.IRI) (*roaring.Bitmap, bool) {
	b, ok := i.m[key]
	return b, ok
}

func (i *colIndex) Add(li vocab.LinkOrIRI) (uint32, error) {
	tok := li.GetLink()

	iris := ExtractCollectionItems(li)
	cref := HashFn(tok)
	if len(iris) == 0 {
		return cref, nil
	}

	if _, ok := i.m[tok]; !ok {
		i.m[tok] = roaring.New()
	}
	for _, iri := range iris {
		ref := HashFn(iri)
		if ref == 0 {
			continue
		}
		i.m[tok].Add(ref)
	}
	return cref, nil
}

// ExtractCollectionItems returns the [vocab.IRI] tokens corresponding to the Items property of
// the received [vocab.Item] collection.
func ExtractCollectionItems(li vocab.LinkOrIRI) []vocab.IRI {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	var iris []vocab.IRI
	_ = vocab.OnCollectionIntf(it, func(c vocab.CollectionInterface) error {
		iris = c.Collection().IRIs()
		return nil
	})
	return iris
}
