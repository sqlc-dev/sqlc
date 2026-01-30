package ast

type RoleSpec struct {
	Roletype RoleSpecType
	Rolename *string
	Location int

	BindRolename Node
}

func (n *RoleSpec) Pos() int {
	return n.Location
}
