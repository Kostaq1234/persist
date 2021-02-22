

package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

var db *sql.DB
func init() {
	var err error
	db,err=sql.Open("postgres","postgres://postgres:Kostaq@localhost/postgres?sslmode=disable")
	if err!=nil{
		panic(err)
	}
	if err=db.Ping(); err!=nil{
		panic(err)
	}
	fmt.Println("You are connected")


}


type Document struct {
	id string
	Inputdocument InputDocument
	originalFileName string
	docsCount int
	arrivedAt time.Time
	state string
}
type InputDocument struct {
	FileName string
	Content string
	ProcessedAt int64
}

func main(){
	mux:=http.NewServeMux()
	mux.HandleFunc("/persist",readUrl)
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
w.Header().Set("Content-Type","application/json")
w.WriteHeader(http.StatusOK)
_,err=w.Write(inputJ)
if err!=nil{
	panic(err)
}

doc:=Document{}
id:=uuid.Must(uuid.NewV4()).String()
stateC:="to_process"
doc.id=id
content:=base64.StdEncoding.EncodeToString([]byte(inputD.Content))
processedAt:=inputD.ProcessedAt
doc.arrivedAt=time.Now()
doc.state=stateC

_,err=db.Exec("INSERT INTO documents(id,\"fileName\",content,\"processedAt\",\"originalFileName\",\"arrivedAt\",\"state \")\n values ($1,$2,$3,$4,$5,$6,$7)",doc.id,inputD.FileName,content,processedAt,"test",doc.arrivedAt,doc.state)
if err!=nil{
	http.Error(w,http.StatusText(500),http.StatusInternalServerError)
	panic(err)
	return
}

}


