package ast

type RoleSpec struct {
	Roletype RoleSpecType
	Rolename *string
	Location int
}

func (n *RoleSpec) Pos() int {
	return n.Location
}
