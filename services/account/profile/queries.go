package profile

const (
	// ID profile id
	ID int = iota
	// AccountID account id
	AccountID
	// AdminCheckout via checkout admin
	AdminCheckout
	// FromLoyalty from square loyalty in admin
	FromLoyalty
)

// IDQuery id query
type IDQuery struct {
	Type  int
	Value string
}
