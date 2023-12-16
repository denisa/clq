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
	Emoji     string `json:"emoji,omitempty"`
}
type ChangeKindsDto []*ChangeKindDto

func (s ChangeKindsDto) Len() int      { return len(s) }
func (s ChangeKindsDto) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByName struct{ ChangeKindsDto }

func (s ByName) Less(i, j int) bool { return s.ChangeKindsDto[i].Name < s.ChangeKindsDto[j].Name }

func (ck *ChangeKind) MarshalJSON() ([]byte, error) {
	var result ChangeKindsDto
	for k, l := range ck.changes {
		result = append(result, &ChangeKindDto{Name: k, Increment: l.semver.String(), Emoji: l.emoji})
	}
	// enforcing arbitrary order for testing
	sort.Sort(ByName{result})
	return json.Marshal(result)
}

func (ck *ChangeKind) UnmarshalJSON(data []byte) error {
	var dto ChangeKindsDto
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}
	for _, val := range dto {
		inc, err := semver.NewIdentifier(val.Increment)
		if err != nil {
			return fmt.Errorf("Error parsing  %q: %q", val.Name, err)
		}
		if err := ck.add(val.Name, inc, val.Emoji); err != nil {
			return err
		}
	}
	return nil
}
