package domain

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Job struct {
	ID               string    `valid:"uuid" gorm:"type:uuid;primary_key" json:"job_id"`
	OutputBucketPath string    `valid:"notnull" json:"output_bucket_path"`
	Status           string    `valid:"notnull" json:"status"`
	Video            *Video    `valid:"-" json:"video"`
	VideoID          string    `valid:"-" gorm:"column:video_id;type:uuid;notnull" json:"-"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `valid:"-" json:"created_at"`
	UpdatedAt        time.Time `valid:"-" json:"updated_at"`
}

func NewJob(output string, status string, video *Video) (*Job, error) {
	job := &Job{
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
	}

	job.prepare()

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (job *Job) prepare() {
	job.ID = uuid.NewV4().String()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)
	if err != nil {
		return err
	}
	return nil
}
