package user

import "github.com/HiroLiang/goat-server/internal/domain/user"

func toDomain(record *UserRecord) (*user.User, error) {
	return &user.User{
		ID:        record.ID,
		Name:      record.Name,
		Email:     record.Email,
		Password:  record.Password,
		Status:    record.UserStatus,
		LastIP:    record.UserIP,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}, nil
}

func toRecord(user *user.User) *UserRecord {
	return &UserRecord{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Password:   user.Password,
		UserStatus: user.Status,
		UserIP:     user.LastIP,
	}
}
