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
	UserStore db.Store
}

func NewUserService(userStore db.Store) *UserServices {
	return &UserServices{
		UserStore: userStore,
	}
}

func (us *UserServices) CreateUser(u User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return 0, err
	}

	statement := `INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)`

	result, err := us.UserStore.Db.Exec(statement, u.Email, string(hashedPassword), u.FirstName, u.LastName)
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	return userId, err
}

func (us *UserServices) CheckEmail(email string) (User, error) {
	query := `SELECT id, email, password, first_name, last_name FROM users WHERE email = ?`

	statement, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer statement.Close()

	user := User{}

	user.Email = email
	err = statement.QueryRow(user.Email).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
