package musicfile

import (
	"reflect"
	"strings"
	"testing"
)

func TestExtractInfo(t *testing.T) {
	type args struct {
		filepath []byte
	}
	tests := []struct {
		name     string
		args     args
		wantInfo Info
	}{
		{
			name: "dot work",
			args: args{
				filepath: []byte("a b/c/02. d & e/02. work name.mp3"),
			},
			wantInfo: Info{
				Author:        "",
				Album:         "",
				Work:          "work name",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "dot dash work",
			args: args{
				filepath: []byte("a b/c/02. d & e/02. - work name.mp3"),
			},
			wantInfo: Info{
				Author:        "",
				Album:         "",
				Work:          "work name",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "dash work",
			args: args{
				filepath: []byte("a b/c/02. d & e/02. - work name.mp3"),
			},
			wantInfo: Info{
				Author:        "",
				Album:         "",
				Work:          "work name",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "space author work",
			args: args{
				filepath: []byte("a b/c/02. d & e/03 - author - work name.mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work name",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "space author work original mix",
			args: args{
				filepath: []byte("a b/c/02. d & e/03 - author - work name (original mix).mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work name (original mix)",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotInfo := ExtractInfo(tt.args.filepath); !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("ExtractInfo() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

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
			if got := ExtractFilenameTags(tt.args.filename); got != tt.want {
				t.Errorf("ExtractTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkExtractInfo(b *testing.B) {
	var info Info
	filepath := []byte("a b/c/02. d & e/03 - author - work name (radio edit live original mix).mp3")
	for i := 0; i < b.N; i++ {
		info = ExtractInfo(filepath)
	}
	nop(info)
}

func BenchmarkExtractTags(b *testing.B) {
	var tags Tags
	filepath := []byte("c cover by d (radio edit live original mix)")
	for i := 0; i < b.N; i++ {
		tags = ExtractFilenameTags(filepath)
	}
	nop(tags)
}

func nop[T any](a T) {}
