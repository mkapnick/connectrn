package profile

// ProfileStore interface
type ProfileStore interface {
	Create(*Profile) (*Profile, error)
	Delete(ID string) error
	DeleteByAccountID(accountID string) error
	Update(*Profile) (*Profile, error)
	Fetch(ID string) (*Profile, error)
	FetchByAccountID(accountID string) (*Profile, error)
	FetchAllByCondition(whereCondition string) ([]*Profile, error)
}
