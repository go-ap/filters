package filters

import (
	"testing"

	vocab "github.com/go-ap/activitypub"
)

func dnl(v string) vocab.LangRefValue {
	return nl("-", v)
}
func nl(ref string, v string) vocab.LangRefValue {
	return vocab.LangRefValue{Ref: vocab.LangRef(ref), Value: vocab.Content(v)}
}

func TestNameIs(t *testing.T) {
	type args struct {
		name  string
		names vocab.NaturalLanguageValues
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{},
			want: false,
		},
		{
			name: "empty item",
			args: args{
				name:  "some name",
				names: nil,
			},
			want: false,
		},
		{
			name: "empty name",
			args: args{
				name:  "",
				names: vocab.NaturalLanguageValues{},
			},
			want: false,
		},
		{
			name: "matching name",
			args: args{
				name:  "name",
				names: vocab.NaturalLanguageValues{dnl("name")},
			},
			want: true,
		},
		{
			name: "matching unicode name",
			args: args{
				name:  "日本語",
				names: vocab.NaturalLanguageValues{dnl("日本語")},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := vocab.Object{Name: tt.args.names}
			if got := NameIs(tt.args.name)(it); tt.want != got {
				t.Errorf("NameIs(%q)(Object.Name=%v) = %v, want %v", tt.args.name, tt.args.names, got, tt.want)
			}
			act := vocab.Actor{PreferredUsername: tt.args.names}
			if got := NameIs(tt.args.name)(act); tt.want != got {
				t.Errorf("NameIs(%q)(Actor.PreferredName=%v) = %v, want %v", tt.args.name, tt.args.names, got, tt.want)
			}
		})
	}
}
