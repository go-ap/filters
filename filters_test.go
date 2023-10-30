package filters

import (
	vocab "github.com/go-ap/activitypub"
	"net/url"
	"reflect"
	"testing"
)

var mockURL = url.URL{
	Scheme: "https",
	Host:   "example.com",
	Path:   "/",
}

func withValues(u url.URL, q url.Values) url.URL {
	u.RawQuery = q.Encode()
	return u
}

func withAuth(u url.URL, user *url.Userinfo) url.URL {
	u.User = user
	return u
}

func TestFromURL(t *testing.T) {
	tests := []struct {
		name string
		u    url.URL
		item vocab.Item
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "no auth, no query values",
			u:    mockURL,
		},
		{
			name: "with user",
			u:    withAuth(mockURL, url.User("https://example.com/jdoe")),
			item: vocab.ItemCollection{
				&vocab.Activity{Type: "Create", AttributedTo: vocab.IRI("https://example.com/jdoe")},
				&vocab.Activity{Type: "Follow"},
				&vocab.Activity{Type: "Like", Actor: vocab.IRI("https://example.com/jdoe")},
				&vocab.Object{Type: "Note", CC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
				&vocab.Tombstone{Type: "Tombstone", Bto: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
				&vocab.Place{Type: "Place", BCC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
				&vocab.Question{Type: "Question", To: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			},
			want: vocab.ItemCollection{
				&vocab.Activity{Type: "Create", AttributedTo: vocab.IRI("https://example.com/jdoe")},
				&vocab.Activity{Type: "Like", Actor: vocab.IRI("https://example.com/jdoe")},
				&vocab.Object{Type: "Note", CC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
				&vocab.Tombstone{Type: "Tombstone", Bto: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
				&vocab.Place{Type: "Place", BCC: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
				&vocab.Question{Type: "Question", To: vocab.ItemCollection{vocab.IRI("https://example.com/jdoe")}},
			},
		},
		{
			name: "some query values",
			u: withValues(mockURL, url.Values{
				"type": []string{"Create", "Follow"},
			}),
			item: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
				vocab.Activity{Type: "Like"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromURL(tt.u).Run(tt.item)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromURL().Run(%v) = %v, want %v", tt.item, got, tt.want)
			}
		})
	}
}

func TestFromValues(t *testing.T) {
	tests := []struct {
		name string
		v    url.Values
		item vocab.Item
		want vocab.Item
	}{
		{name: "empty"},
		{
			name: "no auth, no query values",
			v:    url.Values{},
		},
		{
			name: "type=Create,Follow",
			v: url.Values{
				"type": []string{"Create", "Follow"},
			},
			item: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
				vocab.Activity{Type: "Like"},
			},
			want: vocab.ItemCollection{
				vocab.Activity{Type: "Create"},
				vocab.Activity{Type: "Follow"},
			},
		},
		{
			name: "ID=https://example.com",
			v: url.Values{
				"id": []string{"https://example.com"},
			},
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
			},
		},
		{
			name: "ID=(null)",
			v: url.Values{
				"id": []string{""},
			},
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Follow{Type: "Follow"},
			},
		},
		{
			name: "ID=(not null)",
			v: url.Values{
				"id": []string{"!"},
			},
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
		},
		{
			name: "ID=(not nil IRI)",
			v: url.Values{
				"id": []string{"!-"},
			},
			item: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
		},
		{
			name: "ID=~example.com",
			v: url.Values{
				"id": []string{"~example.com"},
			},
			item: vocab.ItemCollection{
				vocab.Object{ID: "https://activitypub.rocks", Type: "Page"},
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
				vocab.Follow{Type: "Follow"},
			},
			want: vocab.ItemCollection{
				vocab.Actor{ID: "https://example.com"},
				vocab.Activity{ID: "https://example.com/activity"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromValues(tt.v).Run(tt.item)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromURL().Run(%v) = %v, want %v", tt.item, got, tt.want)
			}
		})
	}
}
