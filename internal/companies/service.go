package companies

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("company not found")
	ErrAlreadyExists = errors.New("company already exists")
)

// Company represents a company entity
type Company struct {
	ID        uuid.UUID
	Name      string
	LegalName string
	SIRET     *string
	VATNumber *string
	Status    string
}

// CreateCompanyRequest is the request to create a company
type CreateCompanyRequest struct {
	Name      string  `json:"name"`
	LegalName string  `json:"legal_name"`
	SIRET     *string `json:"siret,omitempty"`
	VATNumber *string `json:"vat_number,omitempty"`
	Status    string  `json:"status"`
}

// UpdateCompanyRequest is the request to update a company
type UpdateCompanyRequest struct {
	Name      string  `json:"name"`
	LegalName string  `json:"legal_name"`
	SIRET     *string `json:"siret,omitempty"`
	VATNumber *string `json:"vat_number,omitempty"`
	Status    string  `json:"status"`
}

// Service handles company business logic
type Service struct {
	db *sql.DB
}

// NewService creates a new company service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetByID retrieves a company by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Company, error) {
	query := `SELECT id, name, legal_name, siret, vat_number, status FROM companies WHERE id = $1`
	
	var company Company
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&company.ID,
		&company.Name,
		&company.LegalName,
		&company.SIRET,
		&company.VATNumber,
		&company.Status,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	if err != nil {
		return nil, err
	}
	
	return &company, nil
}

// List retrieves all active companies
func (s *Service) List(ctx context.Context) ([]*Company, error) {
	query := `SELECT id, name, legal_name, siret, vat_number, status FROM companies WHERE status = 'active' ORDER BY name`
	
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var companies []*Company
	for rows.Next() {
		var company Company
		if err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.LegalName,
			&company.SIRET,
			&company.VATNumber,
			&company.Status,
		); err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}
	
	return companies, nil
}

// Create creates a new company
func (s *Service) Create(ctx context.Context, req CreateCompanyRequest) (*Company, error) {
	query := `
		INSERT INTO companies (name, legal_name, siret, vat_number, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, legal_name, siret, vat_number, status
	`
	
	var company Company
	err := s.db.QueryRowContext(ctx, query,
		req.Name,
		req.LegalName,
		req.SIRET,
		req.VATNumber,
		req.Status,
	).Scan(
		&company.ID,
		&company.Name,
		&company.LegalName,
		&company.SIRET,
		&company.VATNumber,
		&company.Status,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &company, nil
}

// Update updates an existing company
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateCompanyRequest) (*Company, error) {
	query := `
		UPDATE companies
		SET name = $2, legal_name = $3, siret = $4, vat_number = $5, status = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, legal_name, siret, vat_number, status
	`
	
	var company Company
	err := s.db.QueryRowContext(ctx, query,
		id,
		req.Name,
		req.LegalName,
		req.SIRET,
		req.VATNumber,
		req.Status,
	).Scan(
		&company.ID,
		&company.Name,
		&company.LegalName,
		&company.SIRET,
		&company.VATNumber,
		&company.Status,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	if err != nil {
		return nil, err
	}
	
	return &company, nil
}

// Delete deletes a company
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM companies WHERE id = $1`
	
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return ErrNotFound
	}
	
	return nil
}
