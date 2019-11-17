package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)
const deprecatedFlag = "_deprecated_"
var rawContents = regexp.MustCompile(fmt.Sprintf(`- \[([^\]]+)\]\(([^\)]+)\)(%s)?`, deprecatedFlag))

type contentsData struct {
	fileName string
	wikiPageName string
	status status
}

type status int

const (
	availableOnDropbox status = iota
	availableOnWikiOnly
)

type fileContentDetails struct {
	knownPages map[string]*contentsData
}

func newFileContentDetails() *fileContentDetails {
	return &fileContentDetails{
		knownPages: make(map[string]*contentsData),
	}
}

func (f *fileContentDetails) buildContentsFromRaw(raw string) {
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		if rawContents.MatchString(line) {
			match := rawContents.FindStringSubmatch(line)
			fileName := match[1]
			pageName := match[2]
			//log.Printf("Matched input: %v->%v %v", fileName, pageName, deprecated)
			f.knownPages[pageName] = &contentsData{
				fileName:     fileName,
				wikiPageName: pageName,
				status:       availableOnWikiOnly,
			}
		}
	}
}

func (f *fileContentDetails) buildPageFromContents() string {
	var sortedPages = make([]string, 0)
	for pageName := range f.knownPages {
		sortedPages = append(sortedPages, pageName)
	}
	sort.Strings(sortedPages)
	var output = strings.Builder{}
	output.WriteString("# Dropbox Sync\n\n")
	for _, pageName := range sortedPages {
		details := f.knownPages[pageName]
		if details.status == availableOnWikiOnly {
			output.WriteString(fmt.Sprintf("- [%s](%s) %s\n", details.fileName, details.wikiPageName, deprecatedFlag))
		} else {
			output.WriteString(fmt.Sprintf("- [%s](%s)\n", details.fileName, details.wikiPageName))
		}
	}
	return output.String()
}
