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
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Introduction); ok {
				rc.set(h.Title())
			}
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
	projections
}

func (q *changelogQuery) Enter(heading changelog.Heading) (bool, Project) {

	if _, ok := heading.(changelog.Introduction); !ok {
		return false, nil
	}
	return true, q.enter
}

func (q *changelogQuery) Exit(heading changelog.Heading) (bool, Project) {
	if _, ok := heading.(changelog.Introduction); !ok {
		return false, nil
	}
	return true, q.exit
}
