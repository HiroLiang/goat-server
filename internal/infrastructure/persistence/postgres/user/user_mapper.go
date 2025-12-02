package user

import "github.com/HiroLiang/goat-server/internal/domain/user"

func toDomain(rec *UserRecord) (*user.User, error) {
	email, err := user.NewEmail(rec.Email)
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:        user.ID(rec.ID),
		Name:      rec.Name,
		Email:     email,
		Password:  rec.Password,
		Status:    user.Status(rec.UserStatus),
		LastIP:    rec.UserIP,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}, nil
}

func toRecord(u *user.User) *UserRecord {
	return &UserRecord{
		ID:         int64(u.ID),
		Name:       u.Name,
		Email:      string(u.Email),
		Password:   u.Password,
		UserStatus: StatusDB(u.Status),
		UserIP:     u.LastIP,
	}
}
