package reserve

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Service is a public interface for implementing our Reserve service
type Service interface {
	// fetch all reserve dates with understood params
	ReserveTable(req ReserveRequest) (*UserReservation, error)
	ReserveTables(req []*ReserveRequest) ([]*UserReservation, error)
	CancelReservation(req CancelReserveRequest) (*UserReservationCanceled, error)
}

// service is a private implementation of our profile service
type service struct {
	ds ReserveStore
}

// NewService is a constructor for our Reserve service implementation
func NewService(ds ReserveStore) Service {
	return &service{
		ds: ds,
	}
}

func (s *service) ReserveTable(req ReserveRequest) (*UserReservation, error) {
	t, err := s.ds.FetchTable(req.RestaurantID, req.TableID)
	if err != nil {
		return nil, err
	}

	// make sure the date is valid when reserving this table
	today := time.Now()
	ti, err := time.Parse(time.RFC3339, t.StartDate)
	// this shouldn't error out but check anyways
	if err != nil {
		return nil, err
	}
	// compare dates
	if ti.Before(today) {
		return nil, errors.New("cannot reserve a table in the past")
	}

	// make sure enough seats are available for this reservation
	available := t.NumSeatsAvailable - t.NumSeatsReserved

	if req.NumSeatsReserved > available {
		return nil, errors.New("not enough seats available")
	}

	// begin a transaction
	tx, err := s.ds.DB().Beginx()
	if err != nil {
		return nil, err
	}

	// will only rollback if no COMMIT occurs
	defer tx.Rollback()

	// create the user reservation
	u, err := s.ds.CreateUserReservation(tx, UserReservation{
		ID:           uuid.New().String(),
		RestaurantID: req.RestaurantID,
		TableID:      req.TableID,
		ProfileID:    req.ProfileID,
		NumSeats:     req.NumSeatsReserved,
		StartDate:    t.StartDate,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	})

	if err != nil {
		return nil, err
	}

	// update seats available on the table
	t.NumSeatsReserved = t.NumSeatsReserved + u.NumSeats
	_, err = s.ds.UpdateTable(tx, *t)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	return u, err
}

func (s *service) ReserveTables(req []*ReserveRequest) ([]*UserReservation, error) {
	var uu []*UserReservation
	for _, re := range req {
		u, err := s.ReserveTable(*re)
		if err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}

	return uu, nil
}

func (s *service) CancelReservation(req CancelReserveRequest) (*UserReservationCanceled, error) {
	// fetch user reservation
	r, err := s.ds.FetchUserReservation(req.UserReservationID)
	if err != nil {
		return nil, err
	}

	// fetch table
	t, err := s.ds.FetchTable(req.RestaurantID, req.TableID)
	if err != nil {
		return nil, err
	}

	// begin a transaction
	tx, err := s.ds.DB().Beginx()
	if err != nil {
		return nil, err
	}

	// will only rollback if no COMMIT occurs
	defer tx.Rollback()

	// delete the reservation
	err = s.ds.DeleteUserReservation(tx, *r)
	if err != nil {
		return nil, err
	}

	// create a canceled reservation
	uc, err := s.ds.CreateUserReservationCanceled(tx, UserReservationCanceled{
		ID:           r.ID,
		RestaurantID: r.RestaurantID,
		TableID:      r.TableID,
		ProfileID:    r.ProfileID,
		NumSeats:     r.NumSeats,
		StartDate:    r.StartDate,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	})

	if err != nil {
		return nil, err
	}

	// update the seats on the table
	t.NumSeatsReserved = t.NumSeatsReserved + r.NumSeats
	_, err = s.ds.UpdateTable(tx, *t)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	return uc, err
}
