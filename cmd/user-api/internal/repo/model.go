package repo

import (
	"errors"
	"time"
)

// Errors that can be returned from the repo.
var (
	ErrUserNotFound = errors.New("user could not be found")
)

// SortOrder describes how to sort results
type SortOrder struct {
	By         OrderBy
	Descending bool
}

// SQL returns the SQL representation of SortOrder, suitable
// for use in an OrderBy clause.
func (s *SortOrder) SQL() string {
	if s == nil {
		return "create_time"
	}

	var sq string
	switch s.By {
	case OrderByName:
		sq = "name"
	case OrderByCreateTime:
		fallthrough
	default:
		// Default to Create Time
		sq = "create_time"
	}

	if s.Descending {
		sq += " DESC"
	}

	return sq
}

// OrderBy is an enumeration over sort orders
type OrderBy int

// Cover all the supported sort orders
const (
	OrderByName OrderBy = iota
	OrderByCreateTime
)

// User describes the database model for a User.
type User struct {
	ID         string
	Name       string
	CreateTime time.Time
}
