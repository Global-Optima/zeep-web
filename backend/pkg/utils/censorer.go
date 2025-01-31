package utils

import (
	"fmt"
	goaway "github.com/TwiN/go-away"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// Cleaned-up single-line regexes
var russianProfanityRegex = regexp.MustCompile(`\b(
((у|[нз]а|(хитро|не)?вз?[ыьъ]|с[ьъ]|(и|ра)[зс]ъ?|(о[тб]|п[оа]д)[ьъ]?|(.\B)+?[оаеи-])-?)?
(е[бп](ц|тс|а|у|и|ы|о)|и[пб][ае][тцд][ьъ]|бл[эяеа]?([дт][ьъ]?)?)
.*?|
((н[иеа]|(ра|и)[зс]|[зд]?[ао](т|дн[оа])?|с(м[еи])?|а[пб]ч|в[ъы]?|пр[еи])-?)?
ху([яйиеёю]|л+и|йня|е?шк|вший|мло|[её]вый|чок|[её]во|[её]ман)
.*?|
(о[тб]?|про|на|вы)?м(
анд([ауеыи](л(и[сзщ])?[ауеиы])?|ой|[ао]в.*?|юк(ов|[ауи])?|е[нт]ь|ища)|
уд([яаиое].+?|е?н([ьюия]|ей))|
[ао]л[ао]ф[ьъ]([яиюе]|[еёо]й)
)|
елд[ауые].*?|
ля[тд]ь|
([нз]а|по)х
)\b`)
var englishProfanityRegex = regexp.MustCompile(`\b([a@][s$][s$]$|[a@][s$][s$]h[o0][l1][e3][s$]?|b[a@][s$][t+][a@]rd|b[e3][a@][s$][t+][i1][a@]?[l1]([i1][t+]y)?|b[e3][s$][t+][i1][a@][l1]([i1][t+]y)?|b[i1][t+]ch[s$]?|b[i1][t+]ch[e3]r[s$]?|b[i1][t+]ch[e3][s$]|b[i1][t+]ch[i1]ng?|b[l1][o0]wj[o0]b[s$]?|c[l1][i1][t+]|^(c|k|ck|q)[o0](c|k|ck|q)[s$]?$|(c|k|ck|q)[o0](c|k|ck|q)[s$]u|(c|k|ck|q)[o0](c|k|ck|q)[s$]u(c|k|ck|q)[e3]d|(c|k|ck|q)[o0](c|k|ck|q)[s$]u(c|k|ck|q)[e3]r|(c|k|ck|q)[o0](c|k|ck|q)[s$]u(c|k|ck|q)[i1]ng|(c|k|ck|q)[o0](c|k|ck|q)[s$]u(c|k|ck|q)[s$]|^cum[s$]?$|cumm??[e3]r|cumm?[i1]ngcock|(c|k|ck|q)um[s$]h[o0][t+]|(c|k|ck|q)un[i1][l1][i1]ngu[s$]|(c|k|ck|q)un[i1][l1][l1][i1]ngu[s$]|(c|k|ck|q)unn[i1][l1][i1]ngu[s$]|(c|k|ck|q)un[t+][s$]?|(c|k|ck|q)un[t+][l1][i1](c|k|ck|q)|(c|k|ck|q)un[t+][l1][i1](c|k|ck|q)[e3]r|(c|k|ck|q)un[t+][l1][i1](c|k|ck|q)[i1]ng|cyb[e3]r(ph|f)u(c|k|ck|q)|d[a@]mn|d[i1]ck|d[i1][l1]d[o0]|d[i1][l1]d[o0][s$]|d[i1]n(c|k|ck|q)|d[i1]n(c|k|ck|q)[s$]|[e3]j[a@]cu[l1]|(ph|f)[a@]g[s$]?|(ph|f)[a@]gg[i1]ng|(ph|f)[a@]gg?[o0][t+][s$]?|(ph|f)[a@]gg[s$]|(ph|f)[e3][l1][l1]?[a@][t+][i1][o0]|(ph|f)u(c|k|ck|q)|(ph|f)u(c|k|ck|q)[s$]?|g[a@]ngb[a@]ng[s$]?|g[a@]ngb[a@]ng[e3]d|g[a@]y|h[o0]m?m[o0]|h[o0]rny|j[a@](c|k|ck|q)-?[o0](ph|f)(ph|f)?|j[e3]rk-?[o0](ph|f)(ph|f)?|j[i1][s$z][s$z]?m?|[ck][o0]ndum[s$]?|mast(e|ur)b(8|ait|ate)|n+[i1]+[gq]+[e3]*r+[s$]*|[o0]rg[a@][s$][i1]m[s$]?|[o0]rg[a@][s$]m[s$]?|p[e3]nn?[i1][s$]|p[i1][s$][s$]|p[i1][s$][s$][o0](ph|f)(ph|f)|p[o0]rn|p[o0]rn[o0][s$]?|p[o0]rn[o0]gr[a@]phy|pr[i1]ck[s$]?|pu[s$][s$][i1][e3][s$]|pu[s$][s$]y[s$]?|[s$][e3]x|[s$]h[i1][t+][s$]?|[s$][l1]u[t+][s$]?|[s$]mu[t+][s$]?|[s$]punk[s$]?|[t+]w[a@][t+][s$]?)\b`)

var kazakhBadWords = []string{
	"жаман_сөз1",
	"жаман_сөз2",
}

func CensorText(text string) (string, error) {
	censoredText, err := censorWithGoAway(text)
	if err != nil {
		return "", err
	}

	censoredText, err = censorWithRegex(censoredText, russianProfanityRegex)
	if err != nil {
		return "", err
	}

	censoredText, err = censorWithRegex(censoredText, englishProfanityRegex)
	if err != nil {
		return "", err
	}

	censoredText, err = censorLanguage(censoredText, kazakhBadWords)
	if err != nil {
		return "", err
	}

	return censoredText, nil
}

func censorWithGoAway(text string) (string, error) {
	censoredText := goaway.Censor(text)
	if text != censoredText {
		return censoredText, errors.New("text contains profanities (detected by go-away)")
	}
	return censoredText, nil
}

func censorLanguage(text string, badWords []string) (string, error) {
	for _, word := range badWords {
		pattern := fmt.Sprintf(`(?i)\b%s\b`, regexp.QuoteMeta(word))
		re := regexp.MustCompile(pattern)
		if re.MatchString(text) {
			return text, fmt.Errorf("text contains profanities (detected by custom list: %s)", word)
		}
		replacement := strings.Repeat("*", len(word))
		text = re.ReplaceAllString(text, replacement)
	}
	return text, nil
}

func censorWithRegex(text string, pattern *regexp.Regexp) (string, error) {
	if pattern.MatchString(text) {
		// Return the text with the bad portion replaced by stars, plus an error
		return pattern.ReplaceAllStringFunc(text, func(match string) string {
			return strings.Repeat("*", len(match))
		}), fmt.Errorf("text contains profanities (detected by regex pattern)")
	}
	return text, nil
}
