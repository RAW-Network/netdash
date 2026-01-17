package model

import "time"

type TestResult struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Download   float64   `json:"download"`
	Upload     float64   `json:"upload"`
	Ping       float64   `json:"ping"`
	PacketLoss float64   `json:"packet_loss"`
	ISP        string    `json:"isp"`
	ServerID   int       `json:"server_id"`
	ServerName string    `json:"server_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type AppConfig struct {
	ID            uint   `gorm:"primaryKey"`
	OoklaServerID string `json:"ookla_server_id"`
	CronSchedule  string `json:"cron_schedule"`
	HistoryLimit  int    `json:"history_limit"`
}

type DBStats struct {
	TotalTests int64
	DBSize     string
}

type OoklaJSON struct {
	Type string `json:"type"`
	Ping struct {
		Latency float64 `json:"latency"`
	} `json:"ping"`
	PacketLoss float64 `json:"packetLoss"`
	Download   struct {
		Bandwidth int `json:"bandwidth"`
	} `json:"download"`
	Upload struct {
		Bandwidth int `json:"bandwidth"`
	} `json:"upload"`
	Server struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"server"`
	ISP   string `json:"isp"`
	Error string `json:"error"`
}