package tinbu

import (
	"context"
	"net/http"
	"reflect"
	"time"
)

// GameUpdate holds updated value of a single lottery game object or error. If
// field Error is not nil, value of Game field should be ignored.
type GameUpdate struct {
	Game  Game
	Error error
}

// Client holds configuration of TinBu XML lottery feed client.
type Client struct {
	URL        string
	HTTPClient http.Client
}

// Load retrieves a snapshopt of current lottery feed state.
func (c *Client) Load(ctx context.Context) (map[string]Game, error) {
	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ParseXML(resp.Body)
}

// StreamUpdates checks lottery feed state periodically at given time interval
// and returns a channel detected changes. Lottery is marked as changed if any
// of field values have changed, or lottery was missing in previous update.
//
// Testing shows that stream will often send duplicate lottery game objects
// because stream often remove and later re-introduce same game objects.
func (c *Client) StreamUpdates(ctx context.Context, interval time.Duration) <-chan GameUpdate {
	ch := make(chan GameUpdate)

	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		last := c.checkForUpdates(ctx, nil, ch)
		for range ticker.C {
			last = c.checkForUpdates(ctx, last, ch)
		}
	}()
	return ch
}

func (c *Client) checkForUpdates(ctx context.Context, last map[string]Game, output chan<- GameUpdate) map[string]Game {
	current, err := c.Load(ctx)
	if err != nil {
		output <- GameUpdate{Error: err}
		return last
	}
	for id, game := range current {
		if reflect.DeepEqual(game, last[id]) {
			continue
		}
		output <- GameUpdate{Game: game}
	}
	return current
}
