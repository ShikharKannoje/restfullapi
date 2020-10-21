package controllers

import (
	"Restfull/api/auth"
	"Restfull/api/models"
	"Restfull/api/responses"
	"Restfull/api/utils/formaterror"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	loginUser := models.LoginUser{}

	err = json.Unmarshal(body, &loginUser)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	//fmt.Println(user.Email, user.Password)
	err = loginUser.ValidateLogin()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	//err = loginUser.BeforeSigin()
	if err != nil {
		fmt.Println("cannot hash the password")
		fmt.Println(err)
	}
	fmt.Println(loginUser.Email, loginUser.Password)
	token, err := server.SignIn(loginUser.Email, loginUser.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.LoginUser{}

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err := sql.Open(os.Getenv("DB_DRIVER"), DBURL)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `SELECT userid, password_hash FROM "Password" WHERE userid = (SELECT userid FROM "User" WHERE email = $1)`
	var id int
	var pass_hash = ""
	err = db.QueryRow(statement, email).Scan(&id, &pass_hash)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(pass_hash, user.Password)
	err = models.VerifyPassword(pass_hash, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(uint32(id))
}
