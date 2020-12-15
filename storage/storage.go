package storage

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gabihodoroaga/smtpd-email-forward/config"
	"google.golang.org/api/option"
)

type storageConnection struct {
	Client *storage.Client
}

var (
	client *storageConnection
	once   sync.Once
)

// UploadFileToGCSBucket uploads the email message to GCP bucket, one folder for each mailbox
func UploadFileToGCSBucket(ctx context.Context, mailbox string, data []byte) (string, error) {
	filename := generateFileNameForGCS(ctx, "email")
	filepath := fmt.Sprintf("%s/%s.eml", mailbox, filename)
	client, err := getGCSClient(ctx)
	if err != nil {
		return "", err
	}
	wc := client.Bucket(config.Config.GCPBucket).Object(filepath).NewWriter(ctx)
	if _, err = wc.Write(data); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return filepath, nil
}

func generateFileNameForGCS(ctx context.Context, name string) string {
	time := time.Now().UnixNano()
	var strArr []string
	strArr = append(strArr, name)
	strArr = append(strArr, strconv.Itoa(int(time)))
	var filename string
	for _, str := range strArr {
		filename = filename + str
	}
	return filename
}

func getGCSClient(ctx context.Context) (*storage.Client, error) {
	var clientErr error
	once.Do(func() {
		storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(config.Config.GCPCredentials))
		if err != nil {
			clientErr = fmt.Errorf("Failed to create GCS client ERROR:%s", err.Error())
		} else {
			client = &storageConnection{
				Client: storageClient,
			}
		}
	})
	return client.Client, clientErr
}
