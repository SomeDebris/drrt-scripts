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
	Ranks           [][]int
	RankBoxes       [][]string
	Names           [][]string
	Authors         [][]string
	RankPoints      [][]string
	VictoryText     []string
	MatchNumber     int
	NextMatchNumber int
}

// disgusting
func NewStreamTemplateDataQualifications(ships []*rsmships.Ship, shipIdxsToDisplay []int, mlogs []*DRRTStandardMatchLog, nameToRank map[string]int, showrankdelta bool) *StreamTemplateData {
	var out StreamTemplateData
	lastmlog := mlogs[len(mlogs)-1]
	// there are 2 alliances in each match:
	out.Ranks       = make([][]int, 2)
	out.RankBoxes   = make([][]string, 2)
	out.Names       = make([][]string, 2)
	out.Authors     = make([][]string, 2)
	out.RankPoints  = make([][]string, 2)

	// there are two alliances:
	if lastmlog.Record[0].Result == Loss {
		out.VictoryText = []string{"RED ALLIANCE LOSS", "BLUE ALLIANCE VICTORY!"}
	} else {
		out.VictoryText = []string{"RED ALLIANCE VICTORY!", "BLUE ALLIANCE LOSS"}
	}

	out.MatchNumber     = lastmlog.MatchNumber
	out.NextMatchNumber = lastmlog.MatchNumber + 1

	// alliance index: 0 for left, 1 for right
	var alidx int
	for i, idx := range shipIdxsToDisplay {
		if i < len(shipIdxsToDisplay) / 2 {
			alidx = 0
		} else {
			alidx = 1
		}
		
		ship := ships[idx]
		name := ShipAuthorFromCommonNamefmt(ship.Data.Name)[0]
		out.Names[alidx]   = append(out.Names[alidx], name)
		out.Authors[alidx] = append(out.Authors[alidx], ship.Data.Author)

		rank, ok := nameToRank[name]
		if !ok {
			slog.Error("Failed to index the ship name against the rank.", "name", name, "author", ship.Data.Author)
		}
		out.Ranks[alidx] = append(out.Ranks[alidx], rank)

		rankpointsstring := formatRankPointAdditions(lastmlog.Record[i].RankPointsEarned)
		out.RankPoints[alidx] = append(out.RankPoints[alidx], rankpointsstring)

		var rankbox string
		captain := rank <= 8
		if captain {
			rankbox = boxrankneutral_captain
		} else {
			rankbox = boxrankneutral
		}

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
					rankbox = boxrankup_captain
				case -1:
					rankbox = boxrankdown_captain
				}
			} else {
				switch rankdirection {
				case 1: //rank up
					rankbox = boxrankup
				case -1:
					rankbox = boxrankdown
				}
			}
		}
		slog.Debug("Selected rankbox.", "box", rankbox)
		out.RankBoxes[alidx] = append(out.RankBoxes[alidx], rankbox)
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
		return fmt.Sprintf("-%d", rankpointadd)
	}
}

// TODO
func NewStreamTemplateDataPlayoffs(alliances []*rsmships.Fleet, nameToRank map[string]int, nCaptains int, victoryLeft bool) *StreamTemplateData {
	var out StreamTemplateData
	out.MatchNumber = 0
	out.NextMatchNumber = 0
	out.Ranks       = make([][]int, 2)
	out.RankBoxes   = make([][]string, 2)
	out.Names       = make([][]string, 2)
	out.Authors     = make([][]string, 2)
	out.RankPoints  = make([][]string, 2)
	// the first index of the fleet blueprints is the number of the alliance
	allianceNumbers := make([]int, len(alliances))
	for i, alliance := range alliances {
		slog.Error("WHAT", "blueprints", len(alliance.Blueprints))
		out.Ranks[i]       = make([]int,    len(alliance.Blueprints))
		out.RankBoxes[i]   = make([]string, len(alliance.Blueprints))
		out.Names[i]       = make([]string, len(alliance.Blueprints))
		out.Authors[i]     = make([]string, len(alliance.Blueprints))
		out.RankPoints[i]  = make([]string, len(alliance.Blueprints))

		for j, ship := range alliance.Blueprints {
			name := ShipAuthorFromCommonNamefmt(ship.Data.Name)[0]
			out.Names[i][j] = name

			out.Authors[i][j] = ship.Data.Author

			rank, ok := nameToRank[name]
			if !ok {
				slog.Error("ship name not present in ranks map.", "function", "NewStreamTemplateDataPlayoffs", "name", name, "author", ship.Data.Author)
			}
			out.Ranks[i][j] = rank

			if j <= 0 {
				allianceNumbers[i] = rank
			}

			var rankbox string
			captain := rank <= nCaptains
			if rank > 0 && captain {
				rankbox = boxrankneutral_captain
			} else {
				rankbox = boxrankneutral
			}
			slog.Debug("Selected rankbox.", "box", rankbox)
			out.RankBoxes[i][j] = rankbox

			out.RankPoints[i][j] = ""
		}
	}
	if victoryLeft {
		out.VictoryText = []string{fmt.Sprintf("ALLIANCE %d VICTORY!", allianceNumbers[0]), fmt.Sprintf("ALLIANCE %d LOSS", allianceNumbers[1])}
	} else {
		out.VictoryText = []string{fmt.Sprintf("ALLIANCE %d LOSS", allianceNumbers[0]), fmt.Sprintf("ALLIANCE %d VICTORY!", allianceNumbers[1])}
	}

	return &out
}

func writeTemplate(path string, p *StreamTemplateData, template *template.Template) error {
	outfile, err := os.Create(path)
	if err != nil {
		slog.Error("Failed to open template output file", "err", err)
		return err
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = template.Execute(writer, p)
	if err != nil {
		slog.Error("Failed to save template output", "err", err)
		return err
	}
	return nil
}

func UpdateNextUpPlayoffs(outputPath string, alliances []*rsmships.Fleet, nameToRank map[string]int, nCaptains int) {
	p := NewStreamTemplateDataPlayoffs(alliances, nameToRank, 8, false)
	template := next_template
	err := writeTemplate(outputPath, p, template)
	if err != nil {
		slog.Error("Error saving next-up.", "err", err)
	}
}
func UpdateGamePlayoffs(outputPath string, alliances []*rsmships.Fleet, nameToRank map[string]int, nCaptains int) {
	p := NewStreamTemplateDataPlayoffs(alliances, nameToRank, 8, false)
	template := game_template
	err := writeTemplate(outputPath, p, template)
	if err != nil {
		slog.Error("Error saving game.", "err", err)
	}
}
func UpdateVictoryPlayoffs(outputPath string, alliances []*rsmships.Fleet, nameToRank map[string]int, nCaptains int, victoryLeft bool) {
	p := NewStreamTemplateDataPlayoffs(alliances, nameToRank, 8, victoryLeft)
	template := victory_template
	err := writeTemplate(outputPath, p, template)
	if err != nil {
		slog.Error("Error saving game.", "err", err)
	}
}


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
