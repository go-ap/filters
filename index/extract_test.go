package index

import (
	"reflect"
	"testing"
	"time"

	vocab "github.com/go-ap/activitypub"
)

func Test_derefObject(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.Item
		want []vocab.IRI
	}{
		{
			name: "empty",
		},
		{
			name: "item collection",
			arg: vocab.ItemCollection{
				&vocab.Object{ID: "https://example.com"},
				vocab.IRI("https://example.com/1"),
			},
			want: vocab.IRIs{"https://example.com", "https://example.com/1"},
		},
		{
			name: "item",
			arg:  &vocab.Object{ID: "https://example.com/666"},
			want: vocab.IRIs{"https://example.com/666"},
		},
		{
			name: "iri",
			arg:  vocab.IRI("https://example.com/667"),
			want: vocab.IRIs{"https://example.com/667"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := derefObject(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("derefObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func kv(k vocab.LangRef, v vocab.Content) func(values *vocab.NaturalLanguageValues) {
	return func(values *vocab.NaturalLanguageValues) {
		_ = values.Append(k, v)
	}
}
func nlv(fns ...func(values *vocab.NaturalLanguageValues)) vocab.NaturalLanguageValues {
	n := make(vocab.NaturalLanguageValues, 0)
	for _, fn := range fns {
		fn(&n)
	}
	return n
}

func Test_ExtractNatLangVal(t *testing.T) {
	tests := []struct {
		name string
		args vocab.NaturalLanguageValues
		want []string
	}{
		{
			name: "empty",
			args: nil,
			want: nil,
		},
		{
			name: "nil lang ref",
			args: nlv(kv(vocab.NilLangRef, vocab.Content("test"))),
			want: []string{"test"},
		},
		{
			name: "multi word",
			args: nlv(kv(vocab.NilLangRef, vocab.Content("lorem ipsum dolor sic amet"))),
			want: []string{"lorem ipsum dolor sic amet"},
		},
		{
			name: "en-fr",
			args: nlv(
				kv("en", vocab.Content("test")),
				kv("fr", vocab.Content("teste")),
			),
			want: []string{"test", "teste"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractNatLangVal(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractNatLangVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tokenizeNatLangVal(t *testing.T) {
	tests := []struct {
		name string
		args vocab.NaturalLanguageValues
		want []string
	}{
		{
			name: "empty",
			args: nil,
			want: nil,
		},
		{
			name: "nil lang ref",
			args: nlv(kv(vocab.NilLangRef, vocab.Content("test"))),
			want: []string{"test"},
		},
		{
			name: "multi word",
			args: nlv(kv(vocab.NilLangRef, vocab.Content("lorem ipsum dolor sic amet"))),
			want: []string{"lorem", "ipsum", "dolor", "amet"},
		},
		{
			name: "en-fr",
			args: nlv(
				kv("en", vocab.Content("test")),
				kv("fr", vocab.Content("teste")),
			),
			want: []string{"test", "teste"},
		},
		{
			name: "en-fr multi word",
			args: nlv(
				kv("en", vocab.Content("lorem ipsum")),
				kv("fr", vocab.Content("teste de teste")),
			),
			want: []string{"lorem", "ipsum", "teste", "teste"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenizeNatLangVal(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenizeNatLangVal() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_ExtractType(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "*Object type",
			arg:  &vocab.Object{Type: vocab.NoteType},
			want: []string{"Note"},
		},
		{
			name: "Object type",
			arg:  vocab.Object{Type: vocab.NoteType},
			want: []string{"Note"},
		},
		{
			name: "*Link type",
			arg:  &vocab.Link{Type: vocab.MentionType},
			want: []string{"Mention"},
		},
		{
			name: "Link type",
			arg:  vocab.Link{Type: vocab.MentionType},
			want: []string{"Mention"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractType(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractType() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_ExtractName(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "*Object name",
			arg:  &vocab.Object{Name: nlv(kv(vocab.NilLangRef, vocab.Content("John Doe")))},
			want: []string{"John Doe"},
		},
		{
			name: "Object name",
			arg:  vocab.Object{Name: nlv(kv(vocab.NilLangRef, vocab.Content("John Doe")))},
			want: []string{"John Doe"},
		},
		{
			name: "*Link name",
			arg:  &vocab.Link{Name: nlv(kv(vocab.NilLangRef, vocab.Content("The empty page")))},
			want: []string{"The empty page"},
		},
		{
			name: "Link name",
			arg:  vocab.Link{Name: nlv(kv(vocab.NilLangRef, vocab.Content("The empty page")))},
			want: []string{"The empty page"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractName(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractName() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_ExtractPreferredUsername(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "*Actor name",
			arg:  &vocab.Actor{PreferredUsername: nlv(kv(vocab.NilLangRef, vocab.Content("John Doe")))},
			want: []string{"John Doe"},
		},
		{
			name: "Actor name",
			arg:  vocab.Actor{PreferredUsername: nlv(kv(vocab.NilLangRef, vocab.Content("John Doe")))},
			want: []string{"John Doe"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractPreferredUsername(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractPreferredUsername() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_ExtractSummary(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "*Object summary",
			arg:  &vocab.Object{Summary: nlv(kv(vocab.NilLangRef, vocab.Content("Lorem ipsum dolor sic amet")))},
			want: []string{"Lorem", "ipsum", "dolor", "amet"},
		},
		{
			name: "Object summary",
			arg:  vocab.Object{Summary: nlv(kv(vocab.NilLangRef, vocab.Content("Lorem ipsum dolor sic amet")))},
			want: []string{"Lorem", "ipsum", "dolor", "amet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractSummary(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractSummary() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
func Test_ExtractContent(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "*Object content",
			arg:  &vocab.Object{Content: nlv(kv(vocab.NilLangRef, vocab.Content("Lorem ipsum dolor sic amet")))},
			want: []string{"Lorem", "ipsum", "dolor", "amet"},
		},
		{
			name: "Object content",
			arg:  vocab.Object{Content: nlv(kv(vocab.NilLangRef, vocab.Content("Lorem ipsum dolor sic amet")))},
			want: []string{"Lorem", "ipsum", "dolor", "amet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractContent(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractContent() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_ExtractActor(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []uint64
	}{
		{
			name: "empty",
		},
		{
			name: "Activity with nil actor",
			arg:  &vocab.Activity{Actor: nil},
			want: nil,
		},
		{
			name: "IntransitiveActivity with nil actor",
			arg:  &vocab.IntransitiveActivity{Actor: nil},
			want: nil,
		},
		{
			name: "*Activity",
			arg:  &vocab.Activity{Actor: vocab.IRI("https://example.com/~johnDoe")},
			want: []uint64{getRef("https://example.com/~johnDoe")},
		},
		{
			name: "*IntransitiveActivity",
			arg:  &vocab.IntransitiveActivity{Actor: vocab.IRI("https://example.com/~johnDoe")},
			want: []uint64{getRef("https://example.com/~johnDoe")},
		},
		{
			name: "Activity",
			arg:  vocab.Activity{Actor: vocab.IRI("https://example.com/~johnDoe")},
			want: []uint64{getRef("https://example.com/~johnDoe")},
		},
		{
			name: "IntransitiveActivity",
			arg:  vocab.IntransitiveActivity{Actor: vocab.IRI("https://example.com/~johnDoe")},
			want: []uint64{getRef("https://example.com/~johnDoe")},
		},
		{
			name: "*Activity",
			arg: &vocab.Activity{
				Actor: vocab.IRIs{
					"https://example.com/~johnDoe",
					"https://example.com/~alice",
				},
			},
			want: []uint64{getRef("https://example.com/~johnDoe"), getRef("https://example.com/~alice")},
		},
		{
			name: "*IntransitiveActivity",
			arg: &vocab.IntransitiveActivity{
				Actor: vocab.IRIs{
					"https://example.com/~johnDoe",
					"https://example.com/~alice",
				},
			},
			want: []uint64{getRef("https://example.com/~johnDoe"), getRef("https://example.com/~alice")},
		},
		{
			name: "Activity",
			arg: &vocab.Activity{
				Actor: vocab.IRIs{
					"https://example.com/~johnDoe",
					"https://example.com/~alice",
				},
			},
			want: []uint64{getRef("https://example.com/~johnDoe"), getRef("https://example.com/~alice")},
		},
		{
			name: "IntransitiveActivity",
			arg: &vocab.IntransitiveActivity{
				Actor: vocab.IRIs{
					"https://example.com/~johnDoe",
					"https://example.com/~alice",
				},
			},
			want: []uint64{
				getRef("https://example.com/~johnDoe"),
				getRef("https://example.com/~alice"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractActor(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractActor() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_ExtractObject(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []uint64
	}{
		{
			name: "empty",
		},
		{
			name: "Activity with nil object",
			arg:  &vocab.Activity{Object: nil},
			want: nil,
		},
		{
			name: "*Activity",
			arg:  &vocab.Activity{Object: vocab.IRI("https://example.com/~johnDoe")},
			want: []uint64{getRef("https://example.com/~johnDoe")},
		},
		{
			name: "Activity",
			arg:  vocab.Activity{Object: vocab.IRI("https://example.com/~johnDoe")},
			want: []uint64{getRef("https://example.com/~johnDoe")},
		},
		{
			name: "*Activity",
			arg: &vocab.Activity{
				Object: vocab.IRIs{
					"https://example.com/1",
					"https://example.com/2",
				},
			},
			want: []uint64{getRef("https://example.com/1"), getRef("https://example.com/2")},
		},
		{
			name: "Activity",
			arg: &vocab.Activity{
				Object: vocab.IRIs{
					"https://example.com/~johnDoe",
					"https://example.com/~alice",
				},
			},
			want: []uint64{getRef("https://example.com/~johnDoe"), getRef("https://example.com/~alice")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractObject(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractObject() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_ExtractRecipients(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []uint64
	}{
		{
			name: "empty",
		},
		{
			name: "Object with nil recipients",
			arg:  &vocab.Object{},
			want: nil,
		},
		{
			name: "Object",
			arg:  &vocab.Object{To: vocab.ItemCollection{vocab.IRI("https://example.com/~johnDoe")}},
			want: []uint64{getRef("https://example.com/~johnDoe")},
		},
		{
			name: "Object with multiple recipients",
			arg: &vocab.Object{
				To: vocab.ItemCollection{
					vocab.IRI("https://example.com/~johnDoe"),
					vocab.IRI("https://example.com/~alice"),
				},
			},
			want: []uint64{
				getRef("https://example.com/~johnDoe"),
				getRef("https://example.com/~alice"),
			},
		},
		{
			name: "Object with all addressing filled",
			arg: &vocab.Object{
				To: vocab.ItemCollection{
					vocab.IRI("https://example.com/~johnDoe"),
					vocab.IRI("https://example.com/~alice"),
				},
				CC: vocab.ItemCollection{
					vocab.IRI("https://example.com/~bob"),
				},
				Bto: vocab.ItemCollection{
					vocab.IRI("https://example.com"),
				},
				BCC: vocab.ItemCollection{
					vocab.IRI("https://example.com/~pif"),
				},
			},
			want: []uint64{
				getRef("https://example.com/~johnDoe"),
				getRef("https://example.com/~alice"),
				getRef("https://example.com/~bob"),
				getRef("https://example.com"),
				getRef("https://example.com/~pif"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractRecipients(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractRecipients() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExtractInReplyTo(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []uint64
	}{
		{
			name: "empty",
			arg:  nil,
			want: nil,
		},
		{
			name: "empty",
			arg: &vocab.Object{
				InReplyTo: vocab.IRI("https://example.com"),
			},
			want: []uint64{getRef(vocab.IRI("https://example.com"))},
		},
		{
			name: "Object in reply to multiple objects",
			arg: &vocab.Object{
				InReplyTo: vocab.ItemCollection{
					vocab.IRI("https://example.com/one"),
					vocab.IRI("https://example.com/two"),
				},
			},
			want: []uint64{
				getRef("https://example.com/one"),
				getRef("https://example.com/two"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractInReplyTo(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractInReplyTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPublished(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []uint64
	}{
		{
			name: "empty",
			arg:  nil,
			want: nil,
		},
		{
			name: "nil published",
			arg:  &vocab.Object{},
			want: nil,
		},
		{
			name: "non nil published",
			arg: &vocab.Object{
				Published: time.Unix(7213, 0),
			},
			want: []uint64{uint64(time.Unix(7200, 0).UnixMicro())},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractPublished(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractPublished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractUpdated(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []uint64
	}{
		{
			name: "empty",
			arg:  nil,
			want: nil,
		},
		{
			name: "nil updated",
			arg:  &vocab.Object{},
			want: nil,
		},
		{
			name: "non nil updated",
			arg: &vocab.Object{
				Updated: time.Unix(3666, 0),
			},
			want: []uint64{uint64(time.Unix(3600, 0).UnixMicro())},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractUpdated(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractUpdated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractCollectionItems(t *testing.T) {
	tests := []struct {
		name string
		arg  vocab.LinkOrIRI
		want []uint64
	}{
		{
			name: "empty",
		},
		{
			name: "no items in non collection",
			arg:  &vocab.Object{},
			want: nil,
		},
		{
			name: "no items in empty item collection",
			arg:  vocab.ItemCollection{},
			want: nil,
		},
		{
			name: "no items in empty ordered collection",
			arg:  &vocab.OrderedCollection{},
			want: nil,
		},
		{
			name: "one item in an ordered collection",
			arg: &vocab.OrderedCollection{
				Type: vocab.OrderedCollectionType, // NOTE(marius): the type is used in the OnCollectionIntf function
				OrderedItems: vocab.ItemCollection{
					vocab.IRI("https://example.com"),
				},
			},
			want: []uint64{getRef("https://example.com")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractCollectionItems(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractCollectionItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
