package services_test

import (
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
	"videoencoder/application/repositories"
	"videoencoder/application/services"
	"videoencoder/domain"
	"videoencoder/framework/database"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func prepare() (domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()

	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}

	return *video, repo
}

func TestVideoServiceDownload(t *testing.T) {
	video, repo := prepare()
	videoService := services.NewVideoService()
	videoService.Video = &video
	videoService.VideoRepository = &repo

	err := videoService.Download("bucketName")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)
}
