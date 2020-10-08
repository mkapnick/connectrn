package restaurant

// RestaurantStore interface
type RestaurantStore interface {
	CreateRestaurant(r Restaurant) (*Restaurant, error)
	CreateTable(t Table) (*Table, error)
	FetchRestaurant(ID string) (*Restaurant, error)
	FetchAllTables(whereCondition string) ([]*Table, error)
}
