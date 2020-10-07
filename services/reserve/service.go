package reserve

import (
	"fmt"
	"strings"
)

// Service is a public interface for implementing our Reserve service
type Service interface {
	// fetch all reserve dates with understood params
	FetchAll(golfCourseIDs []string, date string) ([]*Reserve, error)
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

func (s *service) FetchAll(golfCourseIDs []string, date string) ([]*Reserve, error) {
	var ids []string
	for _, id := range golfCourseIDs {
		ids = append(ids, fmt.Sprintf("'%s'", id))
	}

	// we want this to be a `between` command
	// we also give 1 day leeway room to take into account UTC times
	and := fmt.Sprintf("'%s' >= start_date::date - '1 day'::interval AND '%s' <= end_date::date + '1 day'::interval", date, date)
	where := fmt.Sprintf("golf_course_id IN (%s) AND %s", strings.Join(ids, ", "), and)

	b, err := s.ds.FetchAllByCondition(where)
	if err != nil {
		return nil, ErrReserveNotFound{}
	}

	return b, nil
}
