package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newReleaseQuery(name string, queryElements []string) error {
	i, err := strconv.Atoi(name)
	if err != nil {
		return fmt.Errorf("Query release selector %q parsing error: %v", name, err)
	}

	queryMe := &releaseQuery{index: i}
	qe.queries = append(qe.queries, queryMe)

	if len(queryElements) == 0 {
		queryMe.enter = func(r changelog.Release) string {
			type Release struct {
				Version string `json:"version"`
				Date    string `json:"date,omitempty"`
			}
			res, _ := json.Marshal(Release{Version: r.Version(), Date: r.Date()})
			return string(res)
		}
		return nil
	}

	if strings.HasPrefix(queryElements[0], "changes[") && strings.HasSuffix(queryElements[0], "]") {
		if err := qe.newChangeQuery(queryElements[0][strings.Index(queryElements[0], "[")+1:len(queryElements[0])-1], queryElements[1:]); err != nil {
			return err
		}
		return nil
	}

	switch queryElements[0] {
	default:
		return fmt.Errorf("Query attribute not recognized %q for a \"release\"", queryElements[0])
	case "date":
		queryMe.enter = func(r changelog.Release) string { return r.Date() }
	case "label":
		queryMe.enter = func(r changelog.Release) string { return r.Label() }
	case "status":
		queryMe.enter = func(r changelog.Release) string {
			if !r.HasBeenReleased() {
				return "unreleased"
			}
			if r.HasBeenYanked() {
				return "yanked"
			}
			if r.IsPrerelease() {
				return "prereleased"
			}
			return "released"
		}
	case "version":
		queryMe.enter = func(r changelog.Release) string { return r.Version() }
	}

	return nil
}

type releaseQuery struct {
	cursor int
	index  int
	enter  func(changelog.Release) string
	exit   func(changelog.Release) string
}

func (q *releaseQuery) Enter(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Release)
	if !ok {
		return false
	}
	selected := q.cursor == q.index
	q.cursor++

	if !selected {
		return false
	}
	if q.enter != nil {
		_, _ = w.WriteString(q.enter(h))
	}
	return true
}

func (q *releaseQuery) Exit(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Release)
	if !ok {
		return false
	}
	if q.exit != nil {
		_, _ = w.WriteString(q.exit(h))
	}
	return true
}
