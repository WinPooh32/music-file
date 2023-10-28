package musicfile

import (
	"bytes"
	"strings"
)

type Info struct {
	Author        string `json:"author,omitempty"`
	Album         string `json:"album,omitempty"`
	Work          string `json:"work,omitempty"`
	Tags          Tags   `json:"tags,omitempty"`
	FileExtension string `json:"file_extension,omitempty"`
}

func ExtractInfo(filepath []byte) (info Info) {
	// Split the file path.
	path := bytes.Split(filepath, []byte("/"))
	return ExtractPathInfo(path)
}

func ExtractPathInfo(path [][]byte) (info Info) {
	if len(path) == 0 {
		return Info{}
	}

	// Extract basename of the file.
	basename := path[len(path)-1]

	info.Author, info.Album, info.Work, info.Tags, info.FileExtension = processBasename(basename)

	for i := 0; i < len(path)-1; i++ {
		dirname := path[i]
		tags := ExtractDirTags(dirname)
		info.Tags = info.Tags.Append(tags)
	}

	return info
}

func ExtractFilenameTags(filename []byte) (tags Tags) {
	if tagsLiveAtRe.Match(filename) {
		tags = tags.Set(Live)
	}
	if tagsInterviewWithRe.Match(filename) {
		tags = tags.Set(Interview)
	}
	if tagsCoverBy.Match(filename) {
		tags = tags.Set(Cover)
	}

	tags = tags.Append(extractParenthesesTags(filename))

	if tagsOriginalMixRe.Match(filename) {
		tags = tags.Del(Remix)
	}
	if tagsMixBy.Match(filename) {
		tags = tags.Set(Remix)
	}

	return tags
}

func ExtractDirTags(dirname []byte) (tags Tags) {
	tags = ExtractFilenameTags(dirname)
	return tags
}

func processBasename(name []byte) (author, album, work string, tags Tags, fileExtension string) {
	// Exclude file extension.
	if i := bytes.LastIndexByte(name, '.'); i >= 0 {
		fileExtension = string(name[i:])
		name = name[0:i]
	}

	fileExtension = strings.TrimSpace(fileExtension)

	if fileExtension == "" {
		fileExtension = "."
	}

	// Fill info struct.

	subexpNames := infoFilenameRe.SubexpNames()

	for _, match := range infoFilenameRe.FindAllSubmatch(name, -1) {
		for groupIdx, group := range match {
			if groupIdx == 0 || len(group) == 0 {
				continue
			}
			groupName := subexpNames[groupIdx]
			if groupName == "" {
				continue
			}
			switch groupName {
			case groupAuthor:
				author = strings.TrimSpace(string(group))
			case groupWork:
				work = strings.TrimSpace(string(group))
			}
		}
	}

	tags = ExtractFilenameTags(name)

	return author, album, work, tags, fileExtension
}

func extractParenthesesTags(name []byte) (tags Tags) {
	for _, match := range parenthesesRe.FindAll(name, -1) {
		tags = tags.Append(extractTagsByRegexp(match))
	}
	return tags
}

func extractTagsByRegexp(name []byte) (tags Tags) {
	re := tagsRe
	groupNames := re.SubexpNames()

	for _, match := range re.FindAllSubmatch(name, -1) {
		for groupIdx, group := range match {
			if groupIdx == 0 {
				continue
			}

			if group == nil {
				continue
			}

			groupName := groupNames[groupIdx]
			if groupName == "" {
				continue
			}

			tags = tags.SetByName(groupName)
		}
	}

	return tags
}
