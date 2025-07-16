package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	if err = db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user create success",
	})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	var storedUser User
	if err := db.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid username"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
