package utils

import (
	"fmt"
	"lamhat/core"
	"log"
	"mime/multipart"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Connect to a minIo server
func connectToMinIo(ctx *gin.Context) *minio.Client {
	endpoint := core.Config.FILE_STORAGE.ENDPOINT
	accessKeyID := core.Config.FILE_STORAGE.ACCESS_KEY_ID
	secretAccessKey := core.Config.FILE_STORAGE.ACCESS_SECRET
	useSSL := core.Config.FILE_STORAGE.SSL
	core.Sugar.Debugf("%s-%s-%s-%s", endpoint, accessKeyID, secretAccessKey, useSSL)
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	core.Sugar.Debug("MinIo client created")
	return minioClient
}

// Define the struct with exported fields
type UploadStatus struct {
	bucketname string
	objectname string
	status     bool
}

func UploadIntoGallery(ctx *gin.Context, gallery_id int, user_id int) ([]UploadStatus, error) {
	client := connectToMinIo(ctx)
	// user-id + gallery-id will be the bucket name, Create the bucket if its not exist
	// Check the bucket exist
	var bucket_name string = fmt.Sprintf("gallery-bucket-%v-%v", gallery_id, user_id)

	err := client.MakeBucket(ctx, bucket_name, minio.MakeBucketOptions{Region: core.Config.FILE_STORAGE.LOCATION})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, bucket_name)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucket_name)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucket_name)
	}
	// file upload logic
	// fire go routines to upload files
	// after upload pass the file path into a slice via channal
	// once all uploaded, collect all values from channel and insert into db

	// Get the files from multipart header
	form, err := ctx.MultipartForm()

	if err != nil {
		core.Sugar.Errorf("Error while Uploading File into gallery %s", err.Error())
		return nil, err
	}

	files := form.File["files"]

	// Create a channel
	ch := make(chan UploadStatus, len(files))
	var wg sync.WaitGroup

	for _, file := range files {
		// Upload into bucket
		wg.Add(1)
		go putImageInBucket(ctx, bucket_name, file, client, ch, &wg)
	}
	// wait for the goroutines
	wg.Wait()
	close(ch)
	// All done.
	// Return the struct
	var statusStruct []UploadStatus
	for each := range ch {
		// Append the received value from the channel to the statusStruct slice
		statusStruct = append(statusStruct, each)

	}

	return statusStruct, nil
}

func putImageInBucket(ctx *gin.Context, bucket_name string, file *multipart.FileHeader, client *minio.Client, ch chan UploadStatus, wg *sync.WaitGroup) {

	var upload_status UploadStatus // struct to store upload status

	object_name := fmt.Sprintf("%s-%s", uuid.NewString(), file.Filename) // uuid + file name (to makesure file name is unique)
	upload_status.bucketname = bucket_name
	upload_status.objectname = object_name

	// open file
	reader, err := file.Open()
	if err != nil {
		upload_status.status = false // failed
		ch <- upload_status
		defer wg.Done()
		core.Sugar.Errorf("Error in file %s", file.Filename)
	}
	defer reader.Close()

	n, err := client.PutObject(ctx, bucket_name, object_name, reader, file.Size, minio.PutObjectOptions{ContentType: "application/image"})

	if err != nil {
		upload_status.status = false // failed
		ch <- upload_status
		defer wg.Done()
	}
	fmt.Printf("upaload status %v", n)
	upload_status.status = true // success
	ch <- upload_status
	defer wg.Done()
}
