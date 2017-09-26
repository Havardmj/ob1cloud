package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"strings"
	//"fmt"
	//"errors"
	"os"
)

/*
AsmStruct the other struct's to this final struct
Languages insert array of languages

Used:
https://mholt.github.io/json-to-go/
to transferring json to right golang struct

*/
//Assemble the other struct's to this final struct
type AsmStruct struct { //AsmStruct final struct for assembly
	Project   string   //project name
	Owner     string   //project owner
	Committer string   //top committer
	Commits   int      //number of commits to top committer
	Language  []string //array of all languages
}

/*
Languages struct used to insert array of code-languages
is exported to get_lang() function where it unmarshal json into a empty interface, map it across a string interface, to be
inserted into the struct
*/
type Languages struct {
	Languages []string //array of strings
}

/*
Repo struct used to insert project name and project owner
is exported to function: get_ProjectAndOwnerID(u string)
where json is unmarshal and inserted.
*/
type Repo struct {
	Project string `json:"name"` //project name
	Owner   struct {
		Login string `json:"login"` //project owner
	} `json:"owner"`
}

/*
Commit struct used to insert Top committer and number of
commits user have produced. Is exported to function:
get_committer() where we unmarshal json and insert to struct
*/
type Commit []struct {
	Committer string `json:"login"`         //committer
	Commits   int    `json:"contributions"` //number of commits
}

/*
er struct used to pickup error messages the may occur when you
type inn wrong url request or get temporary block from the site
for using to many request. Help full for the user to get the
feedback.
*/
type er struct {
	Message          string `json:"message"`           //display error message
	DocumentationURL string `json:"documentation_url"` //doc url
}

/*
This function is only called when one of the other get functions
return a error-message. It's called to pickup error messages
return by the github.com-api. To display use full info.
*/
func getAndDisplayErrorMessage(u string) {

	content, err := http.Get(u) //retrieve content from url
	if err != nil {             //check if error is 0
		panic(err.Error()) //print error to console
	}
	//read all body content
	body, err := ioutil.ReadAll(content.Body)
	if err != nil { //check if readAll retrieve error
		panic(err.Error()) //print error to console
	}

	defer content.Body.Close()             //close content.body after done
	err = json.Unmarshal([]byte(body), &e) //unmarshal json
	if err != nil {                        //check for error
		panic(err.Error()) //if so print error to console
	}
}

/*
This function is called when is time to extract Committer
It read the content of a address. Read the content of the body,
unmarshal the content and return the data with error set/nil
*/
func getCommitter(u string) (Commit, error) {

	content, err := http.Get(u) //get content of url
	if err != nil {             //check if error is 0
		panic(err.Error()) //print error to console
	}
	//readAll content of the body into body er=set/nil
	body, err := ioutil.ReadAll(content.Body)
	if err != nil { //check if error is set
		panic(err.Error()) //output error to console
	}
	defer content.Body.Close()               //close content.body after done
	com := Commit{}                          //init a new Commit struct-obj
	err = json.Unmarshal([]byte(body), &com) //unmarshal into com
	if err != nil {                          //check if error occur
		return Commit{}, err //if so return empty struct
	}
	return com, nil //successfull extract return struct / nil
}

/*
This function is called in the handlerRequest to extract Languages
from the repository. It handles the url string and return the unmarshal
data along a error set/nil
*/
func getLang(u string) (Languages, error) {
	//source: https://blog.golang.org/go-maps-in-action
	content, err := http.Get(u) //get the content of url input
	if err != nil {             //check is err is nil
		panic(err.Error()) //if so output to console
	}
	body, err := ioutil.ReadAll(content.Body) //read all content
	if err != nil {                           //check for error
		panic(err.Error()) //output if true
	}
	defer content.Body.Close() //close body after done
	var a interface{}          //construct empty interface
	lang := Languages{}        //construct member of Languages struct

	//err = json.Unmarshal(body, &a)
	err = json.Unmarshal([]byte(body), &a) //unmarshal into a interface
	if err != nil {                        //check for error
		return Languages{}, err //if error set return empty string
	}
	//a set to map it along string interface
	m := a.(map[string]interface{}) //this is read into variable m
	//p.Fill_struct(m)
	for r := range m { //loop through m
		lang.Languages = append(lang.Languages, r) //append it to Languages
	}
	return lang, nil //return lang with error set to nil
}

/*
This function is called to extract Project name and owner. It read content of
the url, it's body, then unmarshal and read it to the Repo struct.
It do so by either returning data, or empty string, with/without error.
*/
func getProjectAndOwnerID(u string) (Repo, error) {

	//data:=Repo{}
	content, err := http.Get(u) //get url content
	if err != nil {             //check if error is set
		panic(err.Error()) //if error output to console
	}
	Pro := Repo{}                             //initialize empty Repo member
	body, err := ioutil.ReadAll(content.Body) //read content body into body/err
	if err != nil {                           //is error is not equal to nil
		panic(err.Error()) //output to console
	}
	defer content.Body.Close() //close body after done

	err = json.Unmarshal([]byte(body), &Pro) //unmarshal into Pro
	//err = json.Unmarshal(body, &Pro)
	if err != nil { //again check for error is set or not
		return Repo{}, err //if error return empty Repo alongside error
	}
	//fmt.Println(Pro)
	return Pro, nil //if success full return data and error equal nil
}

/*
This function takes two parameters, a reponsewriter that handles the output,
and a Request that handles the request through the socket. This function also
act's like a main function calling and assembling the other function and data
*/
func handlerRequest(w http.ResponseWriter, r *http.Request) {

	//set the content-Type to json format
	w.Header().Set("Content-Type", "application/json")

	/*r.URL.Path gets the url path inserted. With applying strings.Split
	we then can split the string inn several part into parts on every '/'
	*/
	parts := strings.Split(r.URL.Path, "/")

	//"https://api.github.com/repos/apache/kafka"
	if len(parts) != 6 { //allowed length is 6, if requirement not met:
		status := 400                                  //status equal 400
		http.Error(w, http.StatusText(status), status) //output status in browser
		return                                         //then return - rest of function is not executed
	}
	//assemble valid url string:
	u := "https://api.github.com/repos/" + parts[4] + "/" + parts[5]
	u2 := u + "/contributors" //add '/contributors'to get committer function
	u3 := u + "/languages"    //add '/languages to get languages function

	ProjectIDOwner, ok := getProjectAndOwnerID(u) //call and get data from project function
	if ok == nil {                                //check for error
		CommitData, ok := getCommitter(u2) //call and get data pro committer function
		if ok == nil {                     //check for error
			Lang, ok := getLang(u3) //call and get data from language function
			if ok == nil {          //check for error
				Feed := AsmStruct{ //set Feed equal to AsmStruct and insert struct value here.
					Project:   ProjectIDOwner.Project,     //project name insert to project
					Owner:     ProjectIDOwner.Owner.Login, //project owner insert
					Committer: CommitData[0].Committer,    //input committer from committer struct
					Commits:   CommitData[0].Commits,      //input commits from committer
					Language:  Lang.Languages,             //input languages from language struct
				}

				json.NewEncoder(w).Encode(Feed) //if all ok, json encode data for output
				/*
					will only run the else-statement if and where the if-statement fails to comply
				*/
			} else {
				getAndDisplayErrorMessage(u3) //if error in first if-statement get error message
				json.NewEncoder(w).Encode(&e) //encode and output message
			}
		} else {
			getAndDisplayErrorMessage(u2) //if error in second if-statement get error message
			json.NewEncoder(w).Encode(&e) //encode and output message
		}
	} else {
		getAndDisplayErrorMessage(u)  //if error in third if-statement get error message
		json.NewEncoder(w).Encode(&e) //encode and output message
	}
}

var e er //AsmStruct variables to access struct data

func main() { //the main function
	//template for input function:
	// projectinfo/v1/github.com/owner/repo

	port := os.Getenv("PORT")                           //use to export port used in heroku
	http.HandleFunc("/projectinfo/v1/", handlerRequest) //call HandleFunc
	if port == "" {
		port = "8080"
	}
	//this line enable it to run in heroku-app cloud
	http.ListenAndServe(":"+port, nil)

	//enable this line, and comment-out the line over if you want to run it locally instead
	//http.ListenAndServe(":8080", nil)

}
