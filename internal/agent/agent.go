package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sprint2-final-task/pkg/models"
	"strings"
	"sync"
	"time"
)

var (
	TIME_ADDITION_MS        int
	TIME_SUBTRACTION_MS     int
	TIME_MULTIPLICATIONS_MS int
	TIME_DIVISIONS_MS       int
	COMPUTING_POWER         int
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

func (a *Agent) submitResult(result models.TaskResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := a.client.Post(
		a.orchestratorURL+"/internal/task",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		if os.IsTimeout(err) || isConnectionRefused(err) {
			log.Printf("Оркестратор недоступен при отправке результата, повторная попытка...")
			time.Sleep(5 * time.Second)
			return a.submitResult(result)
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (a *Agent) processTask(task models.Task) error {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond) // Имитация длительного вычисления

	arg1 := task.Arg1
	arg2 := task.Arg2

	var result float64
	switch task.Operation {
	case "+":
		result = arg1 + arg2
	case "-":
		result = arg1 - arg2
	case "*":
		result = arg1 * arg2
	case "/":
		if arg2 == 0 {
			return fmt.Errorf("division by zero")
		}
		result = arg1 / arg2
	default:
		return fmt.Errorf("unknown operation: %s", task.Operation)
	}

	log.Printf("Task %s completed: %f %s %f = %f", task.ID, arg1, task.Operation, arg2, result)
	return a.submitResult(models.TaskResult{
		ID:     task.ID,
		Result: result,
	})
}

func (a *Agent) Worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, err := a.getTask()
		if err != nil {
			log.Printf("Error getting task: %v", err)
			continue
		}

		if task == nil {
			time.Sleep(time.Second)
			continue
		}

		if err := a.processTask(*task); err != nil {
			log.Printf("Error processing task: %v", err)
		}
	}
}
