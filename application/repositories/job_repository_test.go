package repositories_test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"videoencoder/application/repositories"
	"videoencoder/domain"
	"videoencoder/framework/database"
)

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()

	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(*video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)
	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(*job)

	j, err := repoJob.Find(job.ID)
	require.Nil(t, err)
	require.Equal(t, job.ID, j.ID)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()

	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(*video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)
	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(*job)

	job.Status = "Complete"

	repoJob.Update(*job)

	j, err := repoJob.Find(job.ID)
	require.Nil(t, err)
	require.Equal(t, job.ID, j.ID)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.VideoID, video.ID)
	require.Equal(t, j.Status, job.Status)
}
