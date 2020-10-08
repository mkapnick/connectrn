package restaurant

// RestaurantStore interface
type RestaurantStore interface {
	CreateRestaurant(r Restaurant) (*Restaurant, error)
	CreateTable(t Table) (*Table, error)
	FetchRestaurant(ID string) (*Restaurant, error)
	FetchTable(ID string) (*Table, error)
	FetchTableByCondition(whereCondition string) (*Table, error)
	FetchRestaurantByCondition(whereCondition string) ([]*Restaurant, error)
	FetchAllTablesByCondition(whereCondition string) ([]*Table, error)
}
