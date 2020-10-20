package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Jira struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

type Fields struct {
	Summary            string    `json:"summary"`
	AcceptanceCriteria string    `json:"acceptanceCriteria"`
	IssueType          IssueType `json:"issuetype"`
	Description        string    `json:"description"`
}

type IssueType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func encodeJql(jql string) string {
	return url.QueryEscape(jql)
}

func createQueryParams() string {
	nameFieldAcceptanceCriteria := viper.GetString("jira.acceptance_fieldname")
	if nameFieldAcceptanceCriteria != "" {
		nameFieldAcceptanceCriteria = "," + nameFieldAcceptanceCriteria
	}
	// limit fields to the few required to create a plan
	fields := "&fields=summary,issuetype,description" + nameFieldAcceptanceCriteria

	limit := viper.GetInt("jira.limit")
	if limit > 50 {
		limit = 50
	}

	// limit the result to a maximum
	maxResults := "&maxResults=" + strconv.Itoa(limit)
	// send back combination all query parameters
	return fields + maxResults
}

func simplifyResponse(bodyBytes []byte) []byte {
	nameFieldAcceptanceCriteria := viper.GetString("jira.acceptance_fieldname")
	if nameFieldAcceptanceCriteria != "" {
		bodyString := string(bodyBytes)
		return []byte(strings.ReplaceAll(bodyString, nameFieldAcceptanceCriteria, "acceptanceCriteria"))
	} else {
		return bodyBytes
	}
}

func getListOfTickets(username string, passwd string, url string, jql string) Jira {
	var target = url + encodeJql(jql) + createQueryParams()

	client := &http.Client{}
	req, err := http.NewRequest("GET", target, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	simpleResponse := simplifyResponse(bodyBytes)

	var jiraRoot Jira
	json.Unmarshal(simpleResponse, &jiraRoot)

	return jiraRoot
}
