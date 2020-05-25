package query

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/util"
)

func (qe *QueryEngine) newIntroductionQuery(name string, queryElements []string) error {
	if name == "title" {
		if len(queryElements) > 0 {
			return fmt.Errorf("no query elements allowed after title")
		}
		qe.queries = append(qe.queries, &changelogQuery{project: func(h changelog.Introduction) string {
			return h.Name()
		}})
		return nil
	}
	if strings.HasPrefix(name, "releases[") && strings.HasSuffix(name, "]") {
		qe.queries = append(qe.queries, &changelogQuery{})
		if err := qe.newReleaseQuery(name[strings.Index(name, "[")+1:len(name)-1], queryElements); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Query attribute not recognized %q", name)
}

type changelogQuery struct {
	project func(changelog.Introduction) string
}

func (q *changelogQuery) Select(w util.BufWriter, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Introduction)
	if !ok {
		return false
	}
	if q.project == nil {
		return true
	} else {
		_, _ = w.WriteString(q.project(h))
		return false
	}
}
