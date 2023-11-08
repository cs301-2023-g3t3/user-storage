package services

import (
	"errors"
	"net/http"
	"user-storage/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

type UserService struct {
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{DB: db}
}


func (t *UserService) GetAllUsers() (*[]models.User, int, error) {
	var users []models.User
	err := t.DB.Find(&users)
	if err.Error != nil {
		return nil, http.StatusInternalServerError, err.Error
	}
    return &users, http.StatusOK, nil
}

func (t *UserService) GetUserByID(id string) (*models.User, int, error) {
	var user models.User
	if id == "" {
		return nil, http.StatusBadRequest, errors.New("User ID cannot be empty")
	}
    res := t.DB.Find(&user, "id = ?", id)

	if res.Error != nil {
        return nil, http.StatusInternalServerError, res.Error
	}

    if res.RowsAffected == 0 {
        return nil, http.StatusNotFound, errors.New("User ID is not found")
    }
    
    return &user, http.StatusOK, nil
}

func (t *UserService) AddUser(user *models.User) (*models.User, int, error) {
    if err := validate.Struct(user); err != nil {
        return nil, http.StatusBadRequest, err
    }

    tx := t.DB.Begin()
    err := tx.Create(&user).Error
    if err != nil {
        tx.Rollback()
        return nil, http.StatusInternalServerError, err
    }
    
    return user, http.StatusCreated, nil
}

func (t *UserService) UpdateUserById(user *models.User, id string) (*models.User, int, error) {
	if id == "" {
		return nil, http.StatusBadRequest, errors.New("User ID cannot be empty")
	}

    if err := validate.Struct(user); err != nil {
        return nil, http.StatusBadRequest, err
    }

    existingUser := models.User{}
    if err := t.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, http.StatusNotFound, errors.New("User not found with given ID")
        } else {
            return nil, http.StatusInternalServerError, err
        }
    }

    user.Id = id

    // Update the user's data
    err := t.DB.Model(models.User{Id: id}).Updates(&user).Error

    if err != nil {
        return nil, http.StatusInternalServerError, err
    }

    return user, http.StatusOK, nil
}

func (t *UserService) DeleteUserById (id string) (*models.User, int, error) {
    if id == "" {
        return nil, http.StatusBadRequest, errors.New("User ID cannot be empty")
    }

    var existingUser models.User
    if err := t.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, http.StatusNotFound, errors.New("User not found with given ID")
        } else {
            return nil, http.StatusInternalServerError, err
        }
    }

    err := t.DB.Where("id = ?", id).Delete(&existingUser).Error
    if err != nil {
        return nil, http.StatusInternalServerError, err
    }

    return &existingUser, http.StatusOK, nil
}

func (t *UserService) GetUsersWithRole(roles []int) (*[]models.User, int, error) {
	var users []models.User
    err := t.DB.Where("role IN ?", roles).Find(&users).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, http.StatusNotFound, errors.New("Cannot find users with given roles")
        } else {
            return nil, http.StatusInternalServerError, err
        }
    }

    return &users, http.StatusOK, nil
}
