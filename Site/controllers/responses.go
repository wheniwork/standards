package controllers

type ErrorAPIResponse struct {
	Message string `json:"message,omitempty"`
	Success bool `json:"success"`
}

type APIResponse struct {
	Success bool `json:"success"`
	Results interface{} `json:"results,omitempty"`
}
