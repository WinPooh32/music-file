package musicfile

import (
	"fmt"
)

const bits = 16

var EmptyTags Tags = 0

type Tags uint16

func (t *Tags) Empty() bool {
	return *t == 0
}

func (t Tags) SetByName(name string) Tags {
	tag, ok := nameToTag[name]
	if !ok {
		panic(fmt.Sprintf("tag with name '%s' not found", name))
	}
	t = t.Set(tag)
	return t
}

func (t Tags) Append(t2 Tags) Tags {
	return t | t2
}

func (t Tags) Set(tag TagBit) Tags {
	t |= 1 << Tags(tag)
	return t
}

func (t Tags) Del(tag TagBit) Tags {
	t &= ^(1 << Tags(tag))
	return t
}

func (t Tags) Has(tag TagBit) bool {
	return t&(1<<Tags(tag)) != 0
}

func (t Tags) Names(dst []string) (n int) {
	if len(dst) < bits {
		panic(fmt.Sprintf("dst len must be not less than %d", bits))
	}

	for i := 0; i < bits; i++ {
		bit := TagBit(i)

		if t.Has(bit) {
			dst[i] = bit.String()
			n++
		}
	}

	return n
}

type TagBit Tags

const (
	Live         TagBit = 0
	Remix        TagBit = 1
	Instrumental TagBit = 2
	Demo         TagBit = 3
	Orchestral   TagBit = 4
	Interview    TagBit = 5
	Interlude    TagBit = 6
	Remaster     TagBit = 7
	Capella      TagBit = 8
	Radio        TagBit = 9
	BackingTrack TagBit = 10
	Fragment     TagBit = 11
	Cover        TagBit = 12
	Rehearsal    TagBit = 13
	Bonus        TagBit = 14
	Draft        TagBit = 15
)

var tagNames = []string{
	"Live",
	"Remix",
	"Instrumental",
	"Demo",
	"Orchestral",
	"Interview",
	"Interlude",
	"Remaster",
	"Capella",
	"Radio",
	"BackingTrack",
	"Fragment",
	"Cover",
	"Rehearsal",
	"Bonus",
	"Draft",
}

var nameToTag = map[string]TagBit{
	"Live":         Live,
	"Remix":        Remix,
	"Instrumental": Instrumental,
	"Demo":         Demo,
	"Orchestral":   Orchestral,
	"Interview":    Interview,
	"Interlude":    Interlude,
	"Remaster":     Remaster,
	"Capella":      Capella,
	"Radio":        Radio,
	"BackingTrack": BackingTrack,
	"Fragment":     Fragment,
	"Cover":        Cover,
	"Rehearsal":    Rehearsal,
	"Bonus":        Bonus,
	"Draft":        Draft,
}

func (tb TagBit) String() string {
	if int(tb) < len(tagNames) {
		return tagNames[int(tb)]
	}
	return "unknown"
}
