package domain

import (
	"github.com/asaskevich/govalidator"
	"time"
)

type Video struct {
	ID         string    `gorm:"type:uuid;primary_key" valid:"uuid" json:"encoded_video_folder"`
	ResourceID string    `gorm:"type:varchar(255)" valid:"notnull" json:"resource_id"`
	FilePath   string    `gorm:"type:varchar(255)" valid:"notnull" json:"file_path"`
	CreatedAt  time.Time `valid:"-" json:"-"`
	Jobs       []*Job    `gorm:"ForeignKey:VideoID" valid:"-" json:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewVideo() *Video {
	return &Video{}
}

func (video *Video) Validate() error {
	_, err := govalidator.ValidateStruct(video)
	if err != nil {
		return err
	}
	return nil
}
