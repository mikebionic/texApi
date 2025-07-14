package scheduler

import (
	"context"
	"log"
	"strconv"
	"texApi/internal/services"
	"time"

	"texApi/database"
)

type AnalyticsScheduler struct {
	ticker   *time.Ticker
	quit     chan bool
	interval time.Duration
}

func NewAnalyticsScheduler() *AnalyticsScheduler {
	return &AnalyticsScheduler{
		quit: make(chan bool),
	}
}

func (s *AnalyticsScheduler) Start() error {
	log.Println("Starting Analytics Scheduler...")

	// Get configuration
	interval, err := s.getLogInterval()
	if err != nil {
		log.Printf("Error getting log interval, using default 24h: %v", err)
		interval = 24 * time.Hour
	}

	s.interval = interval

	if s.shouldRunAnalytics() {
		log.Println("Running initial analytics generation...")
		if err := services.GenerateAnalytics(); err != nil {
			log.Printf("Error in initial analytics generation: %v", err)
		}
	}

	s.ticker = time.NewTicker(s.interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				if s.shouldRunAnalytics() {
					log.Println("Running scheduled analytics generation...")
					if err := services.GenerateAnalytics(); err != nil {
						log.Printf("Error in scheduled analytics generation: %v", err)
					}
				}
			case <-s.quit:
				log.Println("Analytics scheduler stopped")
				return
			}
		}
	}()

	log.Printf("Analytics scheduler started with interval: %v", s.interval)
	return nil
}

func (s *AnalyticsScheduler) Stop() {
	log.Println("Stopping Analytics Scheduler...")
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.quit <- true
}

func (s *AnalyticsScheduler) shouldRunAnalytics() bool {
	enabled, err := s.isAnalyticsEnabled()
	if err != nil || !enabled {
		return false
	}

	lastRun, err := s.getLastRunTime()
	if err != nil {
		log.Printf("Error getting last run time: %v", err)
		return true // If we can't get last run time, assume we should run
	}

	nextRun := lastRun.Add(s.interval)
	return time.Now().After(nextRun)
}

func (s *AnalyticsScheduler) getLogInterval() (time.Duration, error) {
	var intervalDays string
	query := "SELECT value FROM tbl_analytics_config WHERE key = 'log_interval_days'"

	err := database.DB.QueryRow(context.Background(), query).Scan(&intervalDays)
	if err != nil {
		return 0, err
	}

	days, err := strconv.Atoi(intervalDays)
	if err != nil {
		return 0, err
	}

	return time.Duration(days) * 24 * time.Hour, nil
}

func (s *AnalyticsScheduler) isAnalyticsEnabled() (bool, error) {
	var enabled string
	query := "SELECT value FROM tbl_analytics_config WHERE key = 'enabled'"

	err := database.DB.QueryRow(context.Background(), query).Scan(&enabled)
	if err != nil {
		return false, err
	}

	return enabled == "true", nil
}

func (s *AnalyticsScheduler) getLastRunTime() (time.Time, error) {
	var lastRunStr string
	query := "SELECT value FROM tbl_analytics_config WHERE key = 'last_analytics_run'"

	err := database.DB.QueryRow(context.Background(), query).Scan(&lastRunStr)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, lastRunStr)
}

func (s *AnalyticsScheduler) UpdateInterval(newInterval time.Duration) {
	s.interval = newInterval

	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = time.NewTicker(newInterval)
	}

	log.Printf("Analytics scheduler interval updated to: %v", newInterval)
}

func (s *AnalyticsScheduler) GetStatus() map[string]interface{} {
	enabled, _ := s.isAnalyticsEnabled()
	lastRun, _ := s.getLastRunTime()

	return map[string]interface{}{
		"enabled":    enabled,
		"interval":   s.interval.String(),
		"last_run":   lastRun,
		"next_run":   lastRun.Add(s.interval),
		"is_running": s.ticker != nil,
	}
}
