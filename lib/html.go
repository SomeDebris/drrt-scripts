package lib

import (
	"bufio"
	"cmp"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/SomeDebris/rsmships-go"
)

const (
	boxrankup_captain      = `rank_up_captain`
	boxrankneutral_captain = `rank_neutral_captain`
	boxrankdown_captain    = `rank_down_captain`
	boxrankup              = `rank_up`
	boxrankneutral         = `rank_neutral`
	boxrankdown            = `rank_down`
)

var (
	VICTORY_TEMPLATE_PATH = filepath.Join(`html`, `victory_TEMPLATE.html`)
	GAME_TEMPLATE_PATH    = filepath.Join(`html`, `game_TEMPLATE.html`)
	NEXT_TEMPLATE_PATH    = filepath.Join(`html`, `next_TEMPLATE.html`)
	victory_template      = template.Must(template.New("victory_TEMPLATE.html").ParseFiles(VICTORY_TEMPLATE_PATH))
	game_template         = template.Must(template.New("game_TEMPLATE.html").ParseFiles(GAME_TEMPLATE_PATH))
	next_template         = template.Must(template.New("next_TEMPLATE.html").ParseFiles(NEXT_TEMPLATE_PATH))
)

type StreamTemplateData struct {
	Ranks           []int
	RankBoxes       []string
	Names           []string
	Authors         []string
	RankPoints      []string
	VictoryText     []string
	MatchNumber     int
	NextMatchNumber int
}

// disgusting
func NewStreamTemplateDataQualifications(ships []*rsmships.Ship, shipIdxsToDisplay []int, mlogs []*DRRTStandardMatchLog, ranks map[string]int, showrankdelta bool) *StreamTemplateData {
	var out StreamTemplateData
	lastmlog := mlogs[len(mlogs)-1]
	out.Ranks      = make([]int, lastmlog.AllianceLength * 2)
	out.RankBoxes  = make([]string, lastmlog.AllianceLength * 2)
	out.Names      = make([]string, lastmlog.AllianceLength * 2)
	out.Authors    = make([]string, lastmlog.AllianceLength * 2)
	out.RankPoints = make([]string, lastmlog.AllianceLength * 2)

	if lastmlog.Record[0].Result == Loss {
		out.VictoryText = []string{"LOSS", "VICTORY!"}
	} else {
		out.VictoryText = []string{"VICTORY!", "LOSS"}
	}

	out.MatchNumber = lastmlog.MatchNumber
	out.NextMatchNumber = lastmlog.MatchNumber + 1


	for i, idx := range shipIdxsToDisplay {
		ship := ships[idx]
		out.Names[i] = ShipAuthorFromCommonNamefmt(ship.Data.Name)[0]
		out.Authors[i] = ship.Data.Author
		rank, ok := ranks[out.Names[i]]
		if !ok {
			slog.Error("Failed to index the ship name against the rank.")
		}
		out.Ranks[i] = rank
		captain := rank <= 8

		if captain {
			out.RankBoxes[i] = boxrankneutral_captain
		} else {
			out.RankBoxes[i] = boxrankneutral
		}

		out.RankPoints[i] = formatRankPointAdditions(lastmlog.Record[i].RankPointsEarned)

		if showrankdelta {
			// stats now
			stats := NewDRRTShipStats(idx, ships, mlogs)
			// stats before the last match
			statsprev := NewDRRTShipStats(idx, ships, mlogs[0:len(mlogs)-2])

			rankScoreCurrent := stats.RankingScore()
			rankScorePrev := statsprev.RankingScore()
			// did the ship rank up or down?
			rankdirection := cmp.Compare(rankScoreCurrent, rankScorePrev)
			if captain {
				switch rankdirection {
				case 1: //rank up
					out.RankBoxes[i] = boxrankup_captain
				case -1:
					out.RankBoxes[i] = boxrankdown_captain
				}
			} else {
				switch rankdirection {
				case 1: //rank up
					out.RankBoxes[i] = boxrankup
				case -1:
					out.RankBoxes[i] = boxrankdown
				}
			}
		}
	}

	return &out
}

// Add a + sign if ranking points were added, + if 0 were added, and - if they were somehow subtracted. It should never be the case that rank points are negative.
func formatRankPointAdditions(rankpointadd int) string {
	switch cmp.Compare(rankpointadd, 0) {
	case 1:
		return fmt.Sprintf("+%d", rankpointadd)
	case 0:
		return fmt.Sprintf("+%d", rankpointadd)
	default:
		return fmt.Sprintf("%d", rankpointadd)
	}
}

// TODO
// func NewStreamTemplateDataPlayoffs(alliances []*rsmships.Fleet)


func UpdateNextUpQualifications(outputPath string, ships []*rsmships.Ship, shipIdxsToDisplay []int, mlogs []*DRRTStandardMatchLog, ranks map[string]int) {
	p := *NewStreamTemplateDataQualifications(ships, shipIdxsToDisplay, mlogs, ranks, false)
	outfile, err := os.Create(outputPath)
	if err != nil {
		slog.Error("Failed to open template output file", "err", err)
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = next_template.Execute(writer, p)
	if err != nil {
		slog.Error("Failed to save template output", "err", err)
	}
}
func UpdateGameQualifications(outputPath string, ships []*rsmships.Ship, shipIdxsToDisplay []int, mlogs []*DRRTStandardMatchLog, ranks map[string]int) {
	p := *NewStreamTemplateDataQualifications(ships, shipIdxsToDisplay, mlogs, ranks, false)
	outfile, err := os.Create(outputPath)
	if err != nil {
		slog.Error("Failed to open template output file", "err", err)
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = game_template.Execute(writer, p)
	if err != nil {
		slog.Error("Failed to save template output", "err", err)
	}
}
func UpdateVictoryQualifications(outputPath string, ships []*rsmships.Ship, mlogs []*DRRTStandardMatchLog, ranks map[string]int) {
	p := *NewStreamTemplateDataQualifications(ships, mlogs[len(mlogs)-1].ShipIndices, mlogs, ranks, true)
	outfile, err := os.Create(outputPath)
	if err != nil {
		slog.Error("Failed to open template output file", "err", err)
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = victory_template.Execute(writer, p)
	if err != nil {
		slog.Error("Failed to save template output", "err", err)
	}
}
