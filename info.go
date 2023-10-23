package musicfile

import "bytes"

type Info struct {
	Author string `json:"author,omitempty"`
	Album  string `json:"album,omitempty"`
	Work   string `json:"work,omitempty"`
	Tags   Tags   `json:"tags,omitempty"`
}

func ExtractInfo(filepath []byte) (info Info) {
	// Split the file path.
	path := bytes.Split(filepath, []byte("/"))

	if len(path) == 0 {
		return Info{}
	}

	// Extract basename of the file.
	basename := path[len(path)-1]

	// Exclude file extension.
	if i := bytes.LastIndexByte(basename, '.'); i >= 0 {
		basename = basename[0:i]
	}

	// Fill info struct.

	subexpNames := infoFilenameRe.SubexpNames()

	for _, match := range infoFilenameRe.FindAllSubmatch(basename, -1) {
		for groupIdx, group := range match {
			if groupIdx == 0 || len(group) == 0 {
				continue
			}
			groupName := subexpNames[groupIdx]
			if groupName == "" {
				continue
			}
			switch groupName {
			case groupWork:
				info.Work = string(group)
			case groupAuthor:
				info.Author = string(group)
			}
		}
	}

	info.Tags = ExtractTags(basename)

	return info
}

func ExtractTags(filename []byte) (tags Tags) {
	if tagsLiveAtRe.Match(filename) {
		tags = tags.Set(Live)
	}
	if tagsCoverBy.Match(filename) {
		tags = tags.Set(Cover)
	}

	tags = tags.Append(extractParenthesesTags(filename))

	if tagsOriginalMixRe.Match(filename) {
		tags = tags.Del(Remix)
	}

	return tags
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
