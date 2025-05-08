package defaultupdater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"vislab/collector"
	gitlabcollector "vislab/collector/gitlab"
)

type (
	Updater struct {
		collector collector.Collector
		port      string
	}
	Response struct {
		Status string `json:"status"`
	}
	UpdateRequest struct {
		ReleaseTag string `json:"release_tag"`
	}
)

func New(collector collector.Collector, port string) (*Updater, error) {
	if collector == nil {
		return nil, fmt.Errorf("collector not specified")
	}
	if port == "" {
		return nil, fmt.Errorf("port not specified")
	}

	return &Updater{
		collector: collector,
		port:      port,
	}, nil
}

func (u *Updater) Start(ctx context.Context) error {
	http.HandleFunc("/update", u.handleUpdate)

	return http.ListenAndServe(":"+u.port, nil)
}

func (u *Updater) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := u.collector.Update(r.Context(), gitlabcollector.WithReleaseTag(req.ReleaseTag)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := u.collector.Collect(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &Response{
		Status: "ok",
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
