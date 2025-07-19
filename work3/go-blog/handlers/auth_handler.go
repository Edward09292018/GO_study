package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-blog/entity"
	"go-blog/middle"
	"go-blog/utils"
	"net/http"
)

var SendSuccess = utils.SendSuccess
var HandleError = utils.HandleError

func Register(c *gin.Context) {
	var user entity.User
	if err := c.BindJSON(&user); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	hashedPassword, err := middle.HashPassword(user.Password)
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	db := utils.GetDB()
	if err := db.Create(&user).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, nil, "register successfully", http.StatusOK)
}
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	var user entity.User
	db := utils.GetDB()
	if err := db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		HandleError(c, errors.New("invalid username or password"), http.StatusUnauthorized)
		return
	}

	if !middle.CheckPasswordHash(loginData.Password, user.Password) {
		HandleError(c, errors.New("invalid username or password"), http.StatusUnauthorized)
		return
	}

	token, err := middle.GenerateJWT(user.Username)
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, token, "Login successful", http.StatusOK)
}
