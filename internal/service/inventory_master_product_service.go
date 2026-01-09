package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dimensy-bridge/internal/dto"
	"dimensy-bridge/internal/model"
	"dimensy-bridge/internal/repository"
)

type InventoryMasterProductService interface {
	CreateInventory(ctx context.Context, req *dto.CreateInventoryMasterProductRequest) (*model.InventoryMasterProduct, error)
	UpdateInventory(ctx context.Context, id int64, req *dto.UpdateInventoryMasterProductRequest) (*model.InventoryMasterProduct, error)
	GetInventoryByID(ctx context.Context, id int64) (*model.InventoryMasterProduct, error)
	GetAllInventories(ctx context.Context, page, pageSize int) ([]*model.InventoryMasterProduct, int64, error)
	DeleteInventory(ctx context.Context, id int64) error
	AdjustStock(ctx context.Context, id int64, req *dto.AdjustStockRequest) (*model.InventoryMasterProduct, error)
	GetInventoryLogs(ctx context.Context, id int64, page, pageSize int) ([]*model.InventoryMasterProductLog, int64, error)
	MarkAsProcessed(ctx context.Context, id int64) (*model.InventoryMasterProduct, error)
	TogglePriority(ctx context.Context, id int64) (*model.InventoryMasterProduct, error)
	GetLowStockItems(ctx context.Context, threshold int) ([]*model.InventoryMasterProduct, error)
	GetTotalInventoryValue(ctx context.Context) (float64, error)
}

type inventoryMasterProductService struct {
	inventoryRepo repository.InventoryMasterProductRepository
	logRepo       repository.InventoryMasterProductLogRepository
}

func NewInventoryMasterProductService(
	inventoryRepo repository.InventoryMasterProductRepository,
	logRepo repository.InventoryMasterProductLogRepository,
) InventoryMasterProductService {
	return &inventoryMasterProductService{
		inventoryRepo: inventoryRepo,
		logRepo:       logRepo,
	}
}

func (s *inventoryMasterProductService) CreateInventory(ctx context.Context, req *dto.CreateInventoryMasterProductRequest) (*model.InventoryMasterProduct, error) {
	inventory := &model.InventoryMasterProduct{
		MasterProductID: req.MasterProductID,
		VendorName:      req.VendorName,
		Price:           req.Price,
		Quantity:        req.Quantity,
		CurrentStock:    req.Quantity,
		IsProcessed:     false,
		IsPriorityUse:   req.IsPriorityUse,
	}

	if err := s.inventoryRepo.Create(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	return inventory, nil
}

func (s *inventoryMasterProductService) UpdateInventory(ctx context.Context, id int64, req *dto.UpdateInventoryMasterProductRequest) (*model.InventoryMasterProduct, error) {
	inventory, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cek apakah sudah diproses
	if inventory.IsProcessed {
		return nil, errors.New("cannot update inventory that has been processed")
	}

	inventory.VendorName = req.VendorName
	inventory.Price = req.Price
	inventory.IsPriorityUse = req.IsPriorityUse

	if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	return inventory, nil
}

func (s *inventoryMasterProductService) GetInventoryByID(ctx context.Context, id int64) (*model.InventoryMasterProduct, error) {
	return s.inventoryRepo.FindByID(ctx, id)
}

func (s *inventoryMasterProductService) GetAllInventories(ctx context.Context, page, pageSize int) ([]*model.InventoryMasterProduct, int64, error) {
	return s.inventoryRepo.FindAll(ctx, page, pageSize)
}

func (s *inventoryMasterProductService) DeleteInventory(ctx context.Context, id int64) error {
	return s.inventoryRepo.Delete(ctx, id)
}

// AdjustStock handles stock adjustments and creates log entries
// If adjustment is positive (>0), it's a DEBIT (stock in)
// If adjustment is negative (<0), it's a CREDIT (stock out)
// Stock adjustment can only be done if is_processed is true
func (s *inventoryMasterProductService) AdjustStock(ctx context.Context, id int64, req *dto.AdjustStockRequest) (*model.InventoryMasterProduct, error) {
	inventory, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if inventory is processed
	if !inventory.IsProcessed {
		return nil, errors.New("stock adjustment can only be done after inventory has been processed")
	}

	previousStock := inventory.CurrentStock
	newStock := previousStock + req.Adjustment

	// Validate stock doesn't go negative
	if newStock < 0 {
		return nil, errors.New("insufficient stock for this adjustment")
	}

	// Update current stock
	inventory.CurrentStock = newStock
	inventory.Quantity = inventory.Quantity + req.Adjustment

	if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory stock: %w", err)
	}

	// Create log entry
	var debit, credit int
	if req.Adjustment > 0 {
		debit = req.Adjustment // Stock in = Debit
		credit = 0
	} else {
		debit = 0
		credit = -req.Adjustment // Stock out = Credit (converted to positive)
	}

	logEntry := &model.InventoryMasterProductLog{
		InventoryMasterProductID: inventory.ID,
		MasterProductID:          inventory.MasterProductID,
		Debit:                    debit,
		Credit:                   credit,
		PreviousStock:            previousStock,
		CurrentStock:             newStock,
		Time:                     time.Now(),
		Notes:                    &req.Notes,
	}

	if err := s.logRepo.Create(ctx, logEntry); err != nil {
		return nil, fmt.Errorf("failed to create log entry: %w", err)
	}

	return inventory, nil
}

func (s *inventoryMasterProductService) GetInventoryLogs(ctx context.Context, id int64, page, pageSize int) ([]*model.InventoryMasterProductLog, int64, error) {
	return s.logRepo.FindByInventoryID(ctx, id, page, pageSize)
}

func (s *inventoryMasterProductService) MarkAsProcessed(ctx context.Context, id int64) (*model.InventoryMasterProduct, error) {
	inventory, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if inventory.IsProcessed {
		return inventory, nil
	}

	inventory.IsProcessed = true

	if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to mark inventory as processed: %w", err)
	}

	// Create log entry for processed status - in ke debit
	logEntry := &model.InventoryMasterProductLog{
		InventoryMasterProductID: inventory.ID,
		MasterProductID:          inventory.MasterProductID,
		Debit:                    inventory.CurrentStock, // Stock in (processed) = Debit
		Credit:                   0,
		PreviousStock:            inventory.CurrentStock,
		CurrentStock:             inventory.CurrentStock,
		Time:                     time.Now(),
		Notes:                    ptrString("Marked as processed"),
	}

	if err := s.logRepo.Create(ctx, logEntry); err != nil {
		return nil, fmt.Errorf("failed to create log entry: %w", err)
	}

	return inventory, nil
}

func (s *inventoryMasterProductService) TogglePriority(ctx context.Context, id int64) (*model.InventoryMasterProduct, error) {
	inventory, err := s.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	inventory.IsPriorityUse = !inventory.IsPriorityUse

	if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to toggle priority: %w", err)
	}

	return inventory, nil
}

func (s *inventoryMasterProductService) GetLowStockItems(ctx context.Context, threshold int) ([]*model.InventoryMasterProduct, error) {
	return s.inventoryRepo.GetLowStockItems(ctx, threshold)
}

func (s *inventoryMasterProductService) GetTotalInventoryValue(ctx context.Context) (float64, error) {
	return s.inventoryRepo.GetTotalInventoryValue(ctx)
}

// Helper function to create string pointer
func ptrString(s string) *string {
	return &s
}
