package filters

import (
	"reflect"
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func Test_Authorized_Match(t *testing.T) {
	tests := []struct {
		name string
		a    vocab.IRI
		it   vocab.Item
		want bool
	}{
		{
			name: "empty is not authorized",
		},
		{
			name: "example.com is not allowed when not matching To",
			a:    "https://example.com",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: false,
		},
		{
			name: "example.com is allowed when matching To",
			a:    "https://example.com/jdoe",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: true,
		},
		{
			name: "example.com is allowed if object has public To",
			a:    "https://example.com",
			it:   &vocab.Object{To: vocab.ItemCollection{vocab.PublicNS}},
			want: true,
		},
		{
			name: "example.com is not allowed when not matching CC",
			a:    "https://example.com",
			it:   &vocab.Object{CC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: false,
		},
		{
			name: "example.com is allowed when matching CC",
			a:    "https://example.com/jdoe",
			it:   &vocab.Object{CC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: true,
		},
		{
			name: "example.com is allowed if object has public CC",
			a:    "https://example.com",
			it:   &vocab.Object{CC: vocab.ItemCollection{vocab.PublicNS}},
			want: true,
		},
		{
			name: "example.com is not allowed when not matching Bto",
			a:    "https://example.com",
			it:   &vocab.Object{Bto: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: false,
		},
		{
			name: "example.com is allowed when matching Bto",
			a:    "https://example.com/jdoe",
			it:   &vocab.Object{Bto: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: true,
		},
		{
			name: "example.com is allowed if object has public Bto",
			a:    "https://example.com",
			it:   &vocab.Object{Bto: vocab.ItemCollection{vocab.PublicNS}},
			want: true,
		},
		{
			name: "example.com is not allowed when not matching BCC",
			a:    "https://example.com",
			it:   &vocab.Object{BCC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: false,
		},
		{
			name: "example.com is allowed when matching BCC",
			a:    "https://example.com/jdoe",
			it:   &vocab.Object{BCC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: true,
		},
		{
			name: "example.com is allowed if object has public BCC",
			a:    "https://example.com",
			it:   &vocab.Object{BCC: vocab.ItemCollection{vocab.PublicNS}},
			want: true,
		},
		{
			name: "example.com is not allowed when not matching audience",
			a:    "https://example.com",
			it:   &vocab.Object{Audience: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: false,
		},
		{
			name: "example.com is allowed when matching Audience",
			a:    "https://example.com/jdoe",
			it:   &vocab.Object{Audience: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			want: true,
		},
		{
			name: "example.com is allowed if object has public Audience",
			a:    "https://example.com",
			it:   &vocab.Object{Audience: vocab.ItemCollection{vocab.PublicNS}},
			want: true,
		},
		{
			name: "PublicNS should *NOT* be authorized for object with empty recipients list",
			a:    vocab.PublicNS,
			it:   &vocab.Object{Type: vocab.TombstoneType},
			want: false,
		},
		{
			name: "attributedTo should be checked for objects",
			a:    "https://example.com/~jdoe",
			it:   &vocab.Object{Type: vocab.TombstoneType, AttributedTo: vocab.IRI("https://example.com/~jdoe")},
			want: true,
		},
		{
			name: "attributedTo with multiple values should be checked for objects",
			a:    "https://example.com/~jdoe",
			it: &vocab.Object{
				Type:         vocab.TombstoneType,
				AttributedTo: vocab.ItemCollection{vocab.IRI("https://example.com/~jdoe")},
			},
			want: true,
		},
		{
			name: "actor with single values should be checked",
			a:    "https://example.com/~jdoe",
			it: &vocab.Activity{
				Type:  vocab.UpdateType,
				Actor: vocab.IRI("https://example.com/~jdoe"),
			},
			want: true,
		},
		{
			name: "actor with multiple values should be checked",
			a:    "https://example.com/~jdoe",
			it: &vocab.Activity{
				Type:  vocab.UndoType,
				Actor: vocab.ItemCollection{vocab.IRI("https://example.com/~jdoe")},
			},
			want: true,
		},
		{
			name: "object with single values should be checked",
			a:    "https://example.com/~jdoe",
			it: &vocab.Activity{
				Type:   vocab.FollowType,
				Object: vocab.IRI("https://example.com/~jdoe"),
			},
			want: true,
		},
		{
			name: "object with multiple values should be checked",
			a:    "https://example.com/~jdoe",
			it: &vocab.Activity{
				Type:   vocab.FollowType,
				Object: vocab.ItemCollection{vocab.IRI("https://example.com/~jdoe")},
			},
			want: true,
		},
		{
			name: "link matches any IRI",
			a:    "https://example.com/~jdoe",
			it:   &vocab.Link{},
			want: true,
		},
		{
			name: "link matches public collection",
			a:    vocab.PublicNS,
			it:   &vocab.Link{},
			want: true,
		},
		{
			name: "object of block",
			a:    "https://example.com/~jdoe",
			it: &vocab.Activity{
				Type:   vocab.BlockType,
				Object: vocab.IRI("https://example.com/~jdoe"),
			},
			want: false,
		},
		{
			name: "object with multiple values of block",
			a:    "https://example.com/~jdoe",
			it: &vocab.Activity{
				Type:   vocab.BlockType,
				Object: vocab.ItemCollection{vocab.IRI("https://example.com/~jdoe"), vocab.IRI("https://example.com/~alice")},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Authorized(tt.a).Match(tt.it); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizedChecks(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want Checks
	}{
		{
			name: "nil",
			args: nil,
			want: nil,
		},
		{
			name: "empty",
			args: Checks{},
			want: Checks{},
		},
		{
			name: "no authorization checks",
			args: Checks{SameID("http://example.com")},
			want: Checks{},
		},
		{
			name: "just authorization check",
			args: Checks{Authorized("http://example.com/~jdoe")},
			want: Checks{Authorized("http://example.com/~jdoe")},
		},
		{
			name: "with authorization check",
			args: Checks{SameID("http://example.com"), Authorized("http://example.com/~jdoe")},
			want: Checks{Authorized("http://example.com/~jdoe")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthorizedChecks(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthorizedChecks() = %v, want %v", got, tt.want)
			}
		})
	}
}
