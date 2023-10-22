package musicfile

type Info struct {
	Author string `json:"author,omitempty"`
	Album  string `json:"album,omitempty"`
	Work   string `json:"work,omitempty"`
	Tags   Tags   `json:"tags,omitempty"`
}

func ExtractInfo(filepath []byte) Info {
	// TODO
	return Info{}
}

func ExtractTags(filename []byte) (tags Tags) {
	if tagsLiveAtRe.Match(filename) {
		tags = tags.Set(Live)
	}
	if tagsCoverBy.Match(filename) {
		tags = tags.Set(Cover)
	}
	tags = tags.Append(extractParenthesesTags(filename))
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
