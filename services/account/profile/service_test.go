package profile_test

import (
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
)

func TestCreate(t *testing.T) {
	var TestCreateTT = []struct {
		name          string
		createProfile func(ctrl *gomock.Controller, t *testing.T)
	}{
		{
			name: "should throw error if profile exists with account",
			createProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				// should always expect fetchByAccountID to be called
				// FetchByAccountID will return an error
				mockProfileStore.EXPECT().FetchByAccountID(prof.AccountID).Return([]interface{}{
					&profile.Profile{}, fmt.Errorf("profile exists with this account"),
				}...)

				// test service
				s := profile.NewService(mockProfileStore)
				p, err := s.Create(prof)

				assert.NotNil(t, err)
				assert.IsType(t, profile.ErrProfileExists{}, err)
				assert.Nil(t, p)
			},
		}, {
			name: "should create profile",
			createProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				// should always expect fetchByAccountID to be called
				// FetchByAccountID will return an error
				mockProfileStore.EXPECT().FetchByAccountID(prof.AccountID).Return([]interface{}{
					nil, nil,
				}...)

				mockProfileStore.EXPECT().Create(gomock.Any())

				// test service
				s := profile.NewService(mockProfileStore)
				p, err := s.Create(prof)

				assert.NotNil(t, prof)
				assert.Nil(t, err)
				assert.Equal(t, prof.AccountID, p.AccountID)

				assert.NotEmpty(t, p.ID)
				assert.NotEmpty(t, p.CreatedAt)
				assert.NotEmpty(t, p.UpdatedAt)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range TestCreateTT {
		tt.createProfile(ctrl, t)
	}
}

func TestUpdate(t *testing.T) {
	var TestUpdateTT = []struct {
		name          string
		updateProfile func(ctrl *gomock.Controller, t *testing.T)
	}{
		{
			name: "should throw error if profile not found",
			updateProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				mockProfileStore.EXPECT().Update(&prof).Return([]interface{}{
					nil, profile.ErrProfileNotFound{},
				}...)

				// test service
				s := profile.NewService(mockProfileStore)
				p, err := s.Update(prof)

				assert.NotNil(t, err)
				assert.IsType(t, profile.ErrProfileNotFound{}, err)
				assert.Nil(t, p)
			},
		}, {
			name: "should update profile",
			updateProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				mockProfileStore.EXPECT().Update(&prof)

				// test service
				s := profile.NewService(mockProfileStore)
				p, err := s.Update(prof)

				// handle when create passes
				assert.NotNil(t, p)
				assert.Nil(t, err)

				assert.Equal(t, prof.AccountID, p.AccountID)
				assert.Equal(t, prof.ID, p.ID)
				assert.NotEmpty(t, p.UpdatedAt)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range TestUpdateTT {
		tt.updateProfile(ctrl, t)
	}
}

func TestDelete(t *testing.T) {
	var TestDeleteTT = []struct {
		name          string
		deleteProfile func(ctrl *gomock.Controller, t *testing.T)
	}{
		{
			name: "should delete profile",
			deleteProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					ID:        "static-id",
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				query := profile.IDQuery{
					Type:  profile.ID,
					Value: prof.ID,
				}

				mockProfileStore.EXPECT().Delete(query.Value).Return(nil)

				// test service
				s := profile.NewService(mockProfileStore)
				err := s.Delete(query)

				// handle when delete passes
				assert.Nil(t, err)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range TestDeleteTT {
		tt.deleteProfile(ctrl, t)
	}
}

func TestFetch(t *testing.T) {
	var TestFetchTT = []struct {
		name         string
		fetchProfile func(ctrl *gomock.Controller, t *testing.T)
	}{
		{
			name: "should fetch profile by id",
			fetchProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					ID:        "static-id",
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				query := profile.IDQuery{
					Type:  profile.ID,
					Value: "static-id",
				}

				mockProfileStore.EXPECT().Fetch(prof.ID).Return([]interface{}{
					&profile.Profile{
						ID: "static-id",
					}, nil,
				}...)

				// test service
				s := profile.NewService(mockProfileStore)
				p, err := s.Fetch(query)

				assert.Nil(t, err)
				assert.Equal(t, prof.ID, p.ID)
			},
		}, {
			name: "should fetch profile by account id",
			fetchProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				query := profile.IDQuery{
					Type:  profile.AccountID,
					Value: "account-id",
				}

				prof := profile.Profile{
					ID:        "static-id",
					AccountID: "account-id",
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				mockProfileStore.EXPECT().FetchByAccountID(prof.AccountID).Return([]interface{}{
					&profile.Profile{
						ID: "static-id",
					}, nil,
				}...)

				// test service
				s := profile.NewService(mockProfileStore)
				p, err := s.Fetch(query)

				assert.Nil(t, err)
				assert.Equal(t, prof.ID, p.ID)
			},
		}, {
			name: "should throw error if fetch by profile ID not found",
			fetchProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					ID:        "static-id",
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				query := profile.IDQuery{
					Type:  profile.ID,
					Value: "static-id",
				}

				mockProfileStore.EXPECT().Fetch(prof.ID).Return([]interface{}{
					nil, profile.ErrProfileNotFound{},
				}...)

				// test service
				s := profile.NewService(mockProfileStore)
				_, err := s.Fetch(query)

				assert.IsType(t, profile.ErrProfileNotFound{}, err)
			},
		}, {
			name: "should throw error if fetch by profile ID returns an error [even if profile ID found]",
			fetchProfile: func(ctrl *gomock.Controller, t *testing.T) {
				// setup for testing
				mockProfileStore := profile.NewMockProfileStore(ctrl)

				prof := profile.Profile{
					ID:        "static-id",
					AccountID: uuid.New().String(),
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				}

				query := profile.IDQuery{
					Type:  profile.ID,
					Value: "static-id",
				}

				mockProfileStore.EXPECT().Fetch(prof.ID).Return([]interface{}{
					&profile.Profile{}, profile.ErrProfileNotFound{},
				}...)

				// test service
				s := profile.NewService(mockProfileStore)
				_, err := s.Fetch(query)

				assert.NotNil(t, err)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range TestFetchTT {
		tt.fetchProfile(ctrl, t)
	}
}
