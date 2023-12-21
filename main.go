package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {

	ctx := context.Background()
	endpoint := os.Getenv("MINIO_ENDPOINT") // based on what you defined at https://github.com/govindkailas/upload-to-minio/blob/91913c2f84863f3f172f69ea0fab051d6d3601c0/k8s/minio.yaml#L94 
	accessKeyID := os.Getenv("MINIO_ACCESSKEY")
	secretAccessKey := os.Getenv("MINIO_SECRETKEY")
	useSSL := false
	// check if endpoint,acccessKeyID and secretAccessKey are set
	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Fatalln("Error: MINIO_ENDPOINT, MINIO_ACCESSKEY, MINIO_SECRETKEY environment variables are not set")
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	// Make a new bucket called testbucket.
	bucketName := "testbucket"
	location := "us-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Add another handler called /ping to ping the server
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "Yooh, I am alive responding pong!!"}`)

	})

	// Start HTTP server
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		// file, _, err := r.FormFile("file")
		file, header, err := r.FormFile("file") //to upload use this, curl http://localhost:8080/upload -F 'file=@lifecycle.png'
		if err != nil {
			log.Println("Error retrieving file:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		fmt.Printf("File type is %s\n", strings.Split(header.Filename, ".")[1])
		fileName := fmt.Sprintf("file_%d", time.Now().UnixNano())            //adding the timestamp to the filename
		finalname := fileName + "." + strings.Split(header.Filename, ".")[1] // appending the extension to the filename

		// Create a file to write our uploaded file to.
		f, err := os.OpenFile(finalname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		io.Copy(f, file)
		log.Println("File saved to", finalname)

		filePath := "./" + finalname
		contentType := "application/octet-stream"
		log.Println("Uploading file to minio")
		// Upload the test file with FPutObject
		info, err := minioClient.FPutObject(ctx, bucketName, finalname, filePath, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			log.Fatalln(err)
		}

		fileSizeKB := info.Size / 1024
		log.Printf("Successfully uploaded %s of size %d KB\n", finalname, fileSizeKB)

	})

	http.ListenAndServe(":8080", nil)
}
