package ast

type RoleSpec struct {
	Tag NodeTag[RoleSpec] `json:"tag"`

	Roletype RoleSpecType
	Rolename *string `json:",omitempty"`
	Location int
}

func (n *RoleSpec) Pos() int {
	return n.Location
}
