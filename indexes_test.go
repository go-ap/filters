package filters

import (
	"reflect"
	"testing"

	"github.com/go-ap/filters/index"
)

func TestAggregateFilters(t *testing.T) {
	tests := []struct {
		name string
		args []Check
		want []index.BasicFilter
	}{
		{
			name: "empty",
		},
		{
			name: "Object:example.com",
			args: []Check{
				Object(SameID("https://example.com")),
			},
			want: []index.BasicFilter{
				{
					Values: []string{"https://example.com"},
					Type:   index.ByObject,
				},
			},
		},
		{
			name: "Actor:example.com",
			args: []Check{
				Actor(SameID("https://example.com")),
			},
			want: []index.BasicFilter{
				{
					Values: []string{"https://example.com"},
					Type:   index.ByActor,
				},
			},
		},
		{
			name: "multiple actors",
			args: []Check{
				Actor(SameID("https://example.com/~alice"), SameID("https://example.com/~bob")),
			},
			want: []index.BasicFilter{
				{
					Values: []string{"https://example.com/~alice", "https://example.com/~bob"},
					Type:   index.ByActor,
				},
			},
		},
		{
			name: "Recipients:example.com",
			args: []Check{
				Authorized("https://example.com/~alice"),
				Authorized("https://example.com/~bob"),
			},
			want: []index.BasicFilter{
				{
					Values: []string{"https://example.com/~alice"},
					Type:   index.ByRecipients,
				},
				{
					Values: []string{"https://example.com/~bob"},
					Type:   index.ByRecipients,
				},
			},
		},
		{
			name: "name:JaneDoe",
			args: []Check{
				NameIs("JaneDoe"),
			},
			want: []index.BasicFilter{
				{
					Values: []string{"JaneDoe"},
					Type:   index.ByName,
				},
				{
					Values: []string{"JaneDoe"},
					Type:   index.ByPreferredUsername,
				},
			},
		},
		{
			name: "summary",
			args: []Check{
				SummaryIs("Lorem ipsum dolor sic amet."),
			},
			want: []index.BasicFilter{
				{
					Values: []string{"Lorem ipsum dolor sic amet."},
					Type:   index.BySummary,
				},
			},
		},
		{
			name: "content",
			args: []Check{
				ContentIs("Lorem ipsum dolor sic amet, consectetur adipiscing elit."),
			},
			want: []index.BasicFilter{
				{
					Values: []string{"Lorem ipsum dolor sic amet, consectetur adipiscing elit."},
					Type:   index.ByContent,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AggregateFilters(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateFilters() = %+v, want %+v", got, tt.want)
			}
		})
	}
}