package types

type User struct {
	ID      int64
	FName   string
	City    string
	Phone   string
	Height  float64
	Married bool
}

func NewUser(id int64, fname, city, phone string, height float64, married bool) *User {
	return &User{
		ID:      id,
		FName:   fname,
		City:    city,
		Phone:   phone,
		Height:  height,
		Married: married,
	}
}
