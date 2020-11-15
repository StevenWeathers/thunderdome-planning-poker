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

type JiraRequest struct {
	Client     *http.Client
	BaseUrl    string
	EncodedJql string
	Fields     string
	MaxResults int
	Username   string
	Password   string
}

type Jira struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
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

func createFieldsQueryParam() string {
	nameFieldAcceptanceCriteria := viper.GetString("jira.acceptance_fieldname")
	if nameFieldAcceptanceCriteria != "" {
		nameFieldAcceptanceCriteria = "," + nameFieldAcceptanceCriteria
	}
	// limit fields to the few required to create a plan
	return "summary,issuetype,description" + nameFieldAcceptanceCriteria
}

func createMaxResultQueryParam() int {
	limit := viper.GetInt("jira.limit")
	if limit > 50 {
		limit = 50
	}

	// limit the result to a maximum
	return limit
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

	// create request obj
	jiraReq := JiraRequest{
		Client:     &http.Client{},
		BaseUrl:    url,
		EncodedJql: encodeJql(jql),
		Fields:     createFieldsQueryParam(),
		MaxResults: createMaxResultQueryParam(),
		Username:   username,
		Password:   passwd,
	}

	main := processBundle(jiraReq, 0)

	// check if more items need to be loaded (total > max results)
	if main.Total > main.MaxResults {
		// prepare paging and add issues
		var bundles = main.Total / main.MaxResults
		if main.Total%main.MaxResults == 0 {
			bundles--
		}

		// load bundles unless all tickets are loaded
		for bundle := 1; bundle <= bundles; bundle++ {
			// calculate next startAt value
			startAt := bundle * main.MaxResults
			// get page from server
			bundled := processBundle(jiraReq, startAt)
			// replace
			main.Issues = append(main.Issues, bundled.Issues...)
		}
	}
	return main
}

func processBundle(jiraReq JiraRequest, startAt int) Jira {
	req, err := http.NewRequest("GET", jiraReq.BaseUrl, nil)
	req.SetBasicAuth(jiraReq.Username, jiraReq.Password)
	// add query parameters
	q := req.URL.Query()
	q.Add("jql", jiraReq.EncodedJql)
	q.Add("fields", jiraReq.Fields)
	q.Add("maxResults", strconv.Itoa(jiraReq.MaxResults))
	q.Add("startAt", strconv.Itoa(startAt))
	// encode and assign back to the request
	req.URL.RawQuery = q.Encode()
	// send request to server
	resp, err := jiraReq.Client.Do(req)
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
