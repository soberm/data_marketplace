package model

type Role int32

const (
	RoleUnspecified Role = iota
	RoleAdmin
	RoleUser
)
