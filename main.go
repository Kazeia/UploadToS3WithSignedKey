package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// file path
	filePathToUpload := os.Args[0]
	keyOrNameOfFileToUploadS3 := os.Args[1]
	bucket := os.Args[2]

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}
	cfg.Region = "us-east-2"
	s3Svc := s3.New(cfg)

	key := keyOrNameOfFileToUploadS3
	file, err := os.Open(filePathToUpload)
	if err != nil {
		log.Fatal(err)
	}
	f, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	//PUT
	fmt.Println("Received request to presign PutObject for,", key)
	sdkReq := s3Svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),

		// If ContentLength is 0 the header will not be included in the signature.
		ContentLength: aws.Int64(f.Size()),
	})
	u, signedHeaders, _ := sdkReq.PresignRequest(15 * time.Minute)
	log.Println("PUT: ")
	log.Println(u)
	log.Println(signedHeaders)

	//GET
	fmt.Println("Received request to presign GetObject for,", key)
	sdkReqx := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	u, signedHeaders, err = sdkReqx.PresignRequest(15 * time.Minute)
	log.Println("GET: ")
	log.Println(u)
	log.Println(signedHeaders)
}
