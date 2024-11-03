package index

import (
	vocab "github.com/go-ap/activitypub"
	"github.com/jdkato/prose/tokenize"
)

// extractType returns the "type" of the [vocab.LinkOrIRI].
// This works on both [vocab.Link] and [vocab.Item] objects.
func extractType(li vocab.LinkOrIRI) []string {
	switch it := li.(type) {
	case vocab.Item:
		return []string{string(it.GetType())}
	case vocab.Link:
		return []string{string(it.GetType())}
	}
	return nil
}

// extractName returns a single token composed of the "name" property of the [vocab.LinkOrIRI].
// This works on both [vocab.Link] and [vocab.Item] objects.
func extractName(li vocab.LinkOrIRI) []string {
	switch it := li.(type) {
	case vocab.Item:
		result := make([]string, 0)
		_ = vocab.OnObject(it, func(ob *vocab.Object) error {
			result = extractNatLangVal(ob.Name)
			return nil
		})
		return result
	case vocab.Link:
		return extractNatLangVal(it.Name)
	}

	return nil
}

// extractPreferredUsername returns a single token composed of the "preferredUsername" property of the [vocab.Actor]
func extractPreferredUsername(li vocab.LinkOrIRI) []string {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	result := make([]string, 0)
	_ = vocab.OnActor(it, func(act *vocab.Actor) error {
		result = extractNatLangVal(act.PreferredUsername)
		return nil
	})
	return result
}

// extractSummary returns the tokens in the "summary" property of the [vocab.Item]
func extractSummary(li vocab.LinkOrIRI) []string {
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

// extractContent returns the tokens in the "content" property of the [vocab.Item]
func extractContent(li vocab.LinkOrIRI) []string {
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

// extractNatLangVal extracts a single token from the value of the [vocab.NaturalLanguageValues] value.
// This is meant for the properties that contain single words like "preferredUsername" or "name".
func extractNatLangVal(nlv vocab.NaturalLanguageValues) []string {
	result := make([]string, 0)
	for _, cc := range nlv {
		result = append(result, cc.Value.String())
	}
	return result
}

// tokenizeNatLangVal extracts multiple tokens from the value of the [vocab.NaturalLanguageValues] value.
// This is meant for the properties that can contain long texts like "summary" or "content".
// TODO(marius): these usually are HTML, so we should extract the plain text before.
//
//	See something like https://pkg.go.dev/github.com/huantt/plaintext-extractor
func tokenizeNatLangVal(nlv vocab.NaturalLanguageValues) []string {
	result := make([]string, 0)
	for _, cc := range nlv {
		lng := "en"
		if cc.Ref == "en" || cc.Ref == "fr" || cc.Ref == "es" {
			lng = string(cc.Ref)
		}
		tokenizer, _ := tokenize.NewPragmaticSegmenter(lng)
		result = append(result, tokenizer.Tokenize(cc.Value.String())...)
	}
	return result
}

// extractRecipients returns the [vocab.IRI] tokens corresponding to the various addressing properties of
// the received [vocab.Item].
// NOTE(marius): Currently it includes *all* the addressing fields, not removing the "blind" ones (Bto and BCC)
func extractRecipients(li vocab.LinkOrIRI) []vocab.IRI {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}
	iris := make([]vocab.IRI, 0)
	if r, ok := it.(vocab.HasRecipients); ok {
		for _, rec := range r.Recipients() {
			iris = append(iris, rec.GetLink())
		}
	}
	return iris
}

// extractObject returns the [vocab.IRI] tokens corresponding to the "attributedTo" property of
// the received [vocab.Item]
func extractAttributedTo(li vocab.LinkOrIRI) []vocab.IRI {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}
	iris := make([]vocab.IRI, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		iris = derefObject(ob.AttributedTo)
		return nil
	})
	return iris
}

// extractObject returns the [vocab.IRI] tokens corresponding to the "object" property of
// the received [vocab.Activity]
func extractObject(li vocab.LinkOrIRI) []vocab.IRI {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	var iris []vocab.IRI
	_ = vocab.OnActivity(it, func(act *vocab.Activity) error {
		iris = derefObject(act.Object)
		return nil
	})
	return iris
}

// extractActor returns the [vocab.IRI] tokens corresponding to the "actor" property of
// the received [vocab.IntransitiveActivity]
func extractActor(li vocab.LinkOrIRI) []vocab.IRI {
	it, ok := li.(vocab.Item)
	if !ok {
		return nil
	}

	var iris []vocab.IRI
	_ = vocab.OnIntransitiveActivity(it, func(act *vocab.IntransitiveActivity) error {
		iris = derefObject(act.Actor)
		return nil
	})
	return iris
}

// derefObject aggregates the [vocab.IRI] corresponding to received [vocab.Item]
func derefObject(it vocab.Item) []vocab.IRI {
	if vocab.IsNil(it) {
		return nil
	}
	var iris []vocab.IRI
	if it.IsCollection() {
		_ = vocab.OnCollectionIntf(it, func(c vocab.CollectionInterface) error {
			for _, ob := range c.Collection() {
				iris = append(iris, ob.GetLink())
			}
			return nil
		})
	} else {
		iris = []vocab.IRI{it.GetLink()}
	}
	return iris
}
