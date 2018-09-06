package namecheck

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/uwuh/namecheck/util"
)

func check(name, URL string) (string, error) {
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
func Check(name string) (channels []Channel, duration time.Duration) {
	start := time.Now()
	length := len(Channels)
	resultCh := make(chan Channel, length)
	for _, channel := range Channels {
		go func(ch Channel) {
			newChannel := ch
			newChannel.Status, newChannel.Error = check(name, ch.URL)
			resultCh <- newChannel
		}(channel)
	}

	for index := 0; index < length; index++ {
		channels = append(channels, <-resultCh)
	}
	return channels, time.Since(start)
}
