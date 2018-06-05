package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	bucketName = "slidehub-slides01"
	key        = "051db4bdb651a0840f01224429aed949"
)

func exists(downloadFile string, resp *s3.ListObjectsOutput) bool {
	for _, content := range resp.Contents {
		if *content.Key == downloadFile {
			return true
		}
	}
	return false
}

func DownloadFile(filename string) error {
	// download実行
	sess := (session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	downloader := s3manager.NewDownloader(sess)
	os.MkdirAll("./data/tmp/download", 0755)
	f, _ := os.Create("./data/tmp/download/" + filename)
	defer f.Close()

	_, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to Download file, %v", err)
	}
	return nil
}

func FileCheck() error {
	sess := session.Must(session.NewSession())
	svc := s3.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(key),
	})
	if err != nil {
		return err
	}
	if !exists(key, resp) {
		return fmt.Errorf("バケット内にダウンロードするファイルが存在しない")
	}
	return nil
}
