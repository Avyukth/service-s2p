package user

import (
	"errors"
)

var (
	RoleAdmin = Role{"ADMIN"}
	RoleUser  = Role{"USER"}
)

var roles = map[string]Role{RoleAdmin.name: RoleUser, RoleUser.name: RoleUser}

type Role struct {
	name string
}

func ParseRole(value string) (Role, error) {

	role, exists := roles[value]
	if !exists {

		return Role{}, errors.New("invalid role")
	}

	return role, nil
}

func MustParseRole(value string) Role {
	role, err := ParseRole(value)
	if err != nil {
		panic(err)
	}
	return role
}

func (r Role) Name() string {
	return r.name
}

func (r *Role) UnmarshalText(data []byte) error {
	r.name = string(data)
	return nil
}

func (r *Role) MarshalText() ([]byte, error) {
	return []byte(r.name), nil
}

func (r Role) Equal(other Role) bool {
	return r.name == other.name
}
