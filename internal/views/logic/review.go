package logic

import "fmt"

func BuildReviewURL(exerciseID, quality int, fromSession bool, sessionID string) string {
	url := fmt.Sprintf("/exercise/%d/review?quality=%d", exerciseID, quality)
	if fromSession && sessionID != "" {
		url += fmt.Sprintf("&from=session&session=%s", sessionID)
	}
	return url
}
