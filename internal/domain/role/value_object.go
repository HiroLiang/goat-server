package role

import "encoding/json"

type ID int64

type Type string

const (
	Admin  Type = "admin"
	Vendor Type = "vendor"
	User   Type = "user"
	Guest  Type = "guest"
)

func typeFrom(s string) (Type, error) {
	switch Type(s) {
	case Admin, Vendor, User, Guest:
		return Type(s), nil
	default:
		return "", ErrInvalidType
	}
}

func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	v, err := typeFrom(s)
	if err != nil {
		return err
	}

	*t = v
	return nil
}
