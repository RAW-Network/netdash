package service

import (
	"encoding/json"
	"fmt"
	"netdash/internal/logger"
	"netdash/internal/model"
	"netdash/internal/repository"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type SpeedtestService interface {
	Run() (*model.TestResult, error)
	IsRunning() bool
}

type speedtestService struct {
	repo      repository.Repository
	isRunning bool
	mutex     sync.Mutex
}

func NewSpeedtestService(repo repository.Repository) SpeedtestService {
	return &speedtestService{
		repo: repo,
	}
}

func (s *speedtestService) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.isRunning
}

func (s *speedtestService) Run() (*model.TestResult, error) {
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		logger.Log("SPEED", "Test skipped: Another test is running")
		return nil, fmt.Errorf("test is already running")
	}
	s.isRunning = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		s.isRunning = false
		s.mutex.Unlock()
	}()

	logger.Log("SPEED", "Starting Ookla speedtest")

	conf, err := s.repo.GetConfig()
	if err != nil {
		return nil, err
	}

	args := []string{"--accept-license", "--accept-gdpr", "--format=json"}
	if conf.OoklaServerID != "" {
		args = append(args, "--server-id="+conf.OoklaServerID)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("speedtest.exe"); err == nil {
			cmd = exec.Command(".\\speedtest.exe", args...)
		} else {
			cmd = exec.Command("speedtest", args...)
		}
	} else {
		if _, err := os.Stat("speedtest"); err == nil {
			cmd = exec.Command("./speedtest", args...)
		} else {
			cmd = exec.Command("speedtest", args...)
		}
	}

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log("SPEED", "Execution failed: %v", err)
		return nil, fmt.Errorf("execution failed: %v", err)
	}

	var finalData *model.OoklaJSON
	lines := strings.Split(string(outputBytes), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var temp model.OoklaJSON
		if json.Unmarshal([]byte(line), &temp) == nil {
			if temp.Type == "result" {
				finalData = &temp
				break
			}
			if temp.Error != "" {
				logger.Log("SPEED", "Ookla Error: %s", temp.Error)
			}
		}
	}

	if finalData == nil {
		logger.Log("SPEED", "Failed to parse JSON result")
		return nil, fmt.Errorf("failed to parse speedtest result")
	}

	dlMbps := float64(finalData.Download.Bandwidth) * 8 / 1000000
	ulMbps := float64(finalData.Upload.Bandwidth) * 8 / 1000000

	res := &model.TestResult{
		Download:   dlMbps,
		Upload:     ulMbps,
		Ping:       finalData.Ping.Latency,
		PacketLoss: finalData.PacketLoss,
		ISP:        finalData.ISP,
		ServerID:   finalData.Server.ID,
		ServerName: finalData.Server.Name,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.SaveResult(res); err != nil {
		return nil, err
	}

	logger.Log("SPEED", "Result: DL %.2f Mbps | UL %.2f Mbps | Ping %.0f ms | Loss %.1f%% | ISP %s",
		dlMbps, ulMbps, finalData.Ping.Latency, finalData.PacketLoss, finalData.ISP)

	return res, nil
}