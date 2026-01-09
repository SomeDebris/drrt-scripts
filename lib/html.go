package lib

import (
	"bufio"
	"cmp"
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
	victory_template_path = filepath.Join(`html`, `victory_TEMPLATE.html`)
	game_template_path    = filepath.Join(`html`, `game_TEMPLATE.html`)
	next_template_path    = filepath.Join(`html`, `next_TEMPLATE.html`)
)

type MlogOverlayParse struct {
	Ranks       []int
	RankBoxes   []string
	Names       []string
	Authors     []string
	RankPoints  []int
	VictoryText []string
	MatchNumber int
	NextMatchNumber int
}

// disgusting
func NewMlogOverlayParse(ships []*rsmships.Ship, shipidxs []int, mlogs []*DRRTStandardMatchLog, ranks map[string]int, showrankdelta bool) *MlogOverlayParse {
	var out MlogOverlayParse
	lastmlog := mlogs[len(mlogs)-1]
	out.Ranks      = make([]int, lastmlog.AllianceLength * 2)
	out.RankBoxes  = make([]string, lastmlog.AllianceLength * 2)
	out.Names      = make([]string, lastmlog.AllianceLength * 2)
	out.Authors    = make([]string, lastmlog.AllianceLength * 2)
	out.RankPoints = make([]int, lastmlog.AllianceLength * 2)

	if lastmlog.Record[0].Result == Loss {
		out.VictoryText = []string{"LOSS", "VICTORY"}
	} else {
		out.VictoryText = []string{"VICTORY", "LOSS"}
	}

	out.MatchNumber = lastmlog.MatchNumber
	out.NextMatchNumber = lastmlog.MatchNumber + 1


	for i, idx := range shipidxs {
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

		out.RankPoints[i] = lastmlog.Record[i].RankPointsEarned

		if showrankdelta {
			// stats now
			stats := NewDRRTShipStats(lastmlog.ShipIndices[i], ships, mlogs)
			// stats before the last match
			statsprev := NewDRRTShipStats(lastmlog.ShipIndices[i], ships, mlogs[0:len(mlogs)-2])

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

func UpdateNextUp(ships []*rsmships.Ship, shipidxs []int, mlogs []*DRRTStandardMatchLog, ranks map[string]int) {
	p := *NewMlogOverlayParse(ships, shipidxs, mlogs, ranks, false)
	t, err := template.New("next_TEMPLATE.html").ParseFiles(next_template_path)
	if err != nil {
		slog.Error("Failed to parse template.", "err", err)
	}
	outfile, err := os.Create(filepath.Join("./html", "next.html"))
	if err != nil {
		slog.Error("Failed to open template output file", "err", err)
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = t.Execute(writer, p)
	if err != nil {
		slog.Error("Failed to save template output", "err", err)
	}
}
func UpdateGame(ships []*rsmships.Ship, shipidxs []int, mlogs []*DRRTStandardMatchLog, ranks map[string]int) {
	p := *NewMlogOverlayParse(ships, shipidxs, mlogs, ranks, false)
	t, err := template.New("game_TEMPLATE.html").ParseFiles(game_template_path)
	if err != nil {
		slog.Error("Failed to parse template.", "err", err)
	}
	outfile, err := os.Create(filepath.Join("./html", "game.html"))
	if err != nil {
		slog.Error("Failed to open template output file", "err", err)
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = t.Execute(writer, p)
	if err != nil {
		slog.Error("Failed to save template output", "err", err)
	}
}
func UpdateVictory(ships []*rsmships.Ship, mlogs []*DRRTStandardMatchLog, ranks map[string]int) {
	p := *NewMlogOverlayParse(ships, mlogs[len(mlogs)-1].ShipIndices, mlogs, ranks, false)
	t, err := template.New("victory_TEMPLATE.html").ParseFiles(victory_template_path)
	if err != nil {
		slog.Error("Failed to parse template.", "err", err)
	}
	outfile, err := os.Create(filepath.Join("./html", "victory.html"))
	if err != nil {
		slog.Error("Failed to open template output file", "err", err)
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	err = t.Execute(writer, p)
	if err != nil {
		slog.Error("Failed to save template output", "err", err)
	}
}
