package reserve

// ReserveStore interface
type ReserveStore interface {
	FetchAllByCondition(whereCondition string) ([]*Reserve, error)
}
