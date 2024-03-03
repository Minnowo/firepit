package database

import (
	"github.com/EZCampusDevs/firepit/data"
	"github.com/EZCampusDevs/firepit/database/models"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

func IsCredentialsValid(credentials *data.AuthPayload) bool {

	var user models.User

	err := db.Where("username = ?", credentials.Username).First(&user).Error

	if err != nil {

		log.Warn(err)

		return false
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(credentials.Password))

	if err != nil {

		log.Debugf("Password hash missmatch for user %s", credentials.Username)

		return false
	}

	return true
}

func CreateUser(credentials *data.AuthPayload) bool {

	if err := credentials.IsValid(); err != nil {

		log.Warnf("Cannot create new user because invalid credentials: %s", err.Error())

		return false
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)

	if err != nil {

		log.Error(err)

		return false
	}

	var user models.User

	user.Username = credentials.Username
	user.Password = passwordHash

	err = GetDB().Create(&user).Error

	if err != nil {

		log.Warn(err)

		return false
	}

	return true
}
