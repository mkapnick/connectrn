package restaurant

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Service is a public interface for implementing our Restaurant service
type Service interface {
	CreateRestaurant(r RestaurantCreateRequest) (*Restaurant, error)
	CreateTable(t Table) (*Table, error)
	FetchRestaurant(ID string) (*Restaurant, error)
	FetchAllTables(restaurantID string, startDate string) ([]*Table, error)
}

// service is a private implementation of our profile service
type service struct {
	ds RestaurantStore
}

// NewService is a constructor for our Restaurant service implementation
func NewService(ds RestaurantStore) Service {
	return &service{
		ds: ds,
	}
}

func (s *service) CreateRestaurant(r RestaurantCreateRequest) (*Restaurant, error) {
	whereCondition := fmt.Sprintf("name = '%s'", r.Name)
	rr, _ := s.ds.FetchRestaurantByCondition(whereCondition)

	if rr != nil {
		return nil, fmt.Errorf("Restaurant with name %s already exists", r.Name)
	}

	return s.ds.CreateRestaurant(Restaurant{
		ID:        uuid.New().String(),
		Name:      r.Name,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
}

func (s *service) CreateTable(t Table) (*Table, error) {
	t.ID = uuid.New().String()
	t.CreatedAt = uuid.New().String()
	t.UpdatedAt = uuid.New().String()
	return s.ds.CreateTable(t)
}

func (s *service) FetchRestaurant(ID string) (*Restaurant, error) {
	return s.ds.FetchRestaurant(ID)
}

func (s *service) FetchAllTables(restaurantID string, startDate string) ([]*Table, error) {
	_, err := s.ds.FetchRestaurant(restaurantID)
	if err != nil {
		return nil, err
	}

	whereCondition := fmt.Sprintf("restaurant_id = '%s' AND start_date = '%s'", restaurantID, startDate)

	t, err := s.ds.FetchAllTablesByCondition(whereCondition)

	if err != nil {
		return nil, err
	}

	return t, nil
}
