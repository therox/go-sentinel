package backend

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/therox/go-sentinel/tools"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GoogleCloudClient struct {
	indexFile     string
	indexURL      string
	tmpDir        string
	dsList        []Dataset
	storageClient *storage.Client
	infoWriter    io.Writer
	clientContext context.Context
}

type Dataset struct {
	GranuleID            string      // L1C_T56MKT_A029042_20210113T002708
	ProductID            string      // S2A_MSIL1C_20210113T002711_N0209_R016_T56MKT_20210113T021104
	DatatakeIdentifier   string      // GS2A_20210113T002711_029042_N02.09
	MGRSTile             string      // 56MKT
	SensingTime          time.Time   // 2021-01-13T00:28:41.451000Z
	TotalSize            int64       // 248420128
	CloudCover           float64     // 86.8793
	GeometricQualityFlag interface{} // ???
	GenerationTime       time.Time   // 2021-01-13T02:11:04.000000Z
	NorthLat             float64     // -6.325329606275461
	SouthLat             float64     // -7.319001758587432
	WestLon              float64     // 150.28270701458754
	EastLon              float64     // 150.7686688433502
	BaseURL              string      // gs://gcp-public-data-sentinel-2/tiles/56/M/KT/S2A_MSIL1C_20210113T002711_N0209_R016_T56MKT_20210113T021104.SAFE
}

func (gc *GoogleCloudClient) SearchDataset(datasetName string) {
	fmt.Println("Searching on GoogleCloud")
}
func (gc *GoogleCloudClient) Download(datasetName string) {
	fmt.Printf("Downloading dataset %s from GoogleCloud\n", datasetName)
}

func NewGCClient(indexFile string, updateIndex bool) *GoogleCloudClient {
	if indexFile == "" {
		indexFile = "index.csv.gz"
	}

	c := &GoogleCloudClient{
		indexURL:      "https://storage.googleapis.com/gcp-public-data-sentinel-2/index.csv.gz",
		indexFile:     indexFile,
		infoWriter:    os.Stdout,
		clientContext: context.Background(),
	}

	sClient, err := storage.NewClient(c.clientContext)
	if err != nil {
		log.Fatal(err)
	}
	c.storageClient = sClient

	// Creating temp dir
	tdir, err := ioutil.TempDir("", "sen_")
	if err != nil {
		log.Fatal(err)
	}
	c.tmpDir = tdir

	if updateIndex {
		// Download index file from google cloud
		// index file is accessed under https://storage.googleapis.com/gcp-public-data-sentinel-2/index.csv.gz link
		if err = tools.DownloadFile(c.indexURL, c.indexFile); err != nil {
			log.Fatalf("error on file download: %s", err)
		}
	}

	c.dl()
	// Read into structure
	err = c.readIndex()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", c.dsList[1000000])
	return c
}

func (gc *GoogleCloudClient) readIndex() error {
	t := time.Now()
	f, err := os.Open(gc.indexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer gr.Close()

	cr := csv.NewReader(gr)
	cr.Read()
	dsList := []Dataset{}
	isFirst := true
	for {
		row, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				gc.dsList = dsList
				fmt.Printf("Processed %d in %s \n", len(dsList), time.Since(t))
				err = nil
			}
			return err
		}

		st, _ := time.Parse(time.RFC3339, row[4])
		ts, _ := strconv.ParseInt(row[5], 10, 64)
		cc, _ := strconv.ParseFloat(row[6], 64)
		gt, _ := time.Parse(time.RFC3339, row[8])
		nl, _ := strconv.ParseFloat(row[9], 64)
		sl, _ := strconv.ParseFloat(row[10], 64)
		wl, _ := strconv.ParseFloat(row[11], 64)
		el, _ := strconv.ParseFloat(row[12], 64)
		dsList = append(dsList, Dataset{
			GranuleID:            row[0],
			ProductID:            row[1],
			DatatakeIdentifier:   row[2],
			MGRSTile:             row[3],
			SensingTime:          st, // 2021-01-13T00:28:41.451000Z
			TotalSize:            ts,
			CloudCover:           cc,
			GeometricQualityFlag: row[7],
			GenerationTime:       gt,
			NorthLat:             nl,
			SouthLat:             sl,
			WestLon:              wl,
			EastLon:              el,
			BaseURL:              row[13],
		})
		if len(dsList)%1000000 == 0 {
			fmt.Printf("Processed %d records\n", len(dsList))
		}
		if isFirst {
			isFirst = false
			fmt.Printf("%+v\n", dsList[0])
		}
	}
}

func (gc *GoogleCloudClient) dl() {
	log.Println("DOWNLOAD BEGIN")
	files, err := gc.listFilesWithPrefix("gcp-public-data-sentinel-2", "L2/tiles/16/M/BD/S2A_MSIL2A_20181220T162311_N0211_R097_T16MBD_20181220T195456", "")
	// err := gc.listFilesWithPrefix("gcp-public-data-sentinel-2", "L2/tiles/16/M/BD/S2A_MSIL2A_20181220T162311_N0211_R097_T16MBD_20181220T195456", "")
	// err := gc.downloadFile("gcp-public-data-sentinel-2", "L2/tiles/16/M/BD/S2A_MSIL2A_20181220T162311_N0211_R097_T16MBD_20181220T195456.SAFE/GRANULE/L2A_T16MBD_A018255_20181220T162306/QI_DATA/MSK_TECQUA_B03.gml", "/tmp/MSK_TECQUA_B03.gml")
	// err := gc.listFiles("gcp-public-data-sentinel-2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("FILES:")
		for i := range files {
			fmt.Println(files[i])
		}
	}

	log.Println("DOWNLOAD END")

}

func (gc *GoogleCloudClient) Close() error {
	gc.storageClient.Close()
	return os.RemoveAll(gc.tmpDir)
}

// listFilesWithPrefix lists objects using prefix and delimeter.
func (gc *GoogleCloudClient) listFilesWithPrefix(bucket, prefix, delim string) ([]string, error) {
	// bucket := "bucket-name"
	// prefix := "/foo"
	// delim := "_"

	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used to filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimiter, the entire tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	res := make([]string, 0)
	ctx, cancel := context.WithTimeout(gc.clientContext, time.Second*10)
	defer cancel()

	it := gc.storageClient.Bucket(bucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return res, fmt.Errorf("Bucket(%q).Objects(): %v", bucket, err)
		}
		res = append(res, attrs.Name)
	}
	return res, nil
}

func (gc *GoogleCloudClient) listFiles(bucket string) ([]string, error) {
	res := make([]string, 0)
	// bucket := "bucket-name"
	ctx := context.Background()
	// client, err := storage.NewClient(ctx)
	// if err != nil {
	// 		return fmt.Errorf("storage.NewClient: %v", err)
	// }
	// defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := gc.storageClient.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return res, fmt.Errorf("Bucket(%q).Objects: %v", bucket, err)
		}
		res = append(res, attrs.Name)
		// fmt.Fprintln(gc.infoWriter, attrs.Name)
	}
	return res, nil
}

// downloadFile downloads an object to a file.
func (gc *GoogleCloudClient) downloadFile(bucket, object string, destFileName string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	// destFileName := "file.txt"
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	f, err := os.Create(destFileName)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	rc, err := gc.storageClient.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %v", err)
	}

	fmt.Fprintf(gc.infoWriter, "Blob %v downloaded to local file %v\n", object, destFileName)

	return nil

}

// https://storage.googleapis.com/gcp-public-data-sentinel-2/L2/tiles/16/M/BD/S2A_MSIL2A_20181220T162311_N0211_R097_T16MBD_20181220T195456.SAFE/GRANULE/L2A_T16MBD_A018255_20181220T162306/IMG_DATA/R10m/T16MBD_20181220T162311_AOT_10m.jp2
// gs://gcp-public-data-sentinel-2/tiles/13/V/EG/S2B_MSIL1C_20200503T182919_N0209_R027_T13VEG_20200503T220524.SAFE

// https://storage.googleapis.com/gcp-public-data-sentinel-2/L2/tiles/13/V/EG/S2B_MSIL1C_20200503T182919_N0209_R027_T13VEG_20200503T220524.SAFE

// therox.sentinel:Z%M8L2F!jEWY9R

// {
// 	"type": "service_account",
// 	"project_id": "my-project-sentinel-341413",
// 	"private_key_id": "17b2465ee2c50e8af89d67bbd1766bd0b961724c",
// 	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCtrEemoB4H6SQB\n+sJA1O2AqnqiJYx36A9IGuhwVUdkkpMU+DTo5Y7rKnRG7pWO/n+1nd5eZCHAwj2o\nxaLL6rgFX6KK/tPWPci2D0+msw9X6RcSxeOEubJVm6/4EYzRuZReCwNOxCQ3pJWi\nJ+AH/gli/5jSPVNZftoGumqpd3lyENgifQue1RrPdCwBTpX3qEPw0fmia61kZd6I\nFiDIr6E7YaK+UMkSLLiXiiaEgFOeSj2fcFmExFqi2Db3M+PRgAPczm6+4WyeHJ8T\nCVf925fC2KYS639ikNAt6rZNfWQYIUkYWsnXyKKAOfJgSICSDgYhMOUlElgJZoEH\naHkTbK4DAgMBAAECggEAETnS9YemKDoLVJhEUjKMcsRLW5qGfzzTid8ECmL7Aie4\n/Ph4JFg4lRdl5xndAwzif+duWfBr70w/hLuYSnTigcNFl6xvjdgc+KngyEkJuY9+\nDAz4DynCstO5YCKQKXEocoInIYCPkxFOYiy6/4VGoGd3sfVb4gw/SzGMSM4ii+9O\ncNRmmOlfZrYqxsflFnCKL79D4OyOaCE4r/mmIIhz0GzCovDIVtUEwK7xigzAVmAg\nBbURkzq5kfw5X+VcUBePSlTjY+T+DU8+/cSmq3XF1bcY6fMzD9V2Hzdh0JWUYRNg\n1g88nDSWAuSGhGwezScfXq+T6uzp+YrEF+ppQPZX8QKBgQDm0rSpiecuncGAaxJK\nbQATFJza4qxmOdbnouUZD236KqR6YSsy/JwbebJhrzmCkFmVjue5txjno4oJM+F4\nEGFLGqCUo4YVgnD3/927PaDIDmFGvPrei2ZjTcXL8wt97VV/h9x3eSsT7NftdK10\nQGuVAMVn/6W0YkhDeVNyBYoMCwKBgQDAncQ7LC3u690DekzJ13EEdvabwUrPhWpl\n+QOUIOo1pfnvi+q2HhIRSzk7Zrfeac0rlTJq8LH+nriaRTq0UZSUKcmluTNG8wdx\nDfWuJ1Q9qw54xmYomqdIcjvZF4gPTy6WsO2/hU7FJCHV8F9Bq+pi7iFfkCvEVeUp\nxW8njnIo6QKBgAh+wL9sn+GLVM/RIP80O8h3/b3gkc/BWTN1P1FIzBdk8mC+qsRZ\nn7UzKWmG9WrkvSdP9x9Sqcv8rNvh2O5+S9ZjyH0v94NbJaU6+YImpoA2Z0YfNmxy\nXPZSqnDzJP0hewASk/6hqwA0u021+ZxJsqvtN1a4hfS5mUdCc9vCcG45AoGAEvi3\nFjjwwUKb/zaxj4PXzt/ExGySaMQ0SZzwPEvNpcEWUADQ+mUV1Hw7Q92HKLH1EmDc\nmWnhrjmF5zQAZzhN0RYYlWkbYtaXGgyM8T1Qb86EW2L3NvKcnWca+/P0Y8Whyhfd\nIusHaiV367GEQFyldl8MDHSPywPG0YlvvE6scdkCgYAnqkPwjf/MpTUC8qhtmpuP\njzuvIC21N/nGRWKELgh0E6XpF9yzhkbZXizGa/N731dyqgVTekPEcz4EGG+55r7K\n55wJ0KMYIgkG1LK3cBV3OImG55u0/RRpurd0nJlhLlxD+UfstJ6xz9WFcyV87EGh\nHQqaGipgYIeojmp01a4Z0w==\n-----END PRIVATE KEY-----\n",
// 	"client_email": "sentinel@my-project-sentinel-341413.iam.gserviceaccount.com",
// 	"client_id": "111722094265434604179",
// 	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
// 	"token_uri": "https://oauth2.googleapis.com/token",
// 	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
// 	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/sentinel%40my-project-sentinel-341413.iam.gserviceaccount.com"
//   }
// export GOOGLE_APPLICATION_CREDENTIALS=my-project-sentinel-341413-17b2465ee2c5.json
