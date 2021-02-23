package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	_ "github.com/go-pg/pg/v10/orm"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

var Db *pg.DB


func connect()(con *pg.DB) {
	address:=fmt.Sprintf("%s:%s","localhost","5432")
	options :=&pg.Options{
		Addr: address,
		User: "postgres",
		Password: "Kostaq",
		Database: "postgres",

	}
	Db  = pg.Connect(options)
	Db.AddQueryHook(pgdebug.DebugHook{
		Verbose: true,
	})
	err := Db.Ping(context.Background())
	if err != nil  {
		panic(err)
	}
	//createSchema(Db)
	fmt.Println(" connected ")
	return
}
type Document struct {
	Id string
	Filename string
	Content string
	OriginalFileName string
	DocsCount int
	ArrivedAt time.Time
	State string
	ProcessedAt int64
}
type InputDocument struct {
	FileName string
	Content string
	ProcessedAt int64
}

func main(){
	connect()
	mux:=http.NewServeMux()
	mux.HandleFunc("/persistFile",readUrl)
	log.Fatal(http.ListenAndServe(":8080",mux))
}

func readUrl(w http.ResponseWriter,r *http.Request)  {
	if r.Method!="POST"{
		http.Error(w,http.StatusText(405),http.StatusMethodNotAllowed)
		return
	}
	inputD:=InputDocument{}
	err:=json.NewDecoder(r.Body).Decode(&inputD)
	if err!=nil{
		panic(err)
	}
	inputJ,err:=json.Marshal(inputD)
	if err!=nil{
		panic(err)
	}
	fmt.Println(inputD)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_,err=w.Write(inputJ)
	if err!=nil{
		panic(err)
	}
	id:=uuid.NewV4().String()
	doc:=&Document{}
	doc.Id=id
	doc.Filename=inputD.FileName
	doc.Content=base64.StdEncoding.EncodeToString([]byte(inputD.Content))
	doc.OriginalFileName="test"
	doc.ArrivedAt=time.Now()
	doc.State="to_process"
	doc.ProcessedAt=inputD.ProcessedAt

	_,err=Db.Model(doc).Returning("*").Insert()
	if err!=nil{
		panic(err)
	}
}

//func createSchema(db *pg.DB) error{
//	models:=[]interface{}{
//	(*Document)(nil),
//	}
//	for _,model:=range models{
//		err:=Db.Model(model).CreateTable(&orm.CreateTableOptions{
//			IfNotExists: true,
//		})
//		if err!=nil{
//			return err
//		}
//	}
//	return nil
//}
