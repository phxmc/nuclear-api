package domain

const (
	PermDefault = 1 << (iota + 1)
	PermSuper
)

const (
	GroupModeAll int = iota
	GroupModeAny
)

type PermGroup struct {
	Perms     []int
	GroupMode int
}

func HasPerm(perms int, perm int) bool {
	return (perms & perm) == perm
}
