package query

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/util"
)

func (qe *QueryEngine) newReleaseQuery(name string, queryElements []string) error {
	i, err := strconv.Atoi(name)
	if err != nil {
		return fmt.Errorf("Query release selector %q parsing error: %v", name, err)
	}

	var f func(changelog.Release) string
	if len(queryElements) == 0 {
		f = func(r changelog.Release) string {
			type Release struct {
				Version string `json:"version"`
				Date    string `json:"date,omitempty"`
			}
			res, _ := json.Marshal(Release{Version: r.Version(), Date: r.Date()})
			return string(res)
		}
	} else {
		switch queryElements[0] {
		default:
			return fmt.Errorf("Query attribute not recognized %q for a \"release\"", queryElements[0])
		case "date":
			f = func(r changelog.Release) string { return r.Date() }
		case "label":
			f = func(r changelog.Release) string { return r.Label() }
		case "status":
			f = func(r changelog.Release) string {
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
			f = func(r changelog.Release) string { return r.Version() }
		}
	}

	qe.queries = append(qe.queries, &releaseQuery{index: i, project: f})
	return nil
}

type releaseQuery struct {
	cursor  int
	index   int
	project func(changelog.Release) string
}

func (q *releaseQuery) Select(w util.BufWriter, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Release)
	if !ok {
		return false
	}
	selected := q.cursor == q.index
	q.cursor++

	if selected && q.project != nil {
		_, _ = w.WriteString(q.project(h))
	}
	return selected && q.project == nil
}
