package services

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *minio.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")
	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()
	contentType := "application/octet-stream"
	_, err = client.FPutObject(ctx, vu.OutputBucket, path[1], objectPath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)
	err := vu.loadPaths()
	if err != nil {
		return err
	}
	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}
	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(in chan int, returnChannel chan string, uploadClient *minio.Client, ctx context.Context) {
	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)
		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("error during the upload: %v. Error: %v", vu.Paths[x], err)
			returnChannel <- err.Error()
		} else {
			returnChannel <- ""
		}
	}

	returnChannel <- "upload finished"
}

func getClientUpload() (*minio.Client, context.Context, error) {
	ctx := context.Background()
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, nil, err
	}
	return minioClient, ctx, nil
}
