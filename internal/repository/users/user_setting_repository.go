package repository

import (
	"asset-service/internal/models/user"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

type UserSettingRepository interface {
	GetUserSettingBySettingID(settingID uint) (*user.Setting, error)
	GetUserSettingByUserID(userID uint) (*user.Setting, error)
	GetAllUserSettings() ([]user.Setting, error)
	UserSettingExists(userID uint) (bool, error)
}

type userSettingRepository struct {
	db gorm.DB
}

func NewUserSettingRepository(db gorm.DB) UserSettingRepository {
	return &userSettingRepository{db: db}
}

func (r *userSettingRepository) GetUserSettingBySettingID(settingID uint) (*user.Setting, error) {
	var setting user.Setting
	if err := r.db.Table(utils.TableUserSettingName).First(&setting, "setting_id = ?", settingID).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *userSettingRepository) GetUserSettingByUserID(userID uint) (*user.Setting, error) {
	var setting user.Setting
	if err := r.db.Table(utils.TableUserSettingName).First(&setting, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *userSettingRepository) GetAllUserSettings() ([]user.Setting, error) {
	var settings []user.Setting
	if err := r.db.Table(utils.TableUserSettingName).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *userSettingRepository) UserSettingExists(userID uint) (bool, error) {
	var count int64
	err := r.db.Table(utils.TableUserSettingName).Model(&user.Setting{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}
