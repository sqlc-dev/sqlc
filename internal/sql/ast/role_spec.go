package ast

type RoleSpec struct {
	Tag NodeTag[RoleSpec] `json:"tag"`

	Roletype RoleSpecType `json:"roletype"`
	Rolename *string      `json:"rolename,omitempty"`
	Location int          `json:"location"`
}

func (n *RoleSpec) Pos() int {
	return n.Location
}
