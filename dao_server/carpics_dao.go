/*
   Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at
    http://aws.amazon.com/apache2.0/
   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/

package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "fmt"
    "os"
)

func getCarPics(car_img string) () {

  bucket := "cloudhydras3"
  item := car_img

  file, err := os.Create(item)
  if err != nil {
      //exitErrorf("Unable to open file %q, %v", err)
      return
  }
  defer file.Close()

  // Initialize a session in us-west-2 that the SDK will use to load
  // credentials from the shared credentials file ~/.aws/credentials.
  sess, _ := session.NewSession(&aws.Config{
      Region: aws.String("us-west-2")},
  )

  downloader := s3manager.NewDownloader(sess)
  numBytes, err := downloader.Download(file,
      &s3.GetObjectInput{
          Bucket: aws.String(bucket),
          Key:    aws.String(item),
      })
  if err != nil {
      //exitErrorf("Unable to download item %q, %v", item, err)
      return
  }
  fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return
}

func uploadCarPics(car_img string) () {

  bucket := "cloudhydras3"
  filename := car_img

  file, err := os.Open(filename)
  if err != nil {
      //exitErrorf("Unable to open file %q, %v", err)
      return
  }
  defer file.Close()

  // Initialize a session in us-west-2 that the SDK will use to load
  // credentials from the shared credentials file ~/.aws/credentials.
  sess, err := session.NewSession(&aws.Config{
      Region: aws.String("us-west-2")},
  )

  uploader := s3manager.NewUploader(sess)
  // Upload the file's body to S3 bucket as an object with the key being the filename
  _, err = uploader.Upload(&s3manager.UploadInput{
      Bucket: aws.String(bucket),
      Key: aws.String(filename),
      Body: file,
  })
  if err != nil {
      //exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
      return
  }

  fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
	return
}
