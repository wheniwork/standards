package controllers

// The response wrapper objects.
type ErrorAPIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Results interface{} `json:"results,omitempty"`
}
