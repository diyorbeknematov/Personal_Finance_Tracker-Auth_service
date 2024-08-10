package postgres

import (
	"auth-service/config"
	"auth-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	cfg := &config.Config{
		DB_HOST:     "localhost",
		DB_PORT:     5432,
		DB_USER:     "postgres",
		DB_PASSWORD: "your_password",
		DB_NAME:     "auth_service",
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	t.Log("Connected to postgres successfully")
}

func TestRegisterUser(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewAuthenticationRepository(db)
	user := models.RegisterUser{
		Email:     "test_email@test.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "test_password",
	}

	resp, err := repo.RegisterUser(user)
	if err != nil {
        t.Fatal(err)
    }

	assert.Equal(t, "success", resp.Status)
}

func TestLoginUser(t *testing.T) {
	cfg := config.Load()
    db, err := ConnectDB(cfg)
    if err != nil {
        t.Fatal(err)
    }
    repo := NewAuthenticationRepository(db)
    user := models.LoginUserReq{
        Email:    "test_email@test.com",
        Password: "test_password",
    }

    resp, err := repo.LoginUser(user)
    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, resp.Email, user.Email)
}

func TestEmailExists(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewAuthenticationRepository(db)
	
	exists, err := repo.EmailExists("test_email@test.com")
	assert.NoError(t, err)

	assert.Equal(t, exists, true)
}

func TestLogOut(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewAuthenticationRepository(db)

	resp, err := repo.LogOutUser("d70789c8-37e0-4de6-8195-d900abc0afb5")
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
}

func TestResetPassword(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewAuthenticationRepository(db)

	resp, err := repo.ResetPassword("test_email@test.com", "update_password")
	assert.NoError(t, err)

	assert.Equal(t, resp.Status, "success")
}

func TestManageUserRoles(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewAuthenticationRepository(db)

	resp, err := repo.ManageUserRoles("test_email@test.com", "admin")
	assert.NoError(t, err)

	assert.Equal(t, resp.Status, "success")
}

func TestSaveRefreshToken(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewAuthenticationRepository(db)

	resp, err := repo.SaveRefreshToken(models.RefreshToken{
		Email: "test_email@test.com",
		RefreshToken: "test_token",
		ExpiresAt: time.Now().Add(1*24*time.Hour).Format("2006-01-02 15:04:05"),
	})

	assert.NoError(t, err)

	assert.Equal(t, resp.Status, "success")
}

func TestInvalidateRefreshToken(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewAuthenticationRepository(db)

	resp, err := repo.InvalidateRefreshToken("test_email@test.com")
	assert.NoError(t, err)

	assert.Equal(t, resp.Status, "success")
}

func TestIsRefreshTokenValid(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewAuthenticationRepository(db)

	resp, err := repo.IsRefreshTokenValid("test_email@test.com")
	assert.NoError(t, err)

	assert.Equal(t, resp, true)
}