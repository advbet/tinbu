package tinbu

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseXML(t *testing.T) {
	xml := `<?xml version="1.0" encoding="ISO-8859-1"?>
<allgames>
        <StateProv stateprov_name="California" stateprov_id="CA" country="U.S.A.">
                <game game_id="113" game_name="MEGA Millions" update_time="TUE 2007-10-16 23:08:37 EST">
                        <lastdraw_date>10/16/2007</lastdraw_date>
                        <lastdraw_numbers>01-02-05-41-44, Mega Ball: 25</lastdraw_numbers>
                        <nextdraw_date>10/19/2007</nextdraw_date>
                        <jackpot date="10/19/2007">34000000</jackpot>
                </game>
	</StateProv>
        <StateProv stateprov_name="Georgia" stateprov_id="GA" country="U.S.A.">
                <game game_id="113" game_name="MEGA Millions" update_time="TUE 2007-10-16 23:08:37 EST">
                        <lastdraw_date>10/16/2007</lastdraw_date>
                        <lastdraw_numbers>01-02-05-41-44, Mega Ball: 25</lastdraw_numbers>
                        <nextdraw_date>10/19/2007</nextdraw_date>
                        <jackpot date="10/19/2007">34000000</jackpot>
                </game>
                <game game_id="125" game_name="Win For Life" update_time="WED 2007-10-17 23:20:33 EST">
                        <lastdraw_date>10/17/2007</lastdraw_date>
                        <lastdraw_numbers>02-11-24-27-32-37, Free Ball: 19</lastdraw_numbers>
                        <nextdraw_date>10/20/2007</nextdraw_date>
                </game>
	</StateProv>
	<StateProv stateprov_name="Irish" stateprov_id="IE" country="United Kingdom">
                <game game_id="IE4" game_name="Daily Million 9PM" update_time="WED 2017-05-10 16:36:03 ET">
                        <lastdraw_date>05/10/2017</lastdraw_date>
                        <lastdraw_numbers>03-19-22-31-32-34, Bonus: 11</lastdraw_numbers>
                        <nextdraw_date>05/11/2017</nextdraw_date>
                </game>
	</StateProv>
        <StateProv stateprov_name="Atlantic Canada" stateprov_id="AC" country="Canada">
                <game game_id="203" game_name="Lotto Max" update_time="FRI 2019-01-25 23:44:35 $s">
                        <lastdraw_date>01/25/2019</lastdraw_date>
                        <lastdraw_numbers>06-12-20-24-40-42-47, Bonus: 44</lastdraw_numbers>
                        <nextdraw_date>02/01/2019</nextdraw_date>
                </game>
	</StateProv>
</allgames>`
	expected := map[string]Game{
		"113-10/16/2007": Game{
			ID:   "113",
			Name: "MEGA Millions",
			StateProvs: []StateProv{
				{"CA", "California", "U.S.A."},
				{"GA", "Georgia", "U.S.A."},
			},
			UpdateTime:      time.Date(2007, 10, 17, 3, 8, 37, 0, time.UTC),
			LastDrawDate:    time.Date(2007, 10, 16, 0, 0, 0, 0, time.UTC),
			LastDrawNumbers: "01-02-05-41-44, Mega Ball: 25",
			NextDrawDate:    time.Date(2007, 10, 19, 0, 0, 0, 0, time.UTC),
			Jackpot: &Jackpot{
				Date:   time.Date(2007, 10, 19, 0, 0, 0, 0, time.UTC),
				Amount: 34000000,
			},
		},
		"125-10/17/2007": Game{
			ID:   "125",
			Name: "Win For Life",
			StateProvs: []StateProv{
				{"GA", "Georgia", "U.S.A."},
			},
			UpdateTime:      time.Date(2007, 10, 18, 3, 20, 33, 0, time.UTC),
			LastDrawDate:    time.Date(2007, 10, 17, 0, 0, 0, 0, time.UTC),
			LastDrawNumbers: "02-11-24-27-32-37, Free Ball: 19",
			NextDrawDate:    time.Date(2007, 10, 20, 0, 0, 0, 0, time.UTC),
		},
		"IE4-05/10/2017": Game{
			ID:   "IE4",
			Name: "Daily Million 9PM",
			StateProvs: []StateProv{
				{"IE", "Irish", "United Kingdom"},
			},
			UpdateTime:      time.Date(2017, 05, 10, 20, 36, 03, 0, time.UTC),
			LastDrawDate:    time.Date(2017, 05, 10, 0, 0, 0, 0, time.UTC),
			LastDrawNumbers: "03-19-22-31-32-34, Bonus: 11",
			NextDrawDate:    time.Date(2017, 05, 11, 0, 0, 0, 0, time.UTC),
		},
		"203-01/25/2019": Game{
			ID:   "203",
			Name: "Lotto Max",
			StateProvs: []StateProv{
				{"AC", "Atlantic Canada", "Canada"},
			},
			UpdateTime:      time.Date(2019, 01, 26, 04, 44, 35, 0, time.UTC),
			LastDrawDate:    time.Date(2019, 01, 25, 0, 0, 0, 0, time.UTC),
			LastDrawNumbers: "06-12-20-24-40-42-47, Bonus: 44",
			NextDrawDate:    time.Date(2019, 02, 01, 0, 0, 0, 0, time.UTC),
		},
	}
	r := bytes.NewBufferString(xml)

	actual, err := ParseXML(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseXMLSample(t *testing.T) {
	f, err := os.Open("sample.xml")
	require.NoError(t, err)

	_, err = ParseXML(f)
	assert.NoError(t, err, "parsing sample data should not fail")
}

func TestParseXMLInvalidData(t *testing.T) {
	r := bytes.NewBufferString(`{}`)
	_, err := ParseXML(r)
	assert.Error(t, err, "invalid data should trigger parsing error")
}
