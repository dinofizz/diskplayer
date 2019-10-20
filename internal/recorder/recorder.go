package recorder

func Record(spotifyType string, spotifyId string) string {
	return "spotify:" + spotifyType + ":" + spotifyId
}
