package services

import (
	"go-test-2/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type UserServices struct {
	User      User
	UserStore db.Store
}

func NewUserService(u User, userStore db.Store) *UserServices {
	return &UserServices{
		User:      u,
		UserStore: userStore,
	}
}

func (us *UserServices) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	statement := `INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)`

	_, err = us.UserStore.Db.Exec(statement, u.Email, string(hashedPassword), u.FirstName, u.LastName)

	return err
}

func (us *UserServices) CheckEmail(email string) (User, error) {
	query := `SELECT id, email, password, first_name, last_name FROM users WHERE email = ?`

	statement, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer statement.Close()

	us.User.Email = email
	err = statement.QueryRow(us.User.Email).Scan(&us.User.ID, &us.User.Email, &us.User.Password, &us.User.FirstName, &us.User.LastName)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
}
