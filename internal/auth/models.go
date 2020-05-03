package auth

import "golang.org/x/crypto/bcrypt"

// NewUser returns a user with email, pwd and org.
func NewUser(email, password, org string) (*User, error) {
	user := User{
		Email:        email,
		Organization: org,
	}
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// User defines the user model.
type User struct {
	Email        string
	PasswordHash []byte
	Organization string
}

// SetPassword sets password hash with bcrypt.
func (u *User) SetPassword(pwd string) error {
	pwdBytes := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

// VerifyPassword verifies password with bcrypt.
func (u *User) VerifyPassword(pwd string) error {
	pwdBytes := []byte(pwd)
	return bcrypt.CompareHashAndPassword(u.PasswordHash, pwdBytes)
}

// Organization defines the org model.
type Organization struct {
	Name       string
	AdminEmail string
}
