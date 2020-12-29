package changelog

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/denisa/clq/internal/semver"
)

type ChangeKindDto struct {
	Name      string `json:"name"`
	Increment string `json:"increment"`
}
type ChangeKindsDto []*ChangeKindDto

func (s ChangeKindsDto) Len() int      { return len(s) }
func (s ChangeKindsDto) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByName struct{ ChangeKindsDto }

func (s ByName) Less(i, j int) bool { return s.ChangeKindsDto[i].Name < s.ChangeKindsDto[j].Name }

func (c ChangeKind) MarshalJSON() ([]byte, error) {
	var result ChangeKindsDto
	for k, l := range c.changes {
		result = append(result, &ChangeKindDto{Name: k, Increment: l.String()})
	}
	// enforcing arbitrary order for testing
	sort.Sort(ByName{result})
	return json.Marshal(result)
}

func (u *ChangeKind) UnmarshalJSON(data []byte) error {
	var aux ChangeKindsDto
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	for _, val := range aux {
		inc, err := semver.NewIdentifier(val.Increment)
		if err != nil {
			return fmt.Errorf("Error parsing  %q: %q", val.Name, err)
		}
		u.add(val.Name, inc)
	}
	return nil
}
