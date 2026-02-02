package index

import (
	"time"

	vocab "github.com/go-ap/activitypub"
	"github.com/jdkato/prose/tokenize"
)

// ExtractType returns the "type" of the [vocab.LinkOrIRI].
// This works on both [vocab.Link] and [vocab.Item] objects.
func ExtractType(li vocab.LinkOrIRI) []string {
	it, ok := li.(vocab.ActivityObject)
	if !ok {
		return nil
	}
	types := make([]string, 0)
	if typ := it.GetType(); typ != nil {
		for _, tt := range typ.AsTypes() {
			types = append(types, string(tt))
		}
	}
	if len(types) == 0 {
		return nil
	}
	return types
}

// ExtractName returns a single token composed of the "name" property of the [vocab.LinkOrIRI].
// This works on both [vocab.Link] and [vocab.Item] objects.
func ExtractName(li vocab.LinkOrIRI) []string {
	var name vocab.NaturalLanguageValues
	switch it := li.(type) {
	case vocab.Link:
		name = it.Name
	case *vocab.Link:
		name = it.Name
	case vocab.Item:
		_ = vocab.OnObject(it, func(ob *vocab.Object) error {
			name = ob.Name
			return nil
		})
	}
	return ExtractNatLangVal(name)
}

// ExtractPreferredUsername returns a single token composed of the "preferredUsername" property of the [vocab.Actor]
func ExtractPreferredUsername(li vocab.LinkOrIRI) []string {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	result := make([]string, 0)
	_ = vocab.OnActor(it, func(act *vocab.Actor) error {
		result = ExtractNatLangVal(act.PreferredUsername)
		return nil
	})
	return result
}

// ExtractSummary returns the tokens in the "summary" property of the [vocab.Item]
func ExtractSummary(li vocab.LinkOrIRI) []string {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	result := make([]string, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		result = tokenizeNatLangVal(ob.Summary)
		return nil
	})
	return result
}

// ExtractContent returns the tokens in the "content" property of the [vocab.Item]
func ExtractContent(li vocab.LinkOrIRI) []string {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	result := make([]string, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		result = tokenizeNatLangVal(ob.Content)
		return nil
	})
	return result
}

// ExtractNatLangVal extracts a single token from the value of the [vocab.NaturalLanguageValues] value.
// This is meant for the properties that contain single words like "preferredUsername" or "name".
func ExtractNatLangVal(nlv vocab.NaturalLanguageValues) []string {
	if nlv == nil {
		return nil
	}

	result := make([]string, 0)
	for _, cc := range nlv {
		result = append(result, cc.String())
	}
	return result
}

var (
	sentTokenizer = tokenize.NewPunktSentenceTokenizer()
	wordTokenizer = tokenize.NewTreebankWordTokenizer()
)

func textToWords(text string) []string {
	words := make([]string, 0)
	for _, s := range sentTokenizer.Tokenize(text) {
		words = append(words, wordTokenizer.Tokenize(s)...)
	}

	return words
}

// tokenizeNatLangVal extracts multiple tokens from the value of the [vocab.NaturalLanguageValues] value.
// This is meant for the properties that can contain long texts like "summary" or "content".
// TODO(marius): these usually are HTML, so we should extract the plain text before.
//
//	See something like https://pkg.go.dev/github.com/huantt/plaintext-extractor
func tokenizeNatLangVal(nlv vocab.NaturalLanguageValues) []string {
	if nlv == nil {
		return nil
	}

	tt := t{}
	result := make([]string, 0)
	for _, cc := range nlv {
		txt := cc.String()
		for _, tok := range textToWords(txt) {
			if tt.IsTagOrPunctuation(tok) {
				continue
			}
			result = append(result, tok)
		}
	}
	return result
}

type t struct {
	st  bool
	end bool
}

func (tt *t) IsTagOrPunctuation(s string) bool {
	if tt.st {
		if s == ">" {
			tt.end = true
			tt.st = false
		}
		return true
	}
	if len(s) <= 3 {
		return true
	}
	return false
}

// ExtractRecipients returns the [vocab.IRI] tokens corresponding to the various addressing properties of
// the received [vocab.Item].
// NOTE(marius): Currently it includes *all* the addressing fields, not removing the "blind" ones (Bto and BCC)
func ExtractRecipients(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}
	r, ok := it.(vocab.HasRecipients)
	if !ok {
		return nil
	}
	recipients := r.Recipients()
	if len(recipients) == 0 {
		return nil
	}

	iris := make([]vocab.IRI, 0, len(recipients))
	for _, rec := range recipients {
		iris = append(iris, rec.GetLink())
	}
	return iriToRefs(iris...)
}

// ExtractAttributedTo returns the [vocab.IRI] tokens corresponding to the "attributedTo" property of
// the received [vocab.Item]
func ExtractAttributedTo(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}
	iris := make([]vocab.IRI, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		iris = append(iris, derefObject(ob.AttributedTo)...)
		return nil
	})
	return iriToRefs(iris...)
}

func iriRefFn(li vocab.LinkOrIRI) []uint64 {
	return []uint64{HashFn(li)}
}

// ExtractInReplyTo returns the [vocab.IRI] tokens corresponding to the "inReplyTo" property of
// the received [vocab.Item]
func ExtractInReplyTo(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}
	iris := make([]vocab.IRI, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		iris = append(iris, derefObject(ob.InReplyTo)...)
		return nil
	})
	return iriToRefs(iris...)
}

// ExtractObject returns the [vocab.IRI] tokens corresponding to the "object" property of
// the received [vocab.Activity]
func ExtractObject(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	iris := make(vocab.IRIs, 0)
	_ = vocab.OnActivity(it, func(act *vocab.Activity) error {
		iris = append(iris, derefObject(act.Object)...)
		return nil
	})
	if len(iris) == 0 {
		return nil
	}
	return iriToRefs(iris...)
}

// ExtractID returns the [vocab.IRI] token corresponding to the "ID" property of
// the received [vocab.Item]
func ExtractID(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	iris := make(vocab.IRIs, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		iris = append(iris, ob.ID)
		return nil
	})
	return iriToRefs(iris...)
}

func iriToRefs(iris ...vocab.IRI) []uint64 {
	refs := make([]uint64, len(iris))
	for i, iri := range iris {
		refs[i] = HashFn(iri)
	}
	return refs
}

// ExtractActor returns the [vocab.IRI] tokens corresponding to the "actor" property of
// the received [vocab.IntransitiveActivity]
func ExtractActor(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	iris := make(vocab.IRIs, 0)
	_ = vocab.OnIntransitiveActivity(it, func(act *vocab.IntransitiveActivity) error {
		iris = append(iris, derefObject(act.Actor)...)
		return nil
	})
	if len(iris) == 0 {
		return nil
	}
	return iriToRefs(iris...)
}

// derefObject aggregates the [vocab.IRI] corresponding to received [vocab.Item]
func derefObject(it vocab.Item) []vocab.IRI {
	if vocab.IsNil(it) {
		return nil
	}
	iris := make(vocab.IRIs, 0)
	if it.IsCollection() {
		_ = vocab.OnCollectionIntf(it, func(c vocab.CollectionInterface) error {
			for _, ob := range c.Collection() {
				iris = append(iris, ob.GetLink())
			}
			return nil
		})
	} else {
		iris = append(iris, it.GetLink())
	}
	return iris
}

// ExtractPublished returns the [vocab.IRI] tokens corresponding to the Published property
// the received [vocab.Object]
func ExtractPublished(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	var pub time.Time
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		pub = ob.Published
		return nil
	})
	if pub.IsZero() {
		return nil
	}
	return []uint64{uint64(pub.Round(time.Hour).UnixMicro())}
}

// ExtractUpdated returns the [vocab.IRI] tokens corresponding to the Updated property
// the received [vocab.Object]
func ExtractUpdated(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	var upd time.Time
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		upd = ob.Updated
		return nil
	})
	if upd.IsZero() {
		return nil
	}
	return []uint64{uint64(upd.Round(time.Hour).UnixMicro())}
}

// ExtractCollectionItems returns the [vocab.IRI] tokens corresponding to the items in the collection
// of the received [vocab.Item]
func ExtractCollectionItems(li vocab.LinkOrIRI) []uint64 {
	it, ok := li.(vocab.Item)
	if !ok || !it.IsCollection() {
		return nil
	}

	var iris vocab.IRIs
	_ = vocab.OnCollectionIntf(it, func(col vocab.CollectionInterface) error {
		iris = col.Collection().IRIs()
		return nil
	})
	if len(iris) == 0 {
		return nil
	}

	refs := make([]uint64, 0, len(iris))
	for _, iri := range iris {
		refs = append(refs, HashFn(iri))
	}
	return refs
}
