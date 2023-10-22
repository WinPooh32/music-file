package musicfile

import (
	"strings"
	"testing"
)

func TestExtractTags(t *testing.T) {
	type args struct {
		filename []byte
	}
	tests := []struct {
		name string
		args args
		want Tags
	}{
		{
			name: "no tags",
			args: args{filename: []byte(strings.ToLower("05-Radioactive Contamination"))},
			want: EmptyTags,
		},
		{
			name: "live at 1",
			args: args{filename: []byte(strings.ToLower("a - Live At b"))},
			want: EmptyTags.Set(Live),
		},
		{
			name: "live at 2",
			args: args{filename: []byte(strings.ToLower("a - Live in b c d"))},
			want: EmptyTags.Set(Live),
		},
		{
			name: "live 1",
			args: args{filename: []byte(strings.ToLower("a b (Live in c d)"))},
			want: EmptyTags.Set(Live),
		},
		{
			name: "live 2",
			args: args{filename: []byte(strings.ToLower("a (LIVE)"))},
			want: EmptyTags.Set(Live),
		},
		{
			name: "live 3",
			args: args{filename: []byte(strings.ToLower("a (Live)"))},
			want: EmptyTags.Set(Live),
		},
		{
			name: "live 4",
			args: args{filename: []byte(strings.ToLower("a [Live]"))},
			want: EmptyTags.Set(Live),
		},
		{
			name: "radio 1",
			args: args{filename: []byte(strings.ToLower("b (radio edit)"))},
			want: EmptyTags.Set(Radio),
		},
		{
			name: "radio live 1",
			args: args{filename: []byte(strings.ToLower("c (radio edit live)"))},
			want: EmptyTags.Set(Radio).Set(Live),
		},
		{
			name: "radio live 2",
			args: args{filename: []byte(strings.ToLower("c (radio) [live]"))},
			want: EmptyTags.Set(Radio).Set(Live),
		},
		{
			name: "cover by 1",
			args: args{filename: []byte(strings.ToLower("a - b c d Cover by e"))},
			want: EmptyTags.Set(Cover),
		},
		{
			name: "cover by 2",
			args: args{filename: []byte(strings.ToLower("a - b c d на русском e"))},
			want: EmptyTags.Set(Cover),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractTags(tt.args.filename); got != tt.want {
				t.Errorf("ExtractTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
