package utils

import (
	"fmt"
	"lamhat/core"
	"lamhat/model"
	"log"
	"mime/multipart"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Define the struct with exported fields
type UploadStatus struct {
	Bucketname string
	Objectname string
	Status     bool
	Gallery_id int
}

// Connect to a minIo server
func connectToMinIo() *minio.Client {
	endpoint := core.Config.FILE_STORAGE.ENDPOINT
	accessKeyID := core.Config.FILE_STORAGE.ACCESS_KEY_ID
	secretAccessKey := core.Config.FILE_STORAGE.ACCESS_SECRET
	useSSL := core.Config.FILE_STORAGE.SSL

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

func UploadIntoGallery(ctx *gin.Context, gallery_id int, user_id int) ([]UploadStatus, error) {
	client := connectToMinIo()
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
		go putImageInBucket(ctx, bucket_name, gallery_id, file, client, ch, &wg)
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

func putImageInBucket(ctx *gin.Context, bucket_name string, gallery_id int, file *multipart.FileHeader, client *minio.Client, ch chan UploadStatus, wg *sync.WaitGroup) {

	var upload_status UploadStatus // struct to store upload status

	object_name := fmt.Sprintf("%s-%s", uuid.NewString(), file.Filename) // uuid + file name (to makesure file name is unique)
	upload_status.Bucketname = bucket_name
	upload_status.Objectname = object_name
	upload_status.Gallery_id = gallery_id

	// open file
	reader, err := file.Open()
	if err != nil {
		upload_status.Status = false // failed
		ch <- upload_status
		defer wg.Done()
		core.Sugar.Errorf("Error in file %s", file.Filename)
	}
	defer reader.Close()

	n, err := client.PutObject(ctx, bucket_name, object_name, reader, file.Size, minio.PutObjectOptions{ContentType: "application/image"})

	if err != nil {
		upload_status.Status = false // failed
		ch <- upload_status
		defer wg.Done()
	}
	fmt.Printf("upaload status %v", n)
	upload_status.Status = true // success
	ch <- upload_status
	defer wg.Done()
}

func GetPreSignedURL(ctx *gin.Context, data []UploadStatus) []model.PreSignedURLS {
	client := connectToMinIo()

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)

	var allUrl []model.PreSignedURLS

	for _, each := range data {

		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", each.Objectname))

		url, err := client.PresignedGetObject(ctx,
			each.Bucketname, each.Objectname,
			time.Duration(core.Config.FILE_STORAGE.PRESIGNED_DELAY_MNT*int(time.Second)),
			reqParams)

		if err != nil {
			core.Sugar.Errorf("Error while retrieving preSigned URL %s", err.Error())
			continue // Dont append if its error
		}

		allUrl = append(allUrl, model.PreSignedURLS{URL: url.String()})
	}
	return allUrl

}
