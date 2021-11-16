package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	personalAccessToken = "PERSONAL TOKEN"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func RetrieveGitInfo() {
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(context.Background(), "samyak-jn")
	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		return
	}
	d, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Printf("json.MarshlIndent() failed with %s\n", err)
		return
	}
	fmt.Printf("User:\n%s\n", string(d))

	fileContent := []byte("This is the test content of my file\nand the 2nd line of it")
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("This is my test commit message"),
		Content:   fileContent,
		Branch:    github.String("main"),
		Committer: &github.CommitAuthor{Name: github.String("Samyak Jain"), Email: github.String("samyak.jn11@gmail.com")},
	}
	_, _, errCreate := client.Repositories.CreateFile(context.Background(), "samyak-jn", "webhook-test", "test_file.txt", opts)
	if errCreate != nil {
		fmt.Println(errCreate)
		return
	}
	fmt.Println("File Created")

	optsWebhook := &github.Hook{
		Name: github.String("web"),
		URL:  github.String("www.samyak-jn.com"),
		// Events:    []string{},
		// Active:    new(bool),
		Config: map[string]interface{}{
			"url":          "https://github.com/samyak-jn/webhook-test",
			"content_type": "json",
		},
		// ID:        new(int64),
	}
	_, _, errHook := client.Repositories.CreateHook(context.Background(), "samyak-jn", "webhook-test", optsWebhook)
	if errHook != nil {
		fmt.Println(errHook)
		return
	}
	fmt.Println("Hook Created")
}
