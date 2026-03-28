package filters

import (
	"testing"

	"github.com/go-ap/activitypub"
)

func TestValidObjectCollection(t *testing.T) {
	tests := []struct {
		name string
		typ  activitypub.CollectionPath
		want bool
	}{
		{
			name: "empty",
			typ:  "",
			want: false,
		},
		{
			name: "invalid",
			typ:  "test",
			want: false,
		},
		{
			name: "unknown",
			typ:  activitypub.Unknown,
			want: false,
		},
		{
			name: "blocked",
			typ:  BlockedType,
			want: false,
		},
		{
			name: "ignored",
			typ:  IgnoredType,
			want: false,
		},
		{
			name: "activities",
			typ:  ActivitiesType,
			want: false,
		},
		{
			name: "inbox",
			typ:  activitypub.Inbox,
			want: false,
		},
		{
			name: "outbox",
			typ:  activitypub.Outbox,
			want: false,
		},
		{
			name: "shares",
			typ:  activitypub.Shares,
			want: false,
		},
		{
			name: "likes",
			typ:  activitypub.Likes,
			want: false,
		},
		{
			name: "replies",
			typ:  activitypub.Replies,
			want: false,
		},
		{
			name: "actors",
			typ:  ActorsType,
			want: true,
		},
		{
			name: "objects",
			typ:  ObjectsType,
			want: true,
		},
		{
			name: "following",
			typ:  activitypub.Following,
			want: true,
		},
		{
			name: "followers",
			typ:  activitypub.Followers,
			want: true,
		},
		{
			name: "liked",
			typ:  activitypub.Liked,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidObjectCollection(tt.typ); got != tt.want {
				t.Errorf("ValidObjectCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidActivityCollection(t *testing.T) {
	tests := []struct {
		name string
		typ  activitypub.CollectionPath
		want bool
	}{
		{
			name: "empty",
			typ:  "",
			want: false,
		},
		{
			name: "invalid",
			typ:  "test",
			want: false,
		},
		{
			name: "unknown",
			typ:  activitypub.Unknown,
			want: false,
		},
		{
			name: "blocked",
			typ:  BlockedType,
			want: false,
		},
		{
			name: "ignored",
			typ:  IgnoredType,
			want: false,
		},
		{
			name: "actors",
			typ:  ActorsType,
			want: false,
		},
		{
			name: "objects",
			typ:  ObjectsType,
			want: false,
		},
		{
			name: "following",
			typ:  activitypub.Following,
			want: false,
		},
		{
			name: "followers",
			typ:  activitypub.Followers,
			want: false,
		},
		{
			name: "liked",
			typ:  activitypub.Liked,
			want: false,
		},
		{
			name: "activities",
			typ:  ActivitiesType,
			want: true,
		},
		{
			name: "inbox",
			typ:  activitypub.Inbox,
			want: true,
		},
		{
			name: "outbox",
			typ:  activitypub.Outbox,
			want: true,
		},
		{
			name: "shares",
			typ:  activitypub.Shares,
			want: true,
		},
		{
			name: "likes",
			typ:  activitypub.Likes,
			want: true,
		},
		{
			name: "replies",
			typ:  activitypub.Replies,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidActivityCollection(tt.typ); got != tt.want {
				t.Errorf("ValidActivityCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidCollection(t *testing.T) {
	tests := []struct {
		name string
		typ  activitypub.CollectionPath
		want bool
	}{
		{
			name: "empty",
			typ:  "",
			want: false,
		},
		{
			name: "invalid",
			typ:  "test",
			want: false,
		},
		{
			name: "unknown",
			typ:  activitypub.Unknown,
			want: false,
		},
		{
			name: "blocked",
			typ:  BlockedType,
			want: false,
		},
		{
			name: "ignored",
			typ:  IgnoredType,
			want: false,
		},
		{
			name: "actors",
			typ:  ActorsType,
			want: true,
		},
		{
			name: "objects",
			typ:  ObjectsType,
			want: true,
		},
		{
			name: "following",
			typ:  activitypub.Following,
			want: true,
		},
		{
			name: "followers",
			typ:  activitypub.Followers,
			want: true,
		},
		{
			name: "liked",
			typ:  activitypub.Liked,
			want: true,
		},
		{
			name: "activities",
			typ:  ActivitiesType,
			want: true,
		},
		{
			name: "inbox",
			typ:  activitypub.Inbox,
			want: true,
		},
		{
			name: "outbox",
			typ:  activitypub.Outbox,
			want: true,
		},
		{
			name: "shares",
			typ:  activitypub.Shares,
			want: true,
		},
		{
			name: "likes",
			typ:  activitypub.Likes,
			want: true,
		},
		{
			name: "replies",
			typ:  activitypub.Replies,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidCollection(tt.typ); got != tt.want {
				t.Errorf("ValidCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
