package entity

//go:generate go-enum --marshal --names --values

// Role represents user access level
// ENUM(admin, user)
type Role string
