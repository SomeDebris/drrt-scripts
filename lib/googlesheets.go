package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log/slog"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	SHEET_SCOPE         = `https://www.googleapis.com/auth/spreadsheets`
	SHEET_SCOPE_RO      = `https://www.googleapis.com/auth/spreadsheets`
	DRRT_SPREADSHEET_ID = `1RTxlvsUHe6RXdOzsFxBWvxhul5pyJ0O2VaDVaVl-Uow`
	TOKEN_FNAME         = `token.json`
	CREDENTIALS_FNAME   = `credentials.json`
)

var DatasheetService *sheets.Service = nil

type DRRTDatasheet struct {
	Id                 string          `json:"id,omitempty"`
	MatchScheduleRange string          `json:"matchScheduleRange,omitempty"`
	ShipEntryRange     string          `json:"shipEntryRange,omitempty"`
	LogRange           string          `json:"logRange"`
	Service            *sheets.Service `json:"-"`
}

// this was taken from the go quickstart
// https://developers.google.com/workspace/sheets/api/quickstart/go

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := TOKEN_FNAME
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		slog.Warn("Error reading token from file. Attempting to retrieve new token file.", "tokFile", tokFile, "err", err)
		tok, err2 := getTokenFromWeb(config)
		if err2 != nil {
			slog.Error("Could not get token from web.", "tokFile", tokFile, "err", err2)
			return nil, err2
		}
		err2 = saveToken(tokFile, tok)
		if err2 != nil {
			slog.Error("Could not save oauth token to file.", "tokFile", tokFile, "err", err2)
			return nil, err2
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the follo&wing link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(authCode); err != nil {
		slog.Error("Unable to read authorization code.", "err", err)
		return nil, err
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		slog.Error("Unable to retrieve token from web.", "err", err)
		return tok, err
	}
	return tok, nil
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
	if err != nil {
		slog.Error("Error decoding token.", "err", err)
	}
	return &tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	slog.Info("Saving credential file.", "path", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		slog.Error("Unable to cache oauth token.", "err", err)
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		slog.Error("Error encoding token to JSON.", "err", err)
	}
	return err
}


func getSheetsService() (*sheets.Service, error) {
	ctx := context.Background()
	// read contents of secret into memory
	b, err := os.ReadFile(CREDENTIALS_FNAME)
	if err != nil {
		slog.Error("Unable to read credentials filename.", "err", err)
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, SHEET_SCOPE)
	if err != nil {
		slog.Error("Unable to get service configuration.", "err", err)
		return nil, err
	}
	client, err := getClient(config)
	if err != nil {
		slog.Error("Unable to get client information.", "err", err)
		return nil, err
	}
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		slog.Error("Unable to create new sheets service.", "err", err)
	}
	return srv, err
}

func GetGlobalSheetsService() (*sheets.Service, error) {
	if DatasheetService == nil {
		slog.Info("Global datasheet service not set. Setting...")
		srv, err := getSheetsService()
		if err != nil {
			return nil, err
		}
		DatasheetService = srv
	}
	return DatasheetService, nil
}

// convinience functions for working with the DRRT datasheet
// TODO: make the clear request into an update request that replaces the cells
// with empty strings. Then, make both into a batch update.
func (m *DRRTDatasheet) UpdateMatchSchedule(schedule [][]string) error {
	// clear contents of the match schedule location
	req := sheets.ClearValuesRequest{}
	resp, err := m.srv.Spreadsheets.Values.Clear(m.id, m.MatchScheduleRange, &req).Do()
	if err != nil {
		slog.Error("Failed to clear values.", "id", m.id, "range", m.MatchScheduleRange, "err", err)
		return err
	}
	slog.Info("Deleted match schedule range.", "range", resp.ClearedRange, "HTTPStatusCode", resp.HTTPStatusCode, "id", resp.SpreadsheetId)

	// update range with new match schedule
	matchschedulevalrange := sheets.ValueRange{Values: Array2DStringsToInterface(schedule)}
	resp, err := m.srv.Spreadsheets.Values.Update(m.id, m.MatchScheduleRange, &matchschedulevalrange).Do()
	if err != nil {
		slog.Error("Failed to update values.", "id", m.id, "range", m.MatchScheduleRange, "err", err)
	}
	slog.Info("Updated match schedule range.", 


	return nil
}
