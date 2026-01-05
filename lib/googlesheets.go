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
	"github.com/SomeDebris/rsmships-go"
)

const (
	SHEET_SCOPE          = `https://www.googleapis.com/auth/spreadsheets`
	SHEET_SCOPE_RO       = `https://www.googleapis.com/auth/spreadsheets`
	DRRT_SPREADSHEET_ID  = `1RTxlvsUHe6RXdOzsFxBWvxhul5pyJ0O2VaDVaVl-Uow`
	TOKEN_FNAME          = `token.json`
	CREDENTIALS_FNAME    = `credentials_drrt.json`
	RANGE_MATCH_SCHEDULE = `Calc!A1:F`
	RANGE_SHIP_ENTRY     = `Ships!A2:B`
	RANGE_DATA_ENTRY     = `DATA_ENTRY!A2:J`
)

var DatasheetService *sheets.Service = nil

type DRRTDatasheet struct {
	Id                 string          `json:"id,omitempty"`
	MatchScheduleRange string          `json:"matchScheduleRange,omitempty"`
	ShipEntryRange     string          `json:"shipEntryRange,omitempty"`
	LogRange           string          `json:"logRange"`
	Service            *sheets.Service `json:"-"`
}
// TODO: add function for marshalling this to JSON

func NewDRRTDatasheet(id string, matchschedulerange string, shipentryrange string, logrange string) *DRRTDatasheet {
	p := new(DRRTDatasheet)
	p.Id = id
	p.MatchScheduleRange = matchschedulerange
	p.ShipEntryRange = shipentryrange
	p.LogRange = logrange
	srv, err := getSheetsService()
	if err != nil {
		slog.Warn("Failed to get google sheets service. Cannot assign value to this field.", "err", err)
		return p
	}
	p.Service = srv
	return p
}

func NewDRRTDatasheetDefaults() *DRRTDatasheet {
	return NewDRRTDatasheet(DRRT_SPREADSHEET_ID, RANGE_MATCH_SCHEDULE, RANGE_SHIP_ENTRY, RANGE_DATA_ENTRY)
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
	if _, err := fmt.Scan(&authCode); err != nil {
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
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		slog.Error("Error decoding token.", "err", err)
	}
	return tok, err
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

// Call an "Update" on the sheet with at the given range with the given values.
func (m *DRRTDatasheet) UpdateValues(therange string, values [][]any) (*sheets.UpdateValuesResponse, error) {
	valrange := sheets.ValueRange{Values: values}
	call := m.Service.Spreadsheets.Values.Update(m.Id, therange, &valrange)
	call.ValueInputOption("USER_ENTERED")
	return call.Do()
}

// Call a "Clear" on the sheet at the given range.
func (m *DRRTDatasheet) ClearValues(therange string) (*sheets.ClearValuesResponse, error) {
	req := sheets.ClearValuesRequest{}
	return m.Service.Spreadsheets.Values.Clear(m.Id, m.MatchScheduleRange, &req).Do()
}

func (m *DRRTDatasheet) BatchUpdateValues(theranges []string, values [][][]any) (*sheets.BatchUpdateValuesResponse, error) {
	if len(theranges) != len(values) {
		slog.Error("length of ranges is not equivalent to length of values.", "lengthranges", len(theranges), "lengthvalues", len(values))
		return nil, fmt.Errorf("length of ranges is not equivalent to length of values.")
	}
	valueranges := make([]*sheets.ValueRange, len(theranges))
	for i, _ := range valueranges {
		valueranges[i] = &sheets.ValueRange{Values: values[i], Range: theranges[i]}
	}
	req := sheets.BatchUpdateValuesRequest{Data: valueranges, ValueInputOption: "USER_ENTERED"}
	return m.Service.Spreadsheets.Values.BatchUpdate(m.Id, &req).Do()
}

func (m *DRRTDatasheet) BatchClearValues(theranges []string) (*sheets.BatchClearValuesResponse, error) {
	req := sheets.BatchClearValuesRequest{Ranges: theranges}
	return m.Service.Spreadsheets.Values.BatchClear(m.Id, &req).Do()
}

func (m *DRRTDatasheet) ClearMatchSchedule() (*sheets.ClearValuesResponse, error) {
	return m.ClearValues(m.MatchScheduleRange)
}

// TODO: make the clear request into an update request that replaces the cells
// with empty strings. Then, make both into a batch update.
func (m *DRRTDatasheet) UpdateMatchSchedule(schedule [][]any) error {
	// clear contents of the match schedule location
	respclear, err := m.ClearMatchSchedule()
	if err != nil {
		slog.Error("Failed to clear values.", "id", m.Id, "range", m.MatchScheduleRange, "err", err)
		return err
	}
	slog.Info("Deleted match schedule range.", "range", respclear.ClearedRange, "HTTPStatusCode", respclear.HTTPStatusCode, "id", respclear.SpreadsheetId)

	// update range with new match schedule
	respupdate, err := m.UpdateValues(m.MatchScheduleRange, schedule)
	if err != nil {
		slog.Error("Failed to update values.", "id", m.Id, "range", m.MatchScheduleRange, "err", err)
		return err
	}
	slog.Info("Successfully updated match schedule values.", "range", respupdate.UpdatedRange, "HTTPStatusCode", respupdate.HTTPStatusCode, "id", respupdate.SpreadsheetId)

	return nil
}

func (m *DRRTDatasheet) UpdateShipsList(ships []rsmships.Ship) error {
	theupdate := getShipAuthorNamePairInterface(ships)
	respupdate, err := m.UpdateValues(m.ShipEntryRange, theupdate)
	if err != nil {
		slog.Error("Failed to update values.", "id", m.Id, "range", m.MatchScheduleRange, "err", err)
		return err
	}
	slog.Info("Successfully updated ships list.", "range", respupdate.UpdatedRange, "HTTPStatusCode", respupdate.HTTPStatusCode, "id", respupdate.SpreadsheetId)
	return nil
}

func (m *DRRTDatasheet) UpdateShipsAndMatchSchedule(ships []rsmships.Ship, schedule [][]any) error {
	shipauthornamepairs := getShipAuthorNamePairInterface(ships)
	resp, err := m.BatchUpdateValues([]string{m.MatchScheduleRange, m.ShipEntryRange}, [][][]any{schedule, shipauthornamepairs})
	if err != nil {
		slog.Error("Error occured while updating ships list and match schedule.", "err", err)
		return err
	}
	slog.Info("Updated ships list and match schedule.", "id", resp.SpreadsheetId, "TotalUpdatedCells", resp.TotalUpdatedCells)
	return nil
}

func (m *DRRTDatasheet) ClearShipsAndMatchSchedule() error {
	resp, err := m.BatchClearValues([]string{m.MatchScheduleRange, m.ShipEntryRange})
	if err != nil {
		slog.Error("Failed to clear ships list and match schedule.", "err", err)
		return err
	}
	slog.Info("Cleared ships list and match schedule.", "ClearedRanges[0]", resp.ClearedRanges[0], "ClearedRanges[1]", resp.ClearedRanges[1])
	return nil
}

