package jwtauthz

import "time"

// JwtClaims holds the claims of the decoded token
type JwtClaims struct {
	ID    int
	Email string
	Name  string
	Roles []interface{}
	Authz *Authz
}

// Role explains the role of a user
type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   int       `json:"updated_by"`
	Tenant      string    `json:"tenant" `
	IsSuper     bool      `json:"is_super"`
}

// Authz explains the roles + permissions a user has within the RBAC system
type Authz struct {
	Roles       []Role   `json:"roles"`
	Permissions []string `json:"permissions"`
}
