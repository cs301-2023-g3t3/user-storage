package services

import (
	"errors"
	"fmt"
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


func (t *UserService) GetAllUsers(role int, id, name, email string) (*[]models.User, int, error) {
    var users []models.User

    query := t.DB
    if id != "" {
        query = query.Where("id LIKE ?", fmt.Sprint(id,"%"))
    }
    if role != -1 && role != 0{
        query = query.Where("role = ?", role)
    } else if role == 0 {
        query = query.Where("role IS NULL")
    }
    if name != "" {
        query = query.Where("first_name LIKE ? OR last_name LIKE ?", fmt.Sprint(name,"%"), fmt.Sprint(name,"%"))
    }
    if email != "" {
        query = query.Where("email LIKE ?", fmt.Sprint(email,"%"))
    }

    if err := query.Find(&users).Error; err != nil {
        return nil, http.StatusInternalServerError, err
    }
	
    return &users, http.StatusOK, nil
}

func (t *UserService) GetPaginatedUsers(page, pageSize int) (*[]models.User, int, error) {
    var users []models.User

    offset := (page - 1) * pageSize
    err := t.DB.Offset(offset).Limit(pageSize).Find(&users)
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
    err := t.DB.First(&user, "id = ?", id).Error
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, http.StatusNotFound, errors.New("User ID is not found")
        }
        return nil, http.StatusInternalServerError, err
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
    tx.Commit()
    
    return user, http.StatusCreated, nil
}

func (t *UserService) UpdateUserById(user *models.User, id string) (*models.User, int, error) {
	if id == "" {
		return nil, http.StatusBadRequest, errors.New("User ID cannot be empty")
	}

    if err := validate.Struct(user); err != nil {
        return nil, http.StatusBadRequest, err
    }

    tx := t.DB.Begin()
    existingUser := models.User{}
    if err := tx.Where("id = ?", id).First(&existingUser).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, http.StatusNotFound, errors.New("User not found with given ID")
        } else {
            return nil, http.StatusInternalServerError, err
        }
    }

    user.Id = id

    // Update the user's data
    err := tx.Model(models.User{Id: id}).Updates(&user).Error

    if err != nil {
        tx.Rollback()
        return nil, http.StatusInternalServerError, err
    }
    tx.Commit()

    return user, http.StatusOK, nil
}

func (t *UserService) DeleteUserById (id string) (*models.User, int, error) {
    if id == "" {
        return nil, http.StatusBadRequest, errors.New("User ID cannot be empty")
    }

    existingUser, code, err := t.GetUserByID(id)
    if err != nil {
        return nil, code, err
    }

    tx := t.DB.Begin()
 //    var existingUser models.User
 //    err := tx.First(&existingUser, "id = ?", id).Error
	// if err != nil {
 //        tx.Rollback()
 //        if errors.Is(err, gorm.ErrRecordNotFound) {
 //            return nil, http.StatusNotFound, errors.New("User ID is not found")
 //        }
 //        return nil, http.StatusInternalServerError, err
	// }

    err = tx.Where("id = ?", id).Delete(existingUser).Error
    if err != nil {
        tx.Rollback()
        return nil, http.StatusInternalServerError, err
    }
    tx.Commit()

    return existingUser, http.StatusOK, nil
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
