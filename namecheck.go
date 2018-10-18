package namecheck

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/uwuh/namecheck/util"
)

func check(ctx context.Context, name, URL string, ct ChannelType) (ChannelStatus, error) {
	resp, err := util.HTTPClient.Get(strings.Replace(URL, "{name}", name, 1))
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return StatusAvailable, nil
		} else if strings.Contains(err.Error(), "Timeout") {
			return StatusUnknown, errors.New("Timeout")
		}

		return StatusUnknown, err
	}

	if resp.StatusCode == http.StatusOK {
		return StatusNotAvailable, nil
	}

	return StatusAvailable, nil
}

//Check name availability on all channel
func Check(ctx context.Context, name string, baseChannel []*Channel) (channels []*Channel, duration time.Duration) {
	start := time.Now()
	length := len(baseChannel)
	resultCh := make(chan *Channel, length)
	for _, channel := range baseChannel {
		go func(ch *Channel) {
			newChannel := ch
			newChannel.Status, newChannel.Error = check(ctx, name, ch.URL, ch.Type)
			resultCh <- newChannel
		}(channel)
	}

	for i := 0; i < length; i++ {
		channels = append(channels, <-resultCh)
	}
	return channels, time.Since(start)
}
