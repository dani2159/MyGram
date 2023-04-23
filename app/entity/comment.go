package entity

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

// Comment represents Comment
type Comment struct {
	Base
	UserID  uint   `json:"user_id" example:"2"`
	PhotoID uint   `json:"photo_id" form:"photo_id" example:"3"`
	Message       string `gorm:"not null" json:"message" form:"message" valid:"required~Comment is required"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}
	return nil
}

func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(c)

	if errUpdate != nil {
		err = errUpdate
		return
	}
	return nil
}
