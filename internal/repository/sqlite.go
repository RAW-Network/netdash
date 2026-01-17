package repository

import (
	"fmt"
	"netdash/internal/logger"
	"netdash/internal/model"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type sqliteRepo struct {
	db     *gorm.DB
	dbPath string
}

func NewSQLiteRepository() Repository {
	return &sqliteRepo{
		dbPath: "data/netdash.db",
	}
}

func (r *sqliteRepo) Init() error {
	dir := filepath.Dir(r.dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	gl := gormLogger.New(
		nil,
		gormLogger.Config{
			LogLevel: gormLogger.Silent,
		},
	)

	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=auto_vacuum=FULL&_pragma=synchronous=NORMAL", r.dbPath)
	
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gl,
	})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	r.db = db
	
	if err := r.db.Exec("PRAGMA auto_vacuum = FULL").Error; err != nil {
		logger.Log("DB", "Warning: Failed to set auto_vacuum")
	}
	
	if err := r.db.AutoMigrate(&model.TestResult{}, &model.AppConfig{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	var count int64
	r.db.Model(&model.AppConfig{}).Count(&count)
	if count == 0 {
		logger.Log("DB", "Creating default configuration")
		r.db.Create(&model.AppConfig{
			OoklaServerID: "", 
			CronSchedule:  "0 * * * *",
			HistoryLimit:  10,
		})
	} else {
		r.db.Model(&model.AppConfig{}).Where("history_limit = 0").Update("history_limit", 10)
	}

	logger.Log("DB", "Database connected successfully")
	return nil
}

func (r *sqliteRepo) GetLatestResults(limit int) ([]model.TestResult, error) {
	var results []model.TestResult
	err := r.db.Order("created_at desc").Limit(limit).Find(&results).Error
	return results, err
}

func (r *sqliteRepo) GetGraphData(limit int) ([]model.TestResult, error) {
	var results []model.TestResult
	err := r.db.Order("created_at desc").Limit(limit).Find(&results).Error
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}
	return results, nil
}

func (r *sqliteRepo) GetAllResults() ([]model.TestResult, error) {
	var results []model.TestResult
	err := r.db.Order("created_at asc").Find(&results).Error
	return results, err
}

func (r *sqliteRepo) SaveResult(res *model.TestResult) error {
	return r.db.Create(res).Error
}

func (r *sqliteRepo) DeleteResult(id string) error {
	return r.db.Delete(&model.TestResult{}, id).Error
}

func (r *sqliteRepo) ClearResults() error {
	tx := r.db.Begin()
	if err := tx.Exec("DELETE FROM test_results").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	
	r.db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	return nil
}

func (r *sqliteRepo) GetConfig() (*model.AppConfig, error) {
	var conf model.AppConfig
	err := r.db.First(&conf).Error
	return &conf, err
}

func (r *sqliteRepo) UpdateConfig(config *model.AppConfig) error {
	var conf model.AppConfig
	if err := r.db.First(&conf).Error; err != nil {
		return err
	}
	
	conf.OoklaServerID = config.OoklaServerID
	conf.CronSchedule = config.CronSchedule
	conf.HistoryLimit = config.HistoryLimit
	return r.db.Save(&conf).Error
}

func (r *sqliteRepo) GetDBStats() (*model.DBStats, error) {
	var count int64
	r.db.Model(&model.TestResult{}).Count(&count)

	var totalSize int64
	dbDir := filepath.Dir(r.dbPath)
	dbName := filepath.Base(r.dbPath)
	
	files, _ := filepath.Glob(filepath.Join(dbDir, dbName+"*"))
	
	for _, f := range files {
		if info, err := os.Stat(f); err == nil {
			totalSize += info.Size()
		}
	}

	sizeStr := formatBytes(totalSize)

	return &model.DBStats{
		TotalTests: count,
		DBSize:     sizeStr,
	}, nil
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}