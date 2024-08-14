package postgres

import (
	pb "auth-service/generated/user"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	GetUserProfile(id string) (*pb.UserProfile, error)
	UpdateUserProfile(userProfile *pb.UpdateUserProfileReq) (*pb.UpdateUserProfileResp, error)
	GetUsersList(fUser *pb.GetUsersListReq) (*pb.GetUsersListResp, error)
	ChangePassword(change *pb.ChangePasswordReq) (*pb.ChangePasswordResp, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) GetUserProfile(id string) (*pb.UserProfile, error) {
	var userProfile pb.UserProfile
	err := u.db.QueryRow(`
		SELECT 
			id, 
			email, 
			first_name, 
			last_name, 
			role 
		FROM 
			users 
		WHERE id = $1
	`, id).Scan(&userProfile.Id, &userProfile.Email, &userProfile.FirstName, &userProfile.LastName, &userProfile.Role)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (u *userRepositoryImpl) UpdateUserProfile(userProfile *pb.UpdateUserProfileReq) (*pb.UpdateUserProfileResp, error) {
	_, err := u.db.Exec(`
        UPDATE 
            users 
        SET 
            email = $1, 
            first_name = $2, 
            last_name = $3 
        WHERE 
            id = $4
    `, userProfile.Email, userProfile.FirstName, userProfile.LastName, userProfile.Id)

	if err != nil {
		return &pb.UpdateUserProfileResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}
	return &pb.UpdateUserProfileResp{
		Status:  "success",
		Message: "User profile updated successfully",
	}, nil
}

func (u *userRepositoryImpl) ChangePassword(change *pb.ChangePasswordReq) (*pb.ChangePasswordResp, error) {
	_, err := u.db.Exec(`
        UPDATE 
            users 
        SET 
            password_hash = $1 
        WHERE 
            password_hash = $2 AND id = $3
    `, change.NewPassword, change.CurrentPassword, change.Id)

	if err != nil {
		return &pb.ChangePasswordResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}
	return &pb.ChangePasswordResp{
		Status:  "success",
		Message: "Password updated successfully",
	}, nil
}

func (u *userRepositoryImpl) GetUsersList(fUser *pb.GetUsersListReq) (*pb.GetUsersListResp, error) {
	var (
		args   []interface{}
		filter string
	)

	query := `
		SELECT 
			id, 
			email, 
			first_name, 
			last_name, 
			role 
		FROM 
			users
		WHERE
			deleted_at IS NULL `

	if fUser.FirstName != "" {
		filter += fmt.Sprintf(" AND first_name LIKE $%d", len(args)+1)
		args = append(args, fmt.Sprintf("%%%s%%", fUser.FirstName))
	}
	if fUser.LastName != "" {
		filter += fmt.Sprintf(" AND last_name LIKE $%d", len(args)+1)
		args = append(args, fmt.Sprintf("%%%s%%", fUser.LastName))
	}
	if fUser.Role != "" {
		filter += fmt.Sprintf(" AND role = $%d", len(args)+1)
		args = append(args, fUser.Role)
	}

	var totalCount int32
	q := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL` + filter
	err := u.db.QueryRow(q, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	query += filter
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", fUser.Limit, (fUser.Page-1)*fUser.Limit)

	rows, err := u.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var users []*pb.UserProfile
	for rows.Next() {
		var user pb.UserProfile
		err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetUsersListResp{
		Users:      users,
		TotalCount: totalCount,
		Limit: fUser.Limit,
		Page:   fUser.Page,
	}, nil
}
