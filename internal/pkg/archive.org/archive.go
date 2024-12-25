package archiveorg

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type service struct {
}

func New() *service {
	return &service{}
}

func (s *service) GetLastSnapshot(link string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://archive.org/wayback/available?url="+link, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var lastSnapshot AvailableSnapshorts
	err = json.NewDecoder(resp.Body).Decode(&lastSnapshot)
	if err != nil {
		return "", err
	}

	if lastSnapshot.ArchivedSnapshots.Closest.URL == "" {
		slog.Error("no snapshot for", slog.Any("link", link))
		return "", errors.New("no recent snapshot")
	}

	return lastSnapshot.ArchivedSnapshots.Closest.URL, nil
}
