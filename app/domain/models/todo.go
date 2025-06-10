package models

import (
	"errors"

	"gorm.io/gorm"
)

// Enumの代わり
// NotStarted：未着手
// Done：完了
// iotaで自動連番にしています
type Status int

const (
	invalid Status = iota
	NotStarted
	Done
)

type Todo struct {
	gorm.Model
	Title  string
	Status Status
}

// StrToStatus converts a string to Status enum type.
// Args:
//   - target: string to convert
//   - correspond: map of string to Status values
//
// Returns:
//   - status: converted Status value
//   - err: error if conversion fails
//
// Error cases:
//   - target is empty
//   - correspond map is empty
//   - target doesn't match any key in correspond map
//   - matched value is outside valid range
func StrToStatus(target string, correspond map[string]Status) (status Status, err error) {
	if target == "" {
		return invalid, errors.New("target string is empty")
	}
	if len(correspond) == 0 {
		return invalid, errors.New("correspond map is empty")
	}

	for key, value := range correspond {
		if key == target && NotStarted <= value && value <= Done {
			return value, nil
		}
	}
	return invalid, errors.New("invalid status value: " + target)
}
