package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sprint2-final-task/pkg/models"
	"strings"
	"time"
)

type Agent struct {
	orchestratorURL string
	client          *http.Client
}

func NewAgent(orchestratorURL string) *Agent {
	return &Agent{
		orchestratorURL: orchestratorURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  true,
				DisableKeepAlives:   false,
				MaxConnsPerHost:     100,
				MaxIdleConnsPerHost: 100,
			},
		},
	}
}

func (a *Agent) getTask() (*models.Task, error) {
	resp, err := a.client.Get(a.orchestratorURL + "/internal/task")
	if err != nil {
		if os.IsTimeout(err) || isConnectionRefused(err) {
			log.Printf("Оркестратор недоступен, ожидание...")
			time.Sleep(5 * time.Second)
			return nil, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var taskResp models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}

	return taskResp.Task, nil
}

func isConnectionRefused(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "context deadline exceeded"))
}
