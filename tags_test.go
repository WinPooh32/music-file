package musicfile

import "testing"

func TestTags_Intersects(t *testing.T) {
	type args struct {
		tags Tags
	}
	tests := []struct {
		name string
		tr   Tags
		args args
		want bool
	}{
		{
			name: "intersected live+remix vs remix",
			tr:   EmptyTags.Set(Live).Set(Remix),
			args: args{
				tags: EmptyTags.Set(Remix),
			},
			want: true,
		},
		{
			name: "not intersected live+remix vs demo",
			tr:   EmptyTags.Set(Live).Set(Remix),
			args: args{
				tags: EmptyTags.Set(Demo),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Intersects(tt.args.tags); got != tt.want {
				t.Errorf("Tags.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
