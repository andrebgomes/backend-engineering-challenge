package translationdeliverytime

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// ErrNoEvents is returned when input_file contains no events.
var ErrNoEvents = fmt.Errorf("no events")

// TranslationDeliveredJSON represents an event coming from the stream.
type TranslationDeliveredJSON struct {
	Timestamp      string `json:"timestamp"`
	TranslationID  string `json:"translation_id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	ClientName     string `json:"client_name"`
	EventName      string `json:"event_name"`
	Duration       int    `json:"duration"`
	NrWords        int    `json:"nr_words"`
}

// TranslationDelivered holds the necessary info to calculate the moving average.
type TranslationDelivered struct {
	Timestamp time.Time
	Duration  int
}

func (t *TranslationDelivered) UnmarshalJSON(bytes []byte) error {
	var td TranslationDeliveredJSON
	if err := json.Unmarshal(bytes, &td); err != nil {
		return err
	}
	time, err := time.Parse("2006-01-02 15:04:05.000000", td.Timestamp)
	if err != nil {
		return err
	}
	t.Timestamp = time
	t.Duration = td.Duration
	return nil
}

// durations is meant to hold the sum of events duration and their ammount.
type durations struct {
	sum int
	q   int
}

// eventsMap will contain an entry per minute from the events stream.
//
// The key indicates the minute (eg "2006-01-02 15:04").
// The value holds the durations sum and the ammount of events for that minute.
type eventsMap map[string]durations

// AverageDeliveryTime represents how a line of the output will be.
type AverageDeliveryTime struct {
	Date                time.Time
	AverageDeliveryTime float64
}

func (a AverageDeliveryTime) String() string {
	date := a.Date.Format("2006-01-02 15:04:05")

	if math.IsNaN(a.AverageDeliveryTime) {
		return fmt.Sprintf("{\"date\": \"%s\", \"average_delivery_time\": 0}", date)
	}
	i := math.Trunc(a.AverageDeliveryTime)
	iStr := fmt.Sprintf("%.1f", a.AverageDeliveryTime)
	if a.AverageDeliveryTime == i {
		iStr = fmt.Sprintf("%.0f", a.AverageDeliveryTime)
	}
	return fmt.Sprintf("{\"date\": \"%s\", \"average_delivery_time\": %s}", date, iStr)
}
