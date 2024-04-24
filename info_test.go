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
			name: "one slash",
			args: args{
				filepath: []byte("/"),
			},
			wantInfo: Info{FileExtension: "."},
		},
		{
			name: "complex live 1",
			args: args{
				filepath: []byte("/author - 2011 - acoustic live from radio 538/work.mp3"),
			},
			wantInfo: Info{
				Author:        "",
				Album:         "",
				Work:          "work",
				Tags:          EmptyTags.Set(Live),
				FileExtension: ".mp3",
			},
		},
		{
			name: "complex live 2",
			args: args{
				filepath: []byte("/author - live a b c d/work.mp3"),
			},
			wantInfo: Info{
				Author:        "",
				Album:         "",
				Work:          "work",
				Tags:          EmptyTags.Set(Live),
				FileExtension: ".mp3",
			},
		},
		{
			name: "not live 2",
			args: args{
				filepath: []byte("/i like to live/work.mp3"),
			},
			wantInfo: Info{
				Author:        "",
				Album:         "",
				Work:          "work",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "not live 2",
			args: args{
				filepath: []byte("50-60-70-80-90/author - work live work.mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work live work",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "parentheses",
			args: args{
				filepath: []byte("a (e;;moll).mp3"),
			},
			wantInfo: Info{
				Author:        "",
				Album:         "",
				Work:          "a",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "parentheses 2",
			args: args{
				filepath: []byte("03-author-work (abcd-efg. abcd.mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "parentheses 3",
			args: args{
				filepath: []byte("03-author-work (ab (c) d).mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "parentheses 4",
			args: args{
				filepath: []byte("03-author-work (.mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "-",
			args: args{
				filepath: []byte("/a/01-author - work(mix).mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work",
				Tags:          EmptyTags.Set(Remix),
				FileExtension: ".mp3",
			},
		},
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
				Work:          "work name",
				Tags:          EmptyTags,
				FileExtension: ".mp3",
			},
		},
		{
			name: "live at",
			args: args{
				filepath: []byte("zxc live at abvcd/03 - author - work name (original mix).mp3"),
			},
			wantInfo: Info{
				Author:        "author",
				Album:         "",
				Work:          "work name",
				Tags:          EmptyTags.Set(Live),
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
		{
			name: "cover by 3",
			args: args{filename: []byte(strings.ToLower("a (m parody)"))},
			want: EmptyTags.Set(Cover),
		},
		{
			name: "not cover",
			args: args{filename: []byte(strings.ToLower("parody name"))},
			want: EmptyTags,
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

func TestExtractDirTags(t *testing.T) {
	type args struct {
		dirname []byte
	}
	tests := []struct {
		name     string
		args     args
		wantTags Tags
	}{
		{
			name: "cover",
			args: args{
				dirname: []byte("8-й альбом - пародии, посвящённые группе"),
			},
			wantTags: EmptyTags.Set(Cover),
		},
		{
			name: "cover 2",
			args: args{
				dirname: []byte("/дискография/27-й альбом - пародии, посвящённые группе"),
			},
			wantTags: EmptyTags.Set(Cover),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTags := ExtractDirTags(tt.args.dirname); gotTags != tt.wantTags {
				t.Errorf("ExtractDirTags() = %v, want %v", gotTags, tt.wantTags)
			}
		})
	}
}
