package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func loadToS3(imagePaths []string, endpointPaths []string) error {
	return loadToDigitalOceanS3(imagePaths, endpointPaths)
}

func loadToDigitalOceanS3(imagePaths []string, endpointPaths []string) error {
	// Configuration
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			myEnvironment["S3_ACCESS_KEY_ID"],
			myEnvironment["S3_SECRET_ACCESS_KEY"],
			"",
		),
		Endpoint:         aws.String(myEnvironment["S3_ENDPOINT"]),
		S3ForcePathStyle: aws.Bool(false),
		Region:           aws.String(myEnvironment["S3_REGION"]),
	})
	if err != nil {
		log.Println("Error creating the session: ", err)
		return err
	}

	// Create S3 service client
	var client = s3.New(sess)

	// Loop through imagePaths and upload each file to S3
	for i := 0; i < len(imagePaths); i++ {
		// Load file
		file, err := os.Open(imagePaths[i])
		if err != nil {
			log.Println("Error opening the file: ", err)
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println(err)
			}
		}(file)

		// Prepare the S3 upload parameters
		var bucketName = aws.String(myEnvironment["S3_BUCKET_NAME"])
		params := &s3.PutObjectInput{
			Bucket:   bucketName,
			Key:      aws.String(endpointPaths[i]),
			Body:     file,
			ACL:      aws.String("public-read"),
			Metadata: map[string]*string{},
		}

		// Perform the upload
		_, err = client.PutObject(params)
		if err != nil {
			log.Println("Error uploading the file: ", err)
			return err
		}
	}

	return nil
}
