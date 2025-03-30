package services

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"os/exec"
	"videoencoder/application/repositories"
	"videoencoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	if err := minioClient.FGetObject(ctx, bucketName, v.Video.FilePath, os.Getenv("localStoragePath")+"/"+v.Video.ID+".mp4", minio.GetObjectOptions{}); err != nil {
		log.Fatalln(err)
		return err
	}

	log.Printf("Downloaded file: %s", v.Video.ID+".mp4")

	return nil

}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	printOutput(output)
	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	printOutput(output)
	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		log.Println("Erro removing mp4 ", v.Video.ID+".mp4")
		return err
	}
	err = os.Remove(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag")

	if err != nil {
		log.Println("Erro removing frag ", v.Video.ID+".frag")
		return err
	}
	err = os.RemoveAll(os.Getenv("localStoragePath") + "/" + v.Video.ID)

	if err != nil {
		log.Println("Erro removing folder ", v.Video.ID)
		return err
	}

	log.Println("Files have been removed: ", v.Video.ID)
	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n\n", string(out))
	}
}
