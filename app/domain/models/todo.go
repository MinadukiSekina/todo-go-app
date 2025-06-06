package models

import "gorm.io/gorm"

// Enumの代わり
// NotStarted：未着手
// Done：完了
// iotaで自動連番にしています
type Status int

const (
	NotStarted Status = iota
	Done
)

type Todo struct {
	gorm.Model
	Title  string
	Status Status
}
