package profile

const (
	// ID profile id
	ID int = iota
	// AccountID account id
	AccountID
)

// IDQuery id query
type IDQuery struct {
	Type  int
	Value string
}
