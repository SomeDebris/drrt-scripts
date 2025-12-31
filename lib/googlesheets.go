package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	SHEET_SCOPE         = `https://www.googleapis.com/auth/spreadsheets`
	SHEET_SCOPE_RO      = `https://www.googleapis.com/auth/spreadsheets`
	DRRT_SPREADSHEET_ID = `1M7KgW7gcC1pA39ktOCYFZ4yFI2clbzMKrSPzYKQIs3g`
	TOKEN_FNAME         = `token.json`
	CREDENTIALS_FNAME   = `credentials.json`
)

var DatasheetService *sheets.Service = nil


// this was taken from the go quickstart
// https://developers.google.com/workspace/sheets/api/quickstart/go

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := TOKEN_FNAME
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the follo&wing link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token &from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return &tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}


func getSheetsService() (*sheets.Service, error) {
	ctx := context.Background()
	// read contents of secret into memory
	b, err := os.ReadFile(CREDENTIALS_FNAME)
	if err != nil {
		return nil, err
	}
	
	config, err := google.ConfigFromJSON(b, SHEET_SCOPE)
	if err != nil {
		return nil, err
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func GetGlobalSheetsService() (*sheets.Service, error) {
	if DatasheetService == nil {
		srv, err := getSheetsService()
		if err != nil {
			return nil, err
		}
		DatasheetService = srv
	}
	return DatasheetService, nil
}
