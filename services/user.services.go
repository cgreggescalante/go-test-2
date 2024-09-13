package services

import (
	"golang.org/x/crypto/bcrypt"
	"nff-go-htmx/db"
)

type User struct {
	ID        int64
	Email     string
	Password  string
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type UserServices struct {
	UserStore db.Store
}

func NewUserService(userStore db.Store) *UserServices {
	return &UserServices{
		UserStore: userStore,
	}
}

func (us *UserServices) CreateUser(user User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return 0, err
	}

	user.Password = string(hashedPassword)

	statement := `INSERT INTO users (email, password, first_name, last_name) VALUES (:email, :password, :first_name, :last_name)`

	result, err := us.UserStore.Db.NamedExec(statement, user)
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	return userId, err
}

func (us *UserServices) CheckEmail(email string) (User, error) {
	query := `SELECT id, email, password, first_name, last_name FROM users WHERE email = ?`

	var user User

	if err := us.UserStore.Db.Get(&user, query, email); err != nil {
		return User{}, err
	}

	return user, nil
}
