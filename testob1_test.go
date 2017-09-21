package main

import "testing"

func Test_commit(t *testing.T){

	url:="https://api.github.com/repos/golang/go/contributors"
	c, ok :=get_Committer(url)
	if ok != nil && len(c) < 1{
		t.Error("Function faild.\n Fuction get_Committer faild to extract data")
	}

}

func Test_ProjectInfo(t *testing.T){
	url:="https://api.github.com/repos/golang/go"
	c, ok:=get_ProjectAndOwnerID(url)
	if ok != nil && c.Project != "" {
		t.Error("Function faild.\n Function get_ProjectAndOwnerID faild to extract data")
	}
}

func Test_Lang(t *testing.T){
	url:="https://api.github.com/repos/golang/go/languages"
	c, ok:=get_lang(url)
	if ok != nil && len(c.Languages) < 1 {
		t.Error("Function failed.\n Function get_lang faild to extract data")
	}
}
