package services_test

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"videoencoder/application/services"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()
	videoService := services.NewVideoService()
	videoService.Video = &video
	videoService.VideoRepository = &repo

	err := videoService.Download("bucketName")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "bucketName"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID
	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)
	result := <-doneUpload
	require.Equal(t, result, "upload completed")
	err = videoService.Finish()
	require.Nil(t, err)
}
