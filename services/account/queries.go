package account

const (
	EmailID int = iota
	ID
)

type IDQuery struct {
	Type  int
	Value string
}
