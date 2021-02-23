package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/extra/pgdebug"
	//"github.com/go-pg/pg/orm"
	_ "github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/v10"
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
		Database: "gopg",

	}
	Db  = pg.Connect(options)
	Db.AddQueryHook(pgdebug.DebugHook{
		Verbose: true,
	})
	err := Db.Ping(context.Background())
	if err != nil  {
		panic(err)
	}
	//err=createSchema(Db)
	//	if err!=nil{
	//		panic(err)
	//	}


	fmt.Println(" connected ")
	return
}
type dc struct{
	Id string
	RawContent string
	ReportName string
	OrderN int32
	Customer string
	StorageUrl string
	State string
	CreatedAt time.Time
}

type Document struct {
	ID string
	RawContent string
	ParsedContent ParsedContent
}
type ParsedContent struct {
	ReportName string
	OrderN int32
	Customer string
}
type UpdateDocument struct {
	ID string
	Storageurl string

}

func main()  {
	connect()
	mux:=http.NewServeMux()
	mux.HandleFunc("/persist",readUrl)
	mux.HandleFunc("/updateUrl",updateUrl)
	log.Fatal(http.ListenAndServe(":8080",mux))

}

func readUrl(w http.ResponseWriter,r *http.Request)  {
	if r.Method!="POST"{
		http.Error(w,http.StatusText(405),http.StatusMethodNotAllowed)
	}

	d:=Document{}
	err:=json.NewDecoder(r.Body).Decode(&d)
	if err!=nil{
		panic(err)
	}
	dj,err:=json.Marshal(d)
	if err!=nil{
		panic(err)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_,err=w.Write(dj)
	if err!=nil{
		panic(err)
	}
	doc:= &dc{}
		doc.Id= d.ID
		doc.RawContent=base64.StdEncoding.EncodeToString([]byte(d.RawContent))
		doc.ReportName= d.ParsedContent.ReportName
		doc.OrderN=d.ParsedContent.OrderN
		doc.Customer= d.ParsedContent.Customer
		doc.StorageUrl=    "localhost:8080"
		doc.State=        "to_process"
		doc.CreatedAt= time.Now()

	_,err=Db.Model(doc).Insert()
	if err!=nil{
		panic(err)
	}

}

func updateUrl(w http.ResponseWriter,r *http.Request)  {
	if r.Method!="POST"{
		http.Error(w,http.StatusText(405),http.StatusMethodNotAllowed)
	}

	d:=UpdateDocument{}
	err:=json.NewDecoder(r.Body).Decode(&d)
	if err!=nil{
		panic(err)
	}
	dj,err:=json.Marshal(d)
	if err!=nil{
		panic(err)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_,err=w.Write(dj)
	if err!=nil{
		panic(err)
	}
	 doc:=&dc{}
	_,err=Db.Model(doc).Set("storage_url=?",d.Storageurl).Where("id=?0",d.ID).Update()
	if err!=nil{
		panic(err)
	}
}