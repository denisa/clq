package query

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newIntroductionQuery(name string, queryElements []string) error {
	queryMe := &changelogQuery{}
	qe.queries = append(qe.queries, queryMe)
	if name == "title" {
		if len(queryElements) > 0 {
			return fmt.Errorf("no query elements allowed after title")
		}
		queryMe.enter = func(h changelog.Introduction) string {
			return h.Title()
		}
		return nil
	}

	if strings.HasPrefix(name, "releases[") && strings.HasSuffix(name, "]") {
		if err := qe.newReleaseQuery(name[strings.Index(name, "[")+1:len(name)-1], queryElements); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("Query attribute not recognized %q", name)
}

type changelogQuery struct {
	enter func(changelog.Introduction) string
	exit  func(changelog.Introduction) string
}

func (q *changelogQuery) Enter(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Introduction)
	if !ok {
		return false
	}

	if q.enter != nil {
		_, _ = w.WriteString(q.enter(h))
	}
	return true
}

func (q *changelogQuery) Exit(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Introduction)
	if !ok {
		return false
	}

	if q.exit != nil {
		_, _ = w.WriteString(q.exit(h))
	}
	return true
}
