package user

import (
	"errors"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"golang.org/x/crypto/bcrypt"
)

func Create(
	name string,
	surname string,
	phone string,
	email string,
	password string,
) (*models.User, error) {
	token := utils.RandomString(18)

	// check email or phone exists?
	var found_user models.User
	err := database.DB.Where("email = ? OR phone = ?", email, phone).First(&found_user).Error
	if err == nil { // we found someone
		return nil, errors.New("email or phone used by other user")
	}

	// hash password
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name: name,
		Surname: surname,
		Phone: phone,
		Email: email,
		VerifyToken: token,
		Password: string(hashed_password),
	}

	err = database.DB.Create(&user).Error
	if err != nil {
		return nil ,err
	}

	// create role for the user
	err = database.DB.Create(&models.UserRole{
		UserID: user.ID,
		RoleID: 2, // standart user
	}).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}
