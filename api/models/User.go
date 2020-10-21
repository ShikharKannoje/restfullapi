package models

import (
	"context"
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	"html"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/badoux/checkmail"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Userid   uint32 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Phoneno  string `json:"phoneno"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
func (u *LoginUser) BeforeSigin() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.Userid = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Address = html.EscapeString(u.Address)
	u.Phoneno = html.EscapeString(strings.TrimSpace(u.Phoneno))
}

func (u *LoginUser) ValidateLogin() error {
	if u.Password == "" {
		return errors.New("Required Password")
	}
	if u.Email == "" {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("Invalid Email")
	}
	return nil
}

func PasswordValidator(s string) bool {
	var (
		upp, num, sym bool
	)
	fmt.Println(upp, num, sym)
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			fmt.Println(char)
			upp = true
		case unicode.IsNumber(char):
			fmt.Println(char)
			num = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			fmt.Println(char)
			sym = true
		}
	}
	fmt.Println(upp, num, sym)
	if !upp || !num || !sym {
		return false
	}
	return true
}

func (u *User) Validate() error {

	if u.Username == "" {
		return errors.New("Required Username")
	}
	if len(u.Username) > 10 {
		return errors.New("Username too long, should be less than 10 character")
	}

	f := PasswordValidator(u.Password)
	if f == false {
		return errors.New("Password does not satisfy the required contions:Password should contain one Capital letter, one special character, and a number ")
	}

	if u.Email == "" {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("Invalid Email")
	}

	if u.Address == "" {
		return errors.New("Required Address")
	}

	if u.Phoneno == "" {
		return errors.New("Required Phoneno")
	}
	return nil
}

func (u *User) SaveUser() (*User, error) {

	var err error
	err = u.BeforeSave()
	if err != nil {
		fmt.Println("Password was not hashed")
		fmt.Println(err)
	}
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err := sql.Open(os.Getenv("DB_DRIVER"), DBURL)
	if err != nil {
		log.Println(err)
		return u, err
	} else {
		log.Println("Database connected")
	}
	defer db.Close()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	statement := `INSERT into "User"(username, email, address, phoneno) VALUES($1, $2, $3, $4) RETURNING username, userid`
	name := ""
	var id int
	err = tx.QueryRow(statement, u.Username, u.Email, u.Address, u.Phoneno).Scan(&name, &id)

	if err != nil {
		tx.Rollback()
		return u, err
	}
	statement = `INSERT into "Password"(userid, password_hash, status) VALUES($1, $2, $3) RETURNING userid`
	err = tx.QueryRow(statement, id, u.Password, true).Scan(&id)
	if err != nil {
		tx.Rollback()
		return u, err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New User is : ", name, "with userID", id)
	return u, nil
}
