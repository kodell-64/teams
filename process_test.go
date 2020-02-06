package main
import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"strings"
	"log"
)

// This module POSTS (creates) a member with a test entry. Then, retrieves it.
// Finally, the member is deleted and ensured to be removed.
// Many more tests could be added here.

const (
	test_url string = "http://localhost:8000/"
	test_name string = "Theodore Smith"
	test_title string = "Team members for Codelitt, Korey O'Dell"
)

func TestUpdateMember(t *testing.T) {
	data := url.Values{}
	data.Add("name", test_name)
	data.Add("employee_contractor", "employee")
	data.Add("i_role_contract_duration", "Software Engineer")
	data.Add("tags", "C#, Golang, C++")

	req, err := http.NewRequest("POST", "/update_member", strings.NewReader(data.Encode()))
	if err != nil {
		t.Errorf("TestUpdateMember failed: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	UpdateMember(w, req)

}

func TestRetrieveMember(t *testing.T) {
	// query db for our test_name, get id and go from there
	db := dbConn()
	seldb, err := db.Query("select id from members where name=?", test_name)
	seldb.Next()
	var id string
	err = seldb.Scan(&id)
	if err != nil {
		t.Errorf("select failed")
	}
	
	data := url.Values{}
	data.Add("id", id)
	req, err := http.NewRequest("POST", "/retrieve_member", strings.NewReader(data.Encode()))
	if err != nil {
		t.Errorf("TestRetrieveMember POST failed: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	RetrieveMember(w, req)

	got := w.Body.String()
	if !strings.Contains(got, test_name) {
		t.Errorf("did not get: ", test_name)
	}
}


func TestDeleteMember(t *testing.T) {
	// query db for our test_name, delete the member via API, then ensure it is removed
	db := dbConn()
	seldb, err := db.Query("select id from members where name=?", test_name)
	seldb.Next()
	var id string
	err = seldb.Scan(&id)
	if err != nil {
		t.Errorf("select failed")
	}
	
	data := url.Values{}
	data.Add("id", id)
	req, err := http.NewRequest("POST", "/delete_member", strings.NewReader(data.Encode()))
	if err != nil {
		t.Errorf("TestDeleteMember POST failed: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	DeleteMember(w, req)

	seldb, _ = db.Query("select id from members where id=?", id)
	err = seldb.Scan(&id)
	if err == nil {
		t.Errorf("TestDeleteMember failed")
	}
}

func TestMainPage(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", test_url, nil)
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// very simple test - look for known <TITLE> string
	if !strings.Contains(string(f), test_title) {
		t.Errorf("did not get: ", test_title)
	}

	resp.Body.Close()
}
