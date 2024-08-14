package postgres

import (
	"auth-service/config"
	"auth-service/generated/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserProfile(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewUserRepository(db)

	resp, err := repo.GetUserProfile("d70789c8-37e0-4de6-8195-d900abc0afb5")
	assert.NoError(t, err)

	assert.Equal(t, resp.Email, "test_email@test.com")
}

func TestUpdateUserProfile(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewUserRepository(db)

	resp, err := repo.UpdateUserProfile(&user.UpdateUserProfileReq{
		Id: "d70789c8-37e0-4de6-8195-d900abc0afb5",
		Email: "test_update_email@test.com",
		FirstName: "UpdateName",
		LastName: "UpdateLastName",
	})

	assert.NoError(t, err)

	assert.Equal(t, resp.Status, "success")
}

func TestChangePassword(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewUserRepository(db)

	resp, err := repo.ChangePassword(&user.ChangePasswordReq{
		Id: "d70789c8-37e0-4de6-8195-d900abc0afb5",
		CurrentPassword: "test_password",
		NewPassword: "update_password",
	})

	assert.NoError(t, err)

	assert.Equal(t, resp.Status, "success")
}

func TestGetUserList(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewUserRepository(db)

	resp, err := repo.GetUsersList(&user.GetUsersListReq{
		Page: 1,
		Limit: 10,
		FirstName: "Test",
	})

	assert.NoError(t, err)

	assert.Equal(t, resp.TotalCount,  int32(0))
}
