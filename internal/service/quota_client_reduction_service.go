package service

import (
	"context"
	"dimensy-bridge/internal/repository"
	"errors"
	"time"
)

type QuotaClientReductionService interface {
	GetChart(
		ctx context.Context,
		rangeKey string,
		clientID *int64,
		masterProductID *int64,
	) (*ChartResponse, error)
}

type quotaClientReductionService struct {
	reductionRepo repository.QuotaClientReductionRepository
	masterRepo    repository.MasterProductRepository
}

type ChartResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

func NewQuotaClientReductionService(reductionRepo repository.QuotaClientReductionRepository, masterRepo repository.MasterProductRepository) QuotaClientReductionService {
	return &quotaClientReductionService{
		reductionRepo: reductionRepo,
		masterRepo:    masterRepo,
	}
}

func (s *quotaClientReductionService) GetChart(
	ctx context.Context,
	rangeKey string,
	clientID *int64,
	masterProductID *int64,
) (*ChartResponse, error) {

	start, end, err := resolveRange(rangeKey)
	if err != nil {
		return nil, err
	}

	// 1️⃣ ambil semua master product
	products, err := s.masterRepo.FindForChart(ctx, masterProductID)
	if err != nil {
		return nil, err
	}

	// 2️⃣ generate label waktu
	labels := generateTimeLabels(start, end, rangeKey)

	// 3️⃣ init matrix label × product
	matrix := map[string]map[string]interface{}{}
	for _, label := range labels {
		row := map[string]interface{}{
			"label": label,
		}
		for _, p := range products {
			row[p] = int64(0)
		}
		matrix[label] = row
	}

	// 4️⃣ ambil data aktual
	rows, err := s.reductionRepo.GetTimeSeries(
		ctx, start, rangeKey, clientID, masterProductID,
	)
	if err != nil {
		return nil, err
	}

	// 5️⃣ timpa value aktual
	for _, r := range rows {
		matrix[r.Label][r.MasterProductName] = r.TotalQuantity
	}

	// 6️⃣ convert ke slice (order by label)
	var result []map[string]interface{}
	for _, label := range labels {
		result = append(result, matrix[label])
	}

	return &ChartResponse{
		Success: true,
		Message: "Chart usage successfully retrieved",
		Data:    result,
	}, nil
}

func resolveRange(key string) (time.Time, time.Time, error) {
	now := time.Now()

	switch key {
	case "", "last_7_days":
		return now.AddDate(0, 0, -6), now, nil
	case "last_30_days":
		return now.AddDate(0, 0, -29), now, nil
	case "last_12_months":
		return now.AddDate(-1, 1, 0), now, nil
	default:
		return time.Time{}, time.Time{}, errors.New("invalid range")
	}
}

func generateTimeLabels(start, end time.Time, rangeKey string) []string {
	var labels []string
	cursor := start

	for !cursor.After(end) {
		if rangeKey == "last_12_months" {
			labels = append(labels, cursor.Format("2006-01"))
			cursor = cursor.AddDate(0, 1, 0)
		} else {
			labels = append(labels, cursor.Format("2006-01-02"))
			cursor = cursor.AddDate(0, 0, 1)
		}
	}
	return labels
}
