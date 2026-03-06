package recon

import (
	"fmt"
	"strconv"
	"strings"

	"cerber/internal/i18n"
)

type GoogleLinksService struct {
	templates []string
}

func NewGoogleLinksService() *GoogleLinksService {
	return &GoogleLinksService{
		templates: []string{
			"https://www.google.com/search?q=site:{{target}}+ext:doc+|+ext:docx+|+ext:odt+|+ext:rtf+|+ext:sxw+|+ext:psw+|+ext:ppt+|+ext:pptx+|+ext:pps+|+ext:csv",
			"https://www.google.com/search?q=site:{{target}}+intitle:index.of",
			"https://www.google.com/search?q=site:{{target}}+ext:xml+|+ext:conf+|+ext:cnf+|+ext:reg+|+ext:inf+|+ext:rdp+|+ext:cfg+|+ext:txt+|+ext:ora+|+ext:ini+|+ext:env",
			"https://www.google.com/search?q=site:{{target}}+ext:sql+|+ext:dbf+|+ext:mdb",
			"https://www.google.com/search?q=site:{{target}}+ext:log",
			"https://www.google.com/search?q=site:{{target}}+ext:bkf+|+ext:bkp+|+ext:bak+|+ext:old+|+ext:backup",
			"https://www.google.com/search?q=site:{{target}}+inurl:login+|+inurl:signin+|+intitle:Login+|+intitle:%22sign+in%22+|+inurl:auth",
			"https://www.google.com/search?q=site:{{target}}+intext:%22sql+syntax+near%22+|+intext:%22syntax+error+has+occurred%22+|+intext:%22incorrect+syntax+near%22+|+intext:%22unexpected+end+of+SQL+command%22+|+intext:%22Warning:+mysql_connect()%22+|+intext:%22Warning:+mysql_query()%22+|+intext:%22Warning:+pg_connect()%22",
			"https://www.google.com/search?q=site:{{target}}+%22PHP+Parse+error%22+|+%22PHP+Warning%22+|+%22PHP+Error%22",
			"https://www.google.com/search?q=site:{{target}}+ext:php+intitle:phpinfo+%22published+by+the+PHP+Group%22",
			"https://www.google.com/search?q=site:pastebin.com%20|%20site:paste2.org%20|%20site:pastehtml.com%20|%20site:slexy.org%20|%20site:snipplr.com%20|%20site:snipt.net%20|%20site:textsnip.com%20|%20site:bitpaste.app%20|%20site:justpaste.it%20|%20site:heypasteit.com%20|%20site:hastebin.com%20|%20site:dpaste.org%20|%20site:dpaste.com%20|%20site:codepad.org%20|%20site:jsitor.com%20|%20site:codepen.io%20|%20site:jsfiddle.net%20|%20site:dotnetfiddle.net%20|%20site:phpfiddle.org%20|%20site:ide.geeksforgeeks.org%20|%20site:repl.it%20|%20site:ideone.com%20|%20site:paste.debian.net%20|%20site:paste.org%20|%20site:paste.org.ru%20|%20site:codebeautify.org%20%20|%20site:codeshare.io%20|%20site:trello.com%20%22{{target}}%22",
			"https://www.google.com/search?q=site:github.com%20|%20site:gitlab.com%20%22{{target}}%22",
			"https://www.google.com/search?q=site:stackoverflow.com%20%22{{target}}%22+",
			"https://www.google.com/search?q=site:{{target}}+inurl:signup+|+inurl:register+|+intitle:Signup",
			"https://www.google.com/search?q=site:*.{{target}}",
			"https://www.google.com/search?q=site:*.*.{{target}}",
			"https://web.archive.org/web/*/{{target}}/*",
		},
	}
}

func (s *GoogleLinksService) GenerateLinks(target string, modeSpec string) ([]string, error) {
	indices, err := s.parseModes(modeSpec)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(indices))
	for _, idx := range indices {
		template := s.templates[idx-1]
		res = append(res, strings.ReplaceAll(template, "{{target}}", target))
	}
	return res, nil
}

func (s *GoogleLinksService) parseModes(modeSpec string) ([]int, error) {
	spec := strings.TrimSpace(strings.ToLower(modeSpec))
	if spec == "" || spec == "all" {
		indices := make([]int, 0, len(s.templates))
		for i := 1; i <= len(s.templates); i++ {
			indices = append(indices, i)
		}
		return indices, nil
	}

	parts := strings.Split(spec, ",")
	indices := make([]int, 0, len(parts))
	seen := make(map[int]struct{})

	for _, part := range parts {
		part = strings.TrimSpace(part)
		mode, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("err_invalid_google_mode"), part)
		}
		if mode < 1 || mode > len(s.templates) {
			return nil, fmt.Errorf(i18n.T("err_google_mode_out_of_range"), mode, len(s.templates))
		}
		if _, ok := seen[mode]; ok {
			continue
		}
		seen[mode] = struct{}{}
		indices = append(indices, mode)
	}

	return indices, nil
}
