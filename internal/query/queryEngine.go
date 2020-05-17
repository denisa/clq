package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/util"
)

type Query interface {
	Reset()
	Select(w util.BufWriter, heading changelog.Heading) bool
}

type QueryEngine struct {
	queries []Query
	current int
}

func NewQueryEngine(query string) (*QueryEngine, error) {
	qe := &QueryEngine{}
	if len(query) > 0 {
		queryElements := strings.Split(query, ".")
		if err := qe.newChangelogQuery(queryElements[0], queryElements[1:]); err != nil {
			return nil, err
		}
	}
	return qe, nil
}

func (qe *QueryEngine) Apply(w util.BufWriter, heading changelog.Heading) {
	if qe.current < len(qe.queries) && qe.queries[qe.current].Select(w, heading) {
		qe.current++
		qe.queries[qe.current].Reset()
	}
}

func (qe *QueryEngine) newChangelogQuery(name string, queryElements []string) error {
	if name == "title" {
		if len(queryElements) > 0 {
			return fmt.Errorf("no query elements allowed after title")
		}
		qe.queries = append(qe.queries, &changelogQuery{project: func(h changelog.Changelog) string {
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
	project func(changelog.Changelog) string
}

func (q *changelogQuery) Select(w util.BufWriter, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Changelog)
	if !ok {
		return false
	}
	if q.project == nil {
		return true
	} else {
		w.WriteString(q.project(h))
		return false
	}
}

func (q *changelogQuery) Reset() {}

func (qe *QueryEngine) newReleaseQuery(name string, queryElements []string) error {
	if i, err := strconv.Atoi(name); err != nil {
		return fmt.Errorf("Query release selector %q parsing error: %v", name, err)
	} else {
		var f func(changelog.Release) string
		if len(queryElements) == 0 {
			f = func(r changelog.Release) string {
				type Release struct {
					Version string `json:"version"`
					Date    string `json:"date,omitempty"`
				}
				res, _ := json.Marshal(Release{Version: r.Version(), Date: r.Date()})
				var b strings.Builder
				_, _ = b.Write(res)
				return b.String()
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
					if r.Unreleased() {
						return "unreleased"
					}
					if r.Yanked() {
						return "yanked"
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
		w.WriteString(q.project(h))
	}
	return selected && q.project == nil
}

func (q *releaseQuery) Reset() {
	q.cursor = 0
}
