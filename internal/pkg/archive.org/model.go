package archiveorg

type AvailableSnapshorts struct {
	URL               string `json:"url"`
	ArchivedSnapshots struct {
		Closest struct {
			Status    string `json:"status"`
			Available bool   `json:"available"`
			URL       string `json:"url"`
			Timestamp string `json:"timestamp"`
		} `json:"closest"`
	} `json:"archived_snapshots"`
}
