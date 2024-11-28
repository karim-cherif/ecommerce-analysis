package service

import (
	"fmt"
	"math"
	"sort"
	"time"

	"ecommerce-analysis/internal/models"
	"ecommerce-analysis/internal/repository"
	"ecommerce-analysis/internal/utils"
)

type Analyzer struct {
	repo     *repository.Repository
	logger   *utils.Logger
	quantile float64
}

func NewAnalyzer(repo *repository.Repository, logger *utils.Logger, quantile float64) *Analyzer {
	return &Analyzer{
		repo:     repo,
		logger:   logger,
		quantile: quantile,
	}
}

func (a *Analyzer) AnalyzeCustomerRevenue() error {
	startTime := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)

	// Load data
	a.logger.Info("Loading purchase events...")
	events, err := a.repo.GetPurchaseEvents(startTime)
	if err != nil {
		return fmt.Errorf("get purchase events: %w", err)
	}
	a.logger.Info("Loaded %d purchase events", len(events))

	a.logger.Info("Loading content prices...")
	prices, err := a.repo.GetContentPrices()
	if err != nil {
		return fmt.Errorf("get content prices: %w", err)
	}
	a.logger.Info("Loaded prices for %d products", len(prices))

	a.logger.Info("Loading customer emails...")
	emails, err := a.repo.GetCustomerEmails()
	if err != nil {
		return fmt.Errorf("get customer emails: %w", err)
	}
	a.logger.Info("Loaded emails for %d customers", len(emails))

	// Calculate revenue per customer
	customerRevenue := make(map[int64]float64)
	for _, event := range events {
		price, exists := prices[event.ContentID]
		if !exists {
			a.logger.Warning("No price found for content ID: %d", event.ContentID)
			continue
		}

		revenue := float64(event.Quantity) * price
		customerRevenue[event.CustomerID] += revenue
	}

	// Create sorted slice of customers by revenue
	type customerRevenueItem struct {
		CustomerID int64
		Revenue    float64
	}

	var sortedRevenues []customerRevenueItem
	for custID, revenue := range customerRevenue {
		sortedRevenues = append(sortedRevenues, customerRevenueItem{
			CustomerID: custID,
			Revenue:    revenue,
		})
	}

	sort.Slice(sortedRevenues, func(i, j int) bool {
		return sortedRevenues[i].Revenue > sortedRevenues[j].Revenue
	})

	// Print sample entries
	a.logger.Info("Sample of customer revenues:")
	sampleSize := int(math.Min(10, float64(len(sortedRevenues))))
	for i := 0; i < sampleSize; i++ {
		item := sortedRevenues[i]
		a.logger.Info("CustomerID: %d, Revenue: %.2f", item.CustomerID, item.Revenue)
	}

	// Calculate quantile statistics
	quantileSize := int(math.Ceil(float64(len(sortedRevenues)) * a.quantile))
	numQuantiles := int(1 / a.quantile)
	quantileStats := make(map[int]models.QuantileStats)

	for i := 0; i < numQuantiles; i++ {
		start := i * quantileSize
		end := int(math.Min(float64((i+1)*quantileSize), float64(len(sortedRevenues))))

		if start >= len(sortedRevenues) {
			break
		}

		stats := models.QuantileStats{
			NumberOfCustomers: end - start,
			MaxRevenue:        sortedRevenues[start].Revenue,
			MinRevenue:        sortedRevenues[end-1].Revenue,
		}
		quantileStats[i] = stats

		a.logger.Info("Quantile %d (%.1f%%-%.1f%%): Customers=%d, MaxRev=%.2f, MinRev=%.2f",
			i, float64(i)*a.quantile*100, float64(i+1)*a.quantile*100,
			stats.NumberOfCustomers, stats.MaxRevenue, stats.MinRevenue)
	}

	// Prepare export data for top quantile
	var exportRecords []models.ExportRecord
	for i := 0; i < quantileSize && i < len(sortedRevenues); i++ {
		customer := sortedRevenues[i]
		email := emails[customer.CustomerID]
		if email == "" {
			a.logger.Warning("No email found for customer ID: %d", customer.CustomerID)
			email = fmt.Sprintf("unknown_%d@example.com", customer.CustomerID)
		}

		exportRecords = append(exportRecords, models.ExportRecord{
			CustomerID: customer.CustomerID,
			Email:      email,
			Revenue:    customer.Revenue,
		})
	}

	// Export to database
	tableName := fmt.Sprintf("test_export_%s", time.Now().Format("20060102"))
	a.logger.Info("Creating export table: %s", tableName)
	if err := a.repo.CreateExportTable(tableName); err != nil {
		return fmt.Errorf("create export table: %w", err)
	}

	a.logger.Info("Exporting %d records to database", len(exportRecords))
	if err := a.repo.ExportCustomerRevenue(tableName, exportRecords); err != nil {
		return fmt.Errorf("export customer revenue: %w", err)
	}

	return nil
}
