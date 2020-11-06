package main

type JQField struct {
	Name string `json:"name"`
	JQ   string `json:"jq"`
}

type JQFields []JQField

// Query is the JSON object that represents a query
type Query struct {
	Query  string   `json:"query"`
	Fields JQFields `json:"fields"`
}
