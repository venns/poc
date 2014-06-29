package main

import (
	codes "alexhicks.net/manage/lib"
	structs "alexhicks.net/manage/lib"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var config structs.Config
var db *sql.DB

func newPerson(out http.ResponseWriter, req *http.Request) {
	out.Header().Add("Access-Control-Allow-Origin", "*")
	req.ParseForm()
	err := db.Ping()
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlConnectionFailed, err.Error()))
		return
	}
	firstname := req.PostForm.Get("firstname")
	lastname := req.PostForm.Get("lastname")
	if (firstname == "") || (lastname == "") {
		io.WriteString(out, structs.NewStatus(codes.ErrorInvalidArguments, "Invalid arguments"))
		return
	}
	// We really don't need the result here
	_, err = db.Exec("INSERT INTO person (firstname, lastname) VALUES(?, ?)", firstname, lastname)
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlFailedExecute, err.Error()))
		return
	}
	io.WriteString(out, structs.NewStatus(codes.Success, ""))
}

func newTime(out http.ResponseWriter, req *http.Request) {
	out.Header().Add("Access-Control-Allow-Origin", "*")
	req.ParseForm()
	err := db.Ping()
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlConnectionFailed, err.Error()))
		return
	}
	personID := req.PostForm.Get("personid")
	taskID := req.PostForm.Get("taskid")
	startTime := req.PostForm.Get("starttime")
	endTime := req.PostForm.Get("endtime")
	if (personID == "") || (taskID == "") || (startTime == "") || (endTime == "") {
		io.WriteString(out, structs.NewStatus(codes.ErrorInvalidArguments, "Invalid arguments"))
		return
	}
	_, err = db.Exec("INSERT INTO actualtime (personid, taskid, datenow, starttime, endtime) VALUES(?, ?, ?, ?, ?)", personID, taskID, time.Now(), startTime, endTime)
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlFailedExecute, err.Error()))
		return
	}
	io.WriteString(out, structs.NewStatus(codes.Success, ""))
}

func newTask(out http.ResponseWriter, req *http.Request) {
	out.Header().Add("Access-Control-Allow-Origin", "*")
	req.ParseForm()
	err := db.Ping()
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlConnectionFailed, err.Error()))
		return
	}
	name := req.PostForm.Get("name")
	if name == "" {
		io.WriteString(out, structs.NewStatus(codes.ErrorInvalidArguments, "Invalid arguments"))
		return
	}
	_, err = db.Exec("INSERT INTO task (name) VALUES(?)", name)
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlFailedExecute, err.Error()))
		return
	}
	io.WriteString(out, structs.NewStatus(codes.Success, ""))
}

func getTasks(out http.ResponseWriter, req *http.Request) {
	out.Header().Add("Access-Control-Allow-Origin", "*")
	err := db.Ping()
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlConnectionFailed, err.Error()))
		return
	}
	rows, err := db.Query("SELECT id, name FROM task")
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlFailedExecute, err.Error()))
		return
	}
	var tasks []structs.Task
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err = rows.Scan(&id, &name); err != nil {
			io.WriteString(out, structs.NewStatus(codes.ErrorSqlFailedExecute, err.Error()))
			return
		}
		tasks = append(tasks, structs.Task{id, name})
	}
	io.WriteString(out, structs.NewData(codes.Success, tasks))
}

func getPersons(out http.ResponseWriter, req *http.Request) {
	out.Header().Add("Access-Control-Allow-Origin", "*")
	err := db.Ping()
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlConnectionFailed, err.Error()))
		return
	}
	rows, err := db.Query("SELECT id, firstname, lastname FROM person")
	if err != nil {
		io.WriteString(out, structs.NewStatus(codes.ErrorSqlFailedExecute, err.Error()))
		return
	}
	var persons []structs.Person
	defer rows.Close()
	for rows.Next() {
		var id int
		var firstname string
		var lastname string
		if err = rows.Scan(&id, &firstname, &lastname); err != nil {
			io.WriteString(out, structs.NewStatus(codes.ErrorSqlFailedExecute, err.Error()))
			return
		}
		persons = append(persons, structs.Person{id, firstname, lastname})
	}
	io.WriteString(out, structs.NewData(codes.Success, persons))
}

func main() {
	rawJson, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(rawJson, &config)
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", config.MysqlUsername, config.MysqlPassword, config.MysqlHostname, config.MysqlDatabase))
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/newperson", newPerson)
	http.HandleFunc("/newtime", newTime)
	http.HandleFunc("/newtask", newTask)
	http.HandleFunc("/getpersons", getPersons)
	http.HandleFunc("/gettasks", getTasks)
	serverAddress := fmt.Sprintf("%s:%d", config.HttpAddress, config.HttpPort)
	fmt.Printf("Listening on %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
