package client

import "github.com/xanzy/go-gitlab"

var client *gitlab.Client

// InitClient initializes the gitlab client
func InitClient(baseURL, accessToken string) {
	client = gitlab.NewClient(nil, accessToken)
	client.SetBaseURL(baseURL)
}

// GetClient returns the initialized gitlabe client
func GetClient() *gitlab.Client {
	return client
}
