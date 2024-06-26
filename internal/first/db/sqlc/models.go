// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type FirstUserRole string

const (
	FirstUserRoleAdmin    FirstUserRole = "admin"
	FirstUserRolePower    FirstUserRole = "power"
	FirstUserRoleInternal FirstUserRole = "internal"
)

func (e *FirstUserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FirstUserRole(s)
	case string:
		*e = FirstUserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for FirstUserRole: %T", src)
	}
	return nil
}

type NullFirstUserRole struct {
	FirstUserRole FirstUserRole
	Valid         bool // Valid is true if FirstUserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFirstUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.FirstUserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FirstUserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFirstUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FirstUserRole), nil
}

type FirstEntry struct {
	ID            int32
	EntryCode     string
	EntryCategory string
	EntryName     string
	EntryAmount   int32
	EntryWeight   float64
	EntryNote     string
	IsActive      bool
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

type FirstEntryCategory struct {
	ID        int32
	Category  string
	Note      string
	IsActive  bool
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type FirstUser struct {
	ID           int32
	UserName     string
	UserEmail    string
	UserPassword string
	UserRole     FirstUserRole
	IsActive     bool
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}
