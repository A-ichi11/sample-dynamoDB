package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Sample struct {
	UserID string `dynamo:"UserID,hash"`
	Name   string `dynamo:"Name,range"`
	Age    int    `dynamo:"Age"`
	Text   string `dynamo:"Text"`
}

var DYNAMO_ENDPOINT = "http://localhost:8000"

func main() {
	// クライアントの設定
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Endpoint:    aws.String(DYNAMO_ENDPOINT),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	})
	if err != nil {
		panic(err)
	}

	db := dynamo.New(sess)

	db.Table("Samples").DeleteTable().Run()

	// テーブル作成
	err = db.CreateTable("Samples", Sample{}).Run()
	if err != nil {
		panic(err)
	}

	// DBにPutします
	table := db.Table("Samples")
	err = table.Put(&Sample{UserID: "1", Name: "Test1", Age: 20}).Run()
	if err != nil {
		panic(err)
	}

	// DBからGetします
	var sampleDb Sample
	err = table.Get("UserID", "1").Range("Name", dynamo.Equal, "Test1").One(&sampleDb)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetDB%+v\n", sampleDb)

	// DBのデータをUpdateします
	text := "新しいtextです"
	err = table.Update("UserID", "1").Range("Name", "Test1").Set("Text", text).Value(&sampleDb)
	if err != nil {
		panic(err)
	}
	fmt.Printf("UpdateDB%+v\n", sampleDb)

	// DBのデータをDeleteします
	err = table.Delete("UserID", "1").Range("Name", "Test1").Run()
	if err != nil {
		panic(err)
	}

	// Delete出来ているか確認
	err = table.Get("UserID", "1").Range("Name", dynamo.Equal, "Test1").One(&sampleDb)
	if err != nil {
		// Delete出来ていれば、dynamo: no item found のエラーとなる
		fmt.Println("getError:", err)
	}

}
