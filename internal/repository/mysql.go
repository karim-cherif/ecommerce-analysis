package repository

import (
	"database/sql"
	"fmt"
	"time"

	"ecommerce-analysis/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPurchaseEvents(startDate time.Time) ([]models.CustomerEventData, error) {
	query := `
        SELECT 
            ced.EventID,
            ced.CustomerID,
            ced.EventTypeID,
            ced.ContentID,
            ced.Quantity,
            ced.EventDate
        FROM CustomerEventData ced
        WHERE ced.EventTypeID = 6 
        AND ced.EventDate >= ?
    `

	rows, err := r.db.Query(query, startDate)
	if err != nil {
		return nil, fmt.Errorf("query purchase events: %w", err)
	}
	defer rows.Close()

	var events []models.CustomerEventData
	for rows.Next() {
		var event models.CustomerEventData
		err := rows.Scan(
			&event.EventID,
			&event.CustomerID,
			&event.EventTypeID,
			&event.ContentID,
			&event.Quantity,
			&event.EventDate,
		)
		if err != nil {
			return nil, fmt.Errorf("scan purchase event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *Repository) GetContentPrices() (map[int]float64, error) {
	query := `
        SELECT 
            cp.ContentID,
            cp.Price
        FROM ContentPrice cp
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query content prices: %w", err)
	}
	defer rows.Close()

	prices := make(map[int]float64)
	for rows.Next() {
		var contentID int
		var price float64
		if err := rows.Scan(&contentID, &price); err != nil {
			return nil, fmt.Errorf("scan content price: %w", err)
		}
		prices[contentID] = price
	}

	return prices, nil
}

func (r *Repository) GetCustomerEmails() (map[int64]string, error) {
	query := `
        SELECT 
            cd.CustomerID,
            cd.ChannelValue
        FROM CustomerData cd
        WHERE cd.ChannelTypeID = 1
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query customer emails: %w", err)
	}
	defer rows.Close()

	emails := make(map[int64]string)
	for rows.Next() {
		var customerID int64
		var email string
		if err := rows.Scan(&customerID, &email); err != nil {
			return nil, fmt.Errorf("scan customer email: %w", err)
		}
		emails[customerID] = email
	}

	return emails, nil
}

func (r *Repository) CreateExportTable(tableName string) error {
	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            CustomerID BIGINT PRIMARY KEY,
            Email VARCHAR(600),
            Revenue DECIMAL(10,2),
            UpdateDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `, tableName)

	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) ExportCustomerRevenue(tableName string, records []models.ExportRecord) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Prepare the statement for bulk insert
	stmt, err := tx.Prepare(fmt.Sprintf(`
        INSERT INTO %s (CustomerID, Email, Revenue)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE
        Revenue = VALUES(Revenue),
        UpdateDate = CURRENT_TIMESTAMP
    `, tableName))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute batch insert
	for _, record := range records {
		_, err := stmt.Exec(record.CustomerID, record.Email, record.Revenue)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("execute statement: %w", err)
		}
	}

	return tx.Commit()
}
