package tinbu

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientLoad(t *testing.T) {
	respBody := `<?xml version="1.0" encoding="ISO-8859-1"?>
<allgames>
        <StateProv stateprov_name="California" stateprov_id="CA" country="U.S.A.">
                <game game_id="113" game_name="MEGA Millions" update_time="TUE 2007-10-16 23:08:37 EST">
                        <lastdraw_date>10/16/2007</lastdraw_date>
                        <lastdraw_numbers>01-02-05-41-44, Mega Ball: 25</lastdraw_numbers>
                        <nextdraw_date>10/19/2007</nextdraw_date>
                        <jackpot date="10/19/2007">34000000</jackpot>
                </game>
        </StateProv>
</allgames>`
	h := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(respBody))
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	c := Client{URL: srv.URL}
	expected := map[string]Game{
		"113-10/16/2007": Game{
			ID:   "113",
			Name: "MEGA Millions",
			StateProvs: []StateProv{
				{"CA", "California", "U.S.A."},
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
	}
	actual, err := c.Load(context.TODO())
	assert.NoError(t, err, "request should not fail")
	assert.Equal(t, expected, actual)
}

func TestClientLoadEmpty(t *testing.T) {
	respBody := `<?xml version="1.0" encoding="ISO-8859-1"?><allgames></allgames>`
	h := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(respBody))
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	c := Client{URL: srv.URL}
	actual, err := c.Load(context.TODO())
	assert.NoError(t, err, "request should not fail")
	assert.Equal(t, map[string]Game{}, actual)
}

func TestClientLoadFailedContext(t *testing.T) {
	parent, cancel := context.WithTimeout(context.TODO(), time.Second)
	respBody := `<?xml version="1.0" encoding="ISO-8859-1"?><allgames></allgames>`
	h := func(w http.ResponseWriter, r *http.Request) {
		<-parent.Done()
		_, _ = w.Write([]byte(respBody))
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()
	defer cancel()

	ctx, cancel := context.WithTimeout(parent, time.Millisecond)
	defer cancel()

	c := Client{URL: srv.URL}
	_, err := c.Load(ctx)
	assert.Error(t, err)
	assert.NoError(t, parent.Err(), "request was not cancelled by chlid context")
}

func TestClientLoadError(t *testing.T) {
	c := Client{URL: ":::"}
	_, err := c.Load(context.TODO())
	assert.Error(t, err, "invalid URL should fail")
}
