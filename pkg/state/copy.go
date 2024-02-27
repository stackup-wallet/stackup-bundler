package state

import (
	"encoding/json"
)

// Copy creates a deep copy of a given OverrideSet.
func Copy(os OverrideSet) (OverrideSet, error) {
	cpy := OverrideSet{}
	for k, v := range os {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		oa := OverrideAccount{}
		err = json.Unmarshal(b, &oa)
		if err != nil {
			return nil, err
		}

		cpy[k] = oa
	}

	return cpy, nil
}
