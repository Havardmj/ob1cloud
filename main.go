package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"

	"strings"
	//"fmt"
	"errors"
	"os"
)
/*
Used:
https://mholt.github.io/json-to-go/
to transferring json to right golang struct

*/
type AsmStruct struct {
	Project string
	Owner string
	Committer string
	Commits int
	Language []string
}
type Languages struct {
	Languages []string
}
type Repo struct {
	Project   string `json:"name"`
	//Full_name string `json:"full_name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}
type Commit []struct{
	Committer string `json:"login"`
	Commits int `json:"contributions"`
}
type er struct{
	Message string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}
func getAndDisplayErrorMessage(u string) {

	content, err := http.Get(u)
	if err != nil{
		errors.New("something went wrong")
	}
	body, err := ioutil.ReadAll(content.Body)
	if err != nil{
		errors.New("could not read content.body in get and display error message")
	}

	defer content.Body.Close()
	err = json.Unmarshal([]byte(body), &e)
	if err != nil{
		errors.New("could not unmarshal error message to struct")
	}

}

func get_Committer(u string) (Commit, error){

	content, err := http.Get(u)
	if err != nil{
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(content.Body)
	if err != nil {
		panic(err.Error())
	}
	defer content.Body.Close()
	com := Commit{}
	err = json.Unmarshal([]byte(body), &com)
	//err = json.Unmarshal(body, &com)
	if err != nil {
		//panic(err.Error())

		return Commit{}, err
	}
	return com, nil
}

func get_lang(u string) (Languages, error){
	//source: https://blog.golang.org/go-maps-in-action
	content, err := http.Get(u)
	if err != nil{
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(content.Body)
	if err != nil {
		panic(err.Error())
	}
	defer content.Body.Close()
	var a interface{}
	lang := Languages{}

	//err = json.Unmarshal(body, &a)
	err = json.Unmarshal([]byte(body), &a)
	if err != nil {
		panic(err.Error())

		return Languages{}, err
	}

	m := a.(map[string]interface{})
	//p.Fill_struct(m)
	for r := range m{
		lang.Languages = append(lang.Languages,r)
	}

	return lang, nil
}

func get_ProjectAndOwnerID(u string) (Repo, error){

	//data:=Repo{}
	content, err := http.Get(u)
	if err != nil{
		panic(err.Error())
	}
	Pro:= Repo{}
	body, err := ioutil.ReadAll(content.Body)
	if err != nil {
		panic(err.Error())
	}
	defer content.Body.Close()

	err = json.Unmarshal([]byte(body), &Pro)
	//err = json.Unmarshal(body, &Pro)
	if err != nil {
		panic(err.Error())

		return Repo{}, err
	}
	//fmt.Println(Pro)
	return Pro, nil
}
func handlerRequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	parts:= strings.Split(r.URL.Path, "/")

	//"https://api.github.com/repos/apache/kafka"
	if len(parts) != 6{
		status := 400
		http.Error(w, http.StatusText(status), status)
		return
	}

	u := "https://api.github.com/repos/" + parts[4]+ "/" +parts[5]
	u2 := u + "/contributors"
	u3 := u + "/languages"

	ProjectID_Owner, ok:= get_ProjectAndOwnerID(u)
	if ok == nil {
		CommitData, ok := get_Committer(u2)
		if ok == nil {
			Lang, ok := get_lang(u3)
			if ok == nil {
				Feed := AsmStruct{
					Project:   ProjectID_Owner.Project,
					Owner:     ProjectID_Owner.Owner.Login,
					Committer: CommitData[0].Committer,
					Commits:   CommitData[0].Commits,
					Language:  Lang.Languages,
				}

				json.NewEncoder(w).Encode(Feed)
			}else {
				getAndDisplayErrorMessage(u3)
				json.NewEncoder(w).Encode(&e)
			}
		}else {
			getAndDisplayErrorMessage(u2)
			json.NewEncoder(w).Encode(&e)
		}
	}else{
		getAndDisplayErrorMessage(u)
		json.NewEncoder(w).Encode(&e)
	}
}
var e er
func main(){
	//template for input function:
	// projectinfo/v1/github.com/owner/repo
	port := os.Getenv("PORT")
	http.HandleFunc("/projectinfo/v1/", handlerRequest)
	if port == ""{
		port = ":8080"
	}
	http.ListenAndServe(":"+port, nil)
}
