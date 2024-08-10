package postgres

import (
	"auth-service/models"
	"database/sql"
	"fmt"
)

type AuthenticationRepository interface {
	EmailExists(email string) (bool, error)
	RegisterUser(user models.RegisterUser) (*models.Response, error)
	LoginUser(login models.LoginUserReq) (*models.User, error)
	LogOutUser(id string) (*models.Response, error)
	ResetPassword(email string, newPassword string) (*models.Response, error)
	SaveRefreshToken(refreshToken models.RefreshToken) (*models.Response, error)
	InvalidateRefreshToken(email string) (*models.Response, error)
	IsRefreshTokenValid(email string) (bool, error)
	ManageUserRoles(email string, role string) (*models.Response, error)
}

type authenticationRepositoryImpl struct {
	db *sql.DB
}

func NewAuthenticationRepository(db *sql.DB) AuthenticationRepository {
	return &authenticationRepositoryImpl{db: db}
}

func (a *authenticationRepositoryImpl) EmailExists(email string) (bool, error) {
	var exists bool
	err := a.db.QueryRow(`
        SELECT 
			EXISTS (SELECT 1 FROM users WHERE email = $1)
		
    `, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a *authenticationRepositoryImpl) RegisterUser(user models.RegisterUser) (*models.Response, error) {
	_, err := a.db.Exec(`
		INSERT INTO users (
			email, 
			first_name, 
			last_name, 
			password_hash,
			role
		)
			VALUES ($1, $2, $3, $4, $5)
    `, user.Email, user.FirstName, user.LastName, user.Password, "user")

	if err != nil {
		return &models.Response{Status: "error", Message: err.Error()}, err
	}
	return &models.Response{
		Status:  "success",
		Message: "User registered successfully",
	}, nil
}

func (a *authenticationRepositoryImpl) LoginUser(login models.LoginUserReq) (*models.User, error) {
	var user models.User
	err := a.db.QueryRow(`
		SELECT
			id,
			email,
            password_hash,
			role
		FROM
			users
		WHERE 
			deleted_at IS NULL AND email = $1
	`, login.Email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found %s", login.Email)
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *authenticationRepositoryImpl) LogOutUser(id string) (*models.Response, error) {
	_, err := a.db.Exec(`
        UPDATE users
        SET deleted_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `, id)

	if err != nil {
		return &models.Response{Status: "error", Message: err.Error()}, err
	}
	return &models.Response{
		Status:  "success",
		Message: "User logged out successfully",
	}, nil
}

func (a *authenticationRepositoryImpl) ResetPassword(email string, newPassword string) (*models.Response, error) {
	_, err := a.db.Exec(`
        UPDATE users
        SET password_hash = $1
        WHERE email = $2
    `, newPassword, email)

	if err != nil {
		return &models.Response{Status: "error", Message: err.Error()}, err
	}
	return &models.Response{
		Status:  "success",
		Message: "Password reset successfully",
	}, nil
}

func (a *authenticationRepositoryImpl) SaveRefreshToken(refreshToken models.RefreshToken) (*models.Response, error) {
	_, err := a.db.Exec(`
		DELETE FROM 
			refresh_tokens
		WHERE user_email = $1
	`, refreshToken.Email)
	if err != nil {
		return &models.Response{
			Status:  "error",
			Message: err.Error(),
		}, err
	}
	_, err = a.db.Exec(`
		INSERT INTO refresh_tokens (
			user_email,
			token,
			expires_at
		)
		    VALUES ($1, $2, $3)
	`, refreshToken.Email, refreshToken.RefreshToken, refreshToken.ExpiresAt)

	if err != nil {
		return &models.Response{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	return &models.Response{
		Status:  "success",
		Message: "Refresh token saved successfully",
	}, nil
}

func (a *authenticationRepositoryImpl) InvalidateRefreshToken(email string) (*models.Response, error) {
	_, err := a.db.Exec(`
		DELETE FROM 
			refresh_tokens
		WHERE
			user_email = $1
	`, email)

	if err != nil {
		return &models.Response{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	return &models.Response{
		Status:  "success",
		Message: "Refresh token invalidated successfully",
	}, nil
}

func (a *authenticationRepositoryImpl) IsRefreshTokenValid(email string) (bool, error) {
	var count int
	err := a.db.QueryRow(`
		SELECT	
			count(*)
		FROM 
			refresh_tokens
		WHERE 
			user_email = $1 AND expires_at > CURRENT_TIMESTAMP
	`, email).Scan(&count)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (a *authenticationRepositoryImpl) ManageUserRoles(email string, role string) (*models.Response, error) {
	_, err := a.db.Exec(`
        UPDATE users
        SET role = $2
        WHERE email = $1
    `, email, role)

	if err != nil {
		return &models.Response{Status: "error", Message: err.Error()}, err
	}
	return &models.Response{
		Status:  "success",
		Message: "User role updated successfully",
	}, nil
}
