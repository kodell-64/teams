// Golang CRUD for Codelitt, author Korey O'Dell
// February 3, 2020
//
// uses gorilla-mux for routing
//

package main

import (
	"os"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"database/sql"
	"github.com/ziutek/mymysql/godrv"
)

// simple delimiter for member records
var delimiter string = "|||"
var db_driver string = "mymysql"
var db_connection_params string = "team_members/tm/tm"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", MainPage)
	r.HandleFunc("/update_member", UpdateMember)
	r.HandleFunc("/search_member", SearchMember)
	r.HandleFunc("/delete_member", DeleteMember)
	r.HandleFunc("/retrieve_member", RetrieveMember)

	log.Print("starting up")
	http.ListenAndServe(":8000", r)
}

// MainPage, starts HTML off with index.html
func MainPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("index.html")
	if err != nil { 
		log.Print("template parsing error: ", err) // log it
	}
	t.Execute(w, nil)
}

// UpdateMember handles both CRUD tasks of create and update
// if id then update existing member
// if !id then create new member
func UpdateMember(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	db := dbConn()
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		employee_contractor := r.FormValue("employee_contractor")
		role_contract_duration := r.FormValue("i_role_contract_duration")
		tags := r.FormValue("tags")
		
		if len(id) == 0 {
			member, err := db.Prepare("insert members values (?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Print("create(Update) member failed",err);
				return;
			}
			member.Exec(nil, name, employee_contractor, role_contract_duration, tags, now.Unix())
			fmt.Fprintf(w, "Member '%s' added.\n", name)
		} else {
			member, err := db.Prepare("update members set name=?, type=?, role_contract_duration=?, tags=? where id=?")
			if err != nil {
				log.Print("create(Update) member failed",err);
				return;
			}
			member.Exec(name, employee_contractor, role_contract_duration, tags, id)
			fmt.Fprintf(w, "Member '%s' updated.\n", name)
		}
	}
}

// RetrieveMember
// simple retrieve of a member specified by id
func RetrieveMember(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		id := r.FormValue("id")
		seldb, err := db.Query("select id, name, type, role_contract_duration, tags, created from members where id=?", id)
		if err != nil {
			log.Print("retrieve member failed")
			return;
		}
		for seldb.Next() {
			var m_id int
			var name string
			var m_type string
			var role_contract_duration string
			var tags string
			var created int
			seldb.Scan(&m_id, &name, &m_type, &role_contract_duration, &tags, &created)
			// use simple delimited protocol, would use JSON in reality
			fmt.Fprintf(w, "%d%s%s%s%s%s%s%s%s%s%d\n", m_id, delimiter, name, delimiter, m_type,
				delimiter, role_contract_duration, delimiter, tags, delimiter, created)
		}
	}
}

// SearchMember allows for querying of any of the fields
func SearchMember(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		search_params := r.FormValue("search_params")
		q_str := fmt.Sprintf("select id, name, type, role_contract_duration, tags, created from members where name like '%%%s%%' or type like '%%%s%%' or role_contract_duration like '%%%s%%' or tags like '%%%s%%'", search_params, search_params, search_params, search_params)
		rows, err := db.Query(q_str)
		if err != nil {
			log.Print("search member failed")
			return;
		}

		for rows.Next() {
			var id int
			var name string
			var m_type string
			var role_contract_duration string
			var tags string
			var created int
			rows.Scan(&id, &name, &m_type, &role_contract_duration, &tags, &created)
			// use simple delimited protocol, would use JSON in reality
			fmt.Fprintf(w, "%d%s%s%s%s%s%s%s%s%s%d\n", id, delimiter, name, delimiter, m_type,
				delimiter, role_contract_duration, delimiter, tags, delimiter, created)
		}
	}
}

// DeleteMember instantiates deletion of a member
func DeleteMember(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		id := r.FormValue("id")
		delete, err := db.Prepare("delete from members where id=?")
		if err != nil {
			log.Print("delete member failed");
			return;
		}
		delete.Exec(id)
		fmt.Fprintf(w, "Member deleted.\n")
	}
}

func dbConn() (db *sql.DB) {
	godrv.Register("SET NAMES utf8") 
	db, err := sql.Open(db_driver, db_connection_params)
	if err != nil {
		log.Print("can not connect to database, exiting. ", err)
		os.Exit(1)
	}
	return db
}
