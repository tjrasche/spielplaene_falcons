package excel

import (
	"context"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"rasche-thalhofer.cloud/init/model"
	"rasche-thalhofer.cloud/init/yaml"
	"sort"
)

type Reader struct {
	cfBucket          string
	cfAccId           string
	cfAccessKeyId     string
	cfAccessKeySecret string
	db                *gorm.DB
	gameday           yaml.GameDay
}

func NewReader(cfAccId string, cfAccessKeyId string, cfAccessKeySecret string, db *gorm.DB, gameday yaml.GameDay) *Reader {
	return &Reader{cfBucket: gameday.Bucket, cfAccId: cfAccId, cfAccessKeyId: cfAccessKeyId, cfAccessKeySecret: cfAccessKeySecret, db: db, gameday: gameday}
}

func (r Reader) UpdateGames() (err error) {
	file, err := r.readCurrentExcelFile()
	if err != nil {
		return err
	}
	for _, round := range r.gameday.Rounds {
		roundName, err := file.GetCellValue(round.Worksheet, "A1")
		if err != nil {
			return err
		}
		r.db.Save(&model.Round{ID: roundName})
	}
	return nil
}
func (r Reader) readCurrentExcelFile() (file *excelize.File, err error) {

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r.cfAccId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r.cfAccessKeyId, r.cfAccessKeySecret, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	objects, err := client.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: &r.cfBucket,
	})
	if err != nil {
		return nil, err
	}

	files := objects.Contents
	sort.Slice(files, func(i, j int) bool {
		if files[i].LastModified.Compare(*files[j].LastModified) == 1 {
			return true
		}
		return false
	})

	object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &r.cfBucket,
		Key:    files[0].Key,
	})
	if err != nil {
		return nil, err
	}

	file, err = excelize.OpenReader(object.Body)
	if err != nil {
		return nil, err
	}
	return file, nil
}
