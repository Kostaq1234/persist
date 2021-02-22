package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"

	//"log"
	//"net/http"
	//"time"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:Kostaq@localhost/document?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You are connected")
}
type Document struct {
	ID string
	RawContent string
	ParsedContent ParsedContent
	storageUrl string
	state string
	createdAt time.Time
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

	stateC:="to_process"
	storageUrl:="localhost:8080"
	createdAt:=time.Now()
	rawContent:=base64.StdEncoding.EncodeToString([]byte(d.RawContent))

	_,err=db.Exec("INSERT INTO document (id, \"rawContent\", \"reportName\", \"orderN\", customer, \"createdAt\", state, \"storageUrl\") VALUES  ($1,$2,$3,$4,$5,$6,$7,$8)",d.ID,rawContent,d.ParsedContent.ReportName,d.ParsedContent.OrderN,d.ParsedContent.Customer,createdAt,stateC,storageUrl)
	if err!=nil{
		http.Error(w,http.StatusText(500),http.StatusInternalServerError)
		panic(err)
		return
	}

}
func updateUrl(w http.ResponseWriter,r *http.Request){
	if r.Method!="POST"{
		http.Error(w,http.StatusText(405),http.StatusMethodNotAllowed)
	}
	u:=UpdateDocument{}
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		panic(err)
	}
	sqlStatement:=`UPDATE document SET "storageUrl"=$1 WHERE id=$2;`
	_,err=db.Exec(sqlStatement,u.Storageurl,u.ID)
	if err!=nil{
		http.Error(w,http.StatusText(500),http.StatusInternalServerError)
		return
	}
}