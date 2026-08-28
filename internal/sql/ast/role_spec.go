package ast

type RoleSpec struct {
	Tag NodeTag[RoleSpec] `json:"tag"`

	Roletype RoleSpecType
	Rolename *string
	Location int
}

func (n *RoleSpec) Pos() int {
	return n.Location
}
