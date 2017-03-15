package tinbu

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// StateProv is lottery location information. Single lottery might be associated
// with multiple locations for multi-state lotteries.
type StateProv struct {
	ID      string
	Name    string
	Country string
}

// Jackpot describes amount of money accumulated in lottery jackpot for a
// particular draw date.
type Jackpot struct {
	Date   time.Time
	Amount int64
}

// Game is main data structure of the feed. It describes a state of a single
// lottery.
type Game struct {
	ID              string
	Name            string
	StateProvs      []StateProv
	UpdateTime      time.Time
	LastDrawDate    time.Time
	LastDrawNumbers string
	NextDrawDate    time.Time
	Jackpot         *Jackpot
}

type jackpot struct {
	XMLName xml.Name `xml:"jackpot"`
	Date    string   `xml:"date,attr"`
	Amount  int64    `xml:",chardata"`
}

type game struct {
	XMLName      xml.Name `xml:"game"`
	ID           string   `xml:"game_id,attr"`
	Name         string   `xml:"game_name,attr"`
	Updated      string   `xml:"update_time,attr"`
	NextDrawDate string   `xml:"nextdraw_date"`
	LastDrawDate string   `xml:"lastdraw_date"`
	LastDraw     string   `xml:"lastdraw_numbers"`
	Jackpot      *jackpot `xml:"jackpot"`
}

type state struct {
	XMLName xml.Name `xml:"StateProv"`
	ID      string   `xml:"stateprov_id,attr"`
	Name    string   `xml:"stateprov_name,attr"`
	Country string   `xml:"country,attr"`
	Games   []game   `xml:"game"`
}

type document struct {
	XMLName xml.Name `xml:"allgames"`
	States  []state  `xml:"StateProv"`
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "ISO-8859-1":
		return charmap.ISO8859_1.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}
}

func parseKinkyTime(str string) (time.Time, error) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("Mon 2006-01-02 15:04:05 EST", str, loc)
}

func toGameMap(doc *document) (map[string]Game, error) {
	games := make(map[string]Game)

	for _, state := range doc.States {
		for _, game := range state.Games {
			id := fmt.Sprintf("%s-%s", game.ID, game.LastDrawDate)
			updated, err := parseKinkyTime(game.Updated)
			if err != nil {
				return nil, err
			}
			lastDraw, err := time.Parse("01/02/2006", game.LastDrawDate)
			if err != nil {
				return nil, err
			}
			nextDraw, err := time.Parse("01/02/2006", game.NextDrawDate)
			if err != nil {
				return nil, err
			}
			var jpot *Jackpot
			if game.Jackpot != nil {
				time, err := time.Parse("01/02/2006", game.Jackpot.Date)
				if err != nil {
					return nil, err
				}
				jpot = &Jackpot{
					Date:   time,
					Amount: game.Jackpot.Amount,
				}
			}
			old, ok := games[id]
			locations := append(old.StateProvs, StateProv{
				ID:      state.ID,
				Name:    state.Name,
				Country: state.Country,
			})
			if ok && old.LastDrawNumbers != game.LastDraw {
				return nil, fmt.Errorf("duplicate game %s instances with conflicting outcomes", game.ID)
			}
			games[id] = Game{
				ID:              game.ID,
				Name:            game.Name,
				StateProvs:      locations,
				UpdateTime:      updated.UTC(),
				LastDrawDate:    lastDraw,
				LastDrawNumbers: game.LastDraw,
				NextDrawDate:    nextDraw,
				Jackpot:         jpot,
			}
		}
	}

	return games, nil
}

func ParseXML(r io.Reader) (map[string]Game, error) {
	var doc document

	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charsetReader
	err := decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return toGameMap(&doc)
}
