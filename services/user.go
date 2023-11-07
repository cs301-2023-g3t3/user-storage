package services

import (
	"errors"
	"net/http"
	"user-storage/models"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var validate = validator.New()

type UserService struct {
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{DB: db}
}


func (t *UserService) GetAllUsers(users *[]models.User) (int, error) {
	err := t.DB.Find(&users)
	if err.Error != nil {
		return http.StatusInternalServerError, err.Error
	}
    return http.StatusOK, nil
}

func (t *UserService) GetUserByID(user *models.User, id string) (int, error) {
	if id == "" {
		return http.StatusBadRequest, errors.New("User ID cannot be empty")
	}
    res := t.DB.Find(&user, "id = ?", id)

	if res.Error != nil {
        return http.StatusInternalServerError, res.Error
	}

    if res.RowsAffected == 0 {
        return http.StatusNotFound, errors.New("User ID is not found")
    }
    
    return http.StatusOK, nil
}

func (t *UserService) AddUser(user *models.User) (int, error) {
    if err := validate.Struct(user); err != nil {
        return http.StatusBadRequest, err
    }

    user.Id = uuid.NewString()

    res := t.DB.Create(&user)

    if res.Error != nil {
        return http.StatusInternalServerError, res.Error
    }
    
    return http.StatusCreated, nil
}

func (t *UserService) UpdateUserById(user *models.User, id string) (int, error) {
	if id == "" {
		return http.StatusBadRequest, errors.New("User ID cannot be empty")
	}

    if err := validate.Struct(user); err != nil {
        return http.StatusBadRequest, err
    }

    existingUser := models.User{}
    if err := t.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return http.StatusNotFound, errors.New("User not found with given ID")
        } else {
            return http.StatusInternalServerError, err
        }
    }

    user.Id = id

    // Update the user's data
    res := t.DB.Model(models.User{Id: id}).Updates(&user).Error

    if res != nil {
        return http.StatusInternalServerError, res
    }

    return http.StatusOK, nil
}

func (t *UserService) DeleteUserById (existingUser *models.User, id string) (int, error) {
    if id == "" {
        return http.StatusBadRequest, errors.New("User ID cannot be empty")
    }

    if err := t.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return http.StatusNotFound, errors.New("User not found with given ID")
        } else {
            return http.StatusInternalServerError, err
        }
    }

    err := t.DB.Where("id = ?", id).Delete(&existingUser).Error
    if err != nil {
        return http.StatusInternalServerError, err
    }

    return http.StatusOK, nil
}
