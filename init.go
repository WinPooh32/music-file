package musicfile

import (
	"regexp"

	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/dialect/base"
	"github.com/hedhyw/rex/pkg/rex"
)

var (
	tagsRe               *regexp.Regexp
	tagsLiveAtRe         *regexp.Regexp
	tagsFilenameLiveAtRe *regexp.Regexp
	tagsInterviewWithRe  *regexp.Regexp
	tagsCoverBy          *regexp.Regexp
	tagsMixBy            *regexp.Regexp
	tagsOriginalMixRe    *regexp.Regexp
	parenthesesRe        *regexp.Regexp

	infoFilenameRe *regexp.Regexp
)

const (
	groupAuthor = "Author"
	groupWork   = "Work"
)

var groups = map[string][]string{
	Live.String(): {
		"live",
		"(живой )?концерт", "кассета", "радиоэфир",
	},
	Remix.String(): {
		"remix", "mix", "rmx", "alt", "bass", "boost", "disco", "club", "offmix",
		"(metal|rock|piano|guitar|sax|danc) version",
		"ремикс", "микс", "радио", "видео", "клуб", "бас",
	},
	Instrumental.String(): {
		"instrument", "instrumental", "instrumentals", "acoust",
		"инструмент", "инструментал",
	},
	Demo.String(): {
		"demo",
		"демо",
	},
	Orchestral.String(): {
		"orchestra", "orchestral", "orch",
		"оркестр",
	},
	Interview.String(): {
		"interview",
		"интервью",
	},
	Interlude.String(): {
		"interlude",
		"антракт",
	},
	Remaster.String(): {
		"remaster",
		"ремастер",
	},
	Capella.String(): {
		"capella", "acapella",
		"капелла", "акапелла",
	},
	Radio.String(): {
		"radio", "video",
		"радио", "видео", "радиоверсия", "видеоверсия",
	},
	BackingTrack.String(): {
		"backingtrack", "back(ing)? track", "karaok",
		"минус", "караоке",
	},
	Fragment.String(): {
		"fragment", "cut version",
		"фрагмент",
	},
	Cover.String(): {
		"cover",
		"кавер", "ковер", "перепевка", "на русском",
	},
	Rehearsal.String(): {
		"rehearsal",
		"репетиция",
	},
	Bonus.String(): {
		"bonus",
		"бонус",
	},
	Draft.String(): {
		"draft",
		"черновик", "чернов(ое)? сведение",
	},
}

func init() {
	tagsLiveAtRe = rex.New(
		rex.Group.Composite(
			rex.Common.Raw(" -[^-]*live( (from|at|on|in) )?"),
			rex.Common.Raw(" - (живой )?концерт (в|на|у|из) "),
			rex.Common.Raw("na stadione|на стадион(е)?"),
			rex.Common.Raw("концерт(н)?(ные)? запис(и)?"),
			rex.Common.Raw("на рад(ио)? "),
		).NonCaptured(),
	).MustCompile()

	tagsFilenameLiveAtRe = rex.New(
		rex.Group.Composite(
			rex.Common.Raw(" - live (from|at|on|in) "),
			rex.Common.Raw(" - (живой )?концерт (в|на|у|из) "),
			rex.Common.Raw("na stadione|на стадион(е)?"),
			rex.Common.Raw("концерт(н)?(ные)? запис(и)?"),
			rex.Common.Raw("на рад(ио)? "),
		).NonCaptured(),
	).MustCompile()

	tagsInterviewWithRe = rex.New(
		rex.Group.Composite(
			rex.Common.Raw("interview"),
			rex.Common.Raw("intervyu|интерв|интервью"),
		).NonCaptured(),
	).MustCompile()

	tagsCoverBy = rex.New(
		rex.Group.Composite(
			rex.Common.Raw("cover by"),
			rex.Common.Raw("на русском"),
		).NonCaptured(),
	).MustCompile()

	tagsMixBy = rex.New(
		rex.Group.Composite(
			rex.Common.Raw("mix by"),
		).NonCaptured(),
	).MustCompile()

	tagsOriginalMixRe = rex.New(
		rex.Group.Composite(
			rex.Common.Raw("origin(al)? (mix|version)"),
		).NonCaptured(),
	).MustCompile()

	tagsRe = rex.New(tagGroups(groups)).MustCompile()

	parenthesesRe = rex.New(
		rex.Common.Class(
			rex.Chars.Single('('),
			rex.Chars.Single('['),
		),

		rex.Common.NotClass(
			rex.Chars.Single('('),
			rex.Chars.Single(')'),
			rex.Chars.Single('['),
			rex.Chars.Single(']'),
		).Repeat().OneOrMore(),

		rex.Group.Composite(
			rex.Common.Class(
				rex.Chars.Single(')'),
				rex.Chars.Single(']'),
			),
			rex.Chars.End(),
		),
	).MustCompile()

	infoFilenameRe = rex.New(
		rex.Group.NonCaptured(
			rex.Chars.Digits().Repeat().OneOrMore(),
			rex.Chars.Single('.').Repeat().ZeroOrOne(),

			rex.Group.NonCaptured(
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Chars.Single('-'),
			).Repeat().ZeroOrOne(),

			rex.Chars.Whitespace().Repeat().ZeroOrOne(),
		).Repeat().ZeroOrOne(),

		rex.Group.Composite(
			rex.Group.NonCaptured(
				rex.Group.Define(
					rex.Chars.Any().Repeat().OneOrMore(),
				).WithName(groupAuthor),

				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Chars.Single('-'),
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),

				rex.Group.Define(
					rex.Chars.Any().Repeat().OneOrMore(),
				).WithName(groupWork),
			).NonCaptured(),

			rex.Group.Composite(
				rex.Chars.Any().Repeat().OneOrMore(),
			).WithName(groupWork),
		).NonCaptured(),
	).MustCompile()

}

func tagGroups(groups map[string][]string) base.GroupToken {
	var tkns []dialect.Token

	for groupName, tokens := range groups {
		grp := tagRawGroup(tokens...).WithName(groupName)
		tkns = append(tkns, grp)
	}

	return rex.Group.Composite(tkns...).NonCaptured()
}

func tagRawGroup(raws ...string) base.GroupToken {
	var tkns []dialect.Token

	for _, r := range raws {
		tkns = append(tkns, rex.Common.Raw(r))
	}

	return rex.Group.Composite(tkns...)
}
