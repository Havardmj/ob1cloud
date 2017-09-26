package main

import "testing"

func Test_commit(t *testing.T) {

	url := "https://api.github.com/repos/golang/go/contributors"
	c, ok := getCommitter(url)
	if ok != nil && len(c) < 1 {
		t.Error("Function faild. Fuction getCommitter faild to extract data")
	}
}

func Test_ProjectInfo(t *testing.T) {
	url := "https://api.github.com/repos/golang/go"
	c, ok := getProjectAndOwnerID(url)
	if ok != nil && c.Project != "" {
		t.Error("Function faild. Function getProjectAndOwnerID faild to extract data")
	}
}

func Test_Lang(t *testing.T) {
	url := "https://api.github.com/repos/golang/go/languages"
	c, ok := getLang(url)
	if ok != nil && len(c.Languages) < 1 {
		t.Error("Function failed. Function getLang faild to extract data")
	}
}
