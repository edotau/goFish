package gcloud

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"cloud.google.com/go/storage"
// 	"google.golang.org/api/iterator"
// 	"google.golang.org/api/option"
// )

// // implicit uses Application Default Credentials to authenticate.
// func AuthGcp(jsonPath string) (context.Context, *storage.Client) {
// 	ctx := context.Background()

// 	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return ctx, client
// }

// func GetAllBuckets(jsonPath, projectID string) {
// 	ctx, client := AuthGcp(jsonPath)

// 	it := client.Buckets(ctx, projectID)
// 	defer client.Close()

// 	for buckets, err := it.Next(); err != iterator.Done; buckets, err = it.Next() {
// 		if err == nil {
// 			fmt.Println(buckets.Name)
// 		} else {
// 			log.Fatal(err)
// 		}
// 	}
// }

// // downloadFile downloads an object to a file.
// func DownloadFile(bucket, object string, destFileName string) {
// 	// bucket := "bucket-name"
// 	// object := "object-name"
// 	// destFileName := "file.txt"
// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)

// 	if err != nil {
// 		log.Fatalf("storage.NewClient: %v", err)
// 	}

// 	defer client.Close()

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
// 	defer cancel()

// 	f, err := os.Create(destFileName)

// 	if err != nil {
// 		log.Fatalf("os.Create: %v", err)
// 	}

// 	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)

// 	if err != nil {
// 		log.Fatalf("Object(%q).NewReader: %v", object, err)
// 	}
// 	defer rc.Close()

// 	if _, err := io.Copy(f, rc); err != nil {
// 		log.Fatalf("io.Copy: %v", err)
// 	}

// 	if err = f.Close(); err != nil {
// 		log.Fatalf("f.Close: %v", err)
// 	}
// }

// func ListFiles(bucket string, folder string) {
// 	ans := make(chan string, 248)

// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)

// 	if err != nil {
// 		log.Fatalf("storage.NewClient: %v", err)
// 	}

// 	defer client.Close()

// 	query := &storage.Query{Prefix: folder}

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
// 	defer cancel()

// 	it := client.Bucket(bucket).Objects(ctx, query)
// 	go func() {
// 		for attrs, err := it.Next(); err != iterator.Done; attrs, err = it.Next() {
// 			if err == nil {
// 				ans <- attrs.Name
// 			} else {
// 				log.Fatalf("Bucket(%q).Objects: %v", bucket, err)
// 			}
// 		}
// 		close(ans)
// 	}()

// 	for i := range ans {
// 		fmt.Printf("%s\n", i)
// 	}
// }

// func BucketDownload(bucket, folder, destFolder string) {
// 	ans := make(chan string, 248)

// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)

// 	if err != nil {
// 		log.Fatalf("storage.NewClient: %v", err)
// 	}

// 	defer client.Close()

// 	query := &storage.Query{Prefix: folder}

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
// 	defer cancel()

// 	it := client.Bucket(bucket).Objects(ctx, query)
// 	go func() {
// 		for attrs, err := it.Next(); err != iterator.Done; attrs, err = it.Next() {
// 			if err == nil {
// 				ans <- attrs.Name
// 			} else {
// 				log.Fatalf("Bucket(%q).Objects: %v", bucket, err)
// 			}
// 		}
// 		close(ans)
// 	}()

// 	for i := range ans {
// 		DownloadFile(bucket, i, fmt.Sprintf("%s/%s", destFolder, filepath.Base(i)))
// 	}
// }

// func DownloadPdb(bucket, folder, destFolder string) {
// 	ans := make(chan string, 248)

// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)

// 	if err != nil {
// 		log.Fatalf("storage.NewClient: %v", err)
// 	}

// 	defer client.Close()

// 	query := &storage.Query{Prefix: folder}

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
// 	defer cancel()

// 	it := client.Bucket(bucket).Objects(ctx, query)
// 	go func() {
// 		for attrs, err := it.Next(); err != iterator.Done; attrs, err = it.Next() {
// 			if err == nil {
// 				if strings.HasSuffix(attrs.Name, ".pdb") {
// 					ans <- attrs.Name
// 				}
// 			} else {
// 				log.Fatalf("Bucket(%q).Objects: %v", bucket, err)
// 			}
// 		}
// 		close(ans)
// 	}()

// 	for i := range ans {
// 		DownloadFile(bucket, i, fmt.Sprintf("%s/%s", destFolder, filepath.Base(i)))
// 	}
// }
