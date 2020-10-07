package profile

// ProfileStore interface
type ProfileStore interface {
	Create(*Profile) (*Profile, error)
	Fetch(ID string) (*Profile, error)
	FetchByAccountID(accountID string) (*Profile, error)
}
