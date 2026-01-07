package repository

import "netdash/internal/model"

type Repository interface {
	Init() error
	GetLatestResults(limit int) ([]model.TestResult, error)
	GetGraphData(limit int) ([]model.TestResult, error)
	GetAllResults() ([]model.TestResult, error)
	SaveResult(result *model.TestResult) error
	DeleteResult(id string) error
	ClearResults() error
	GetConfig() (*model.AppConfig, error)
	UpdateConfig(config *model.AppConfig) error
	GetDBStats() (*model.DBStats, error)
}