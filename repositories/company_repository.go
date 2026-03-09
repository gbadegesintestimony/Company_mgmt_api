package repositories

import (
	"company_mgmt_api/models"
	"context"
	"database/sql"
	"time"
)

type CompanyRepository struct {
	DB *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

func (r *CompanyRepository) Create(
	ctx context.Context,
	id string,
	name string,
	domain string,
) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO companies (id, name, domain, created_at)
		VALUES ($1, $2, $3, $4)
	`,
		id,
		name,
		domain,
		time.Now(),
	)

	return err
}

func (r *CompanyRepository) FindByDomain(
	ctx context.Context,
	domain string,
) (*models.Company, error) {

	var c models.Company

	err := r.DB.QueryRowContext(
		ctx,
		`SELECT id, name, domain
		FROM companies
		WHERE domain = $1`,
		domain,
	).Scan(&c.ID, &c.Name, &c.Domain)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CompanyRepository) FindByID(ctx context.Context, id string) (*models.Company, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT id, name, domain, created_at FROM companies WHERE id = $1`,
		id,
	)

	var c models.Company
	err := row.Scan(&c.ID, &c.Name, &c.Domain, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) Update(ctx context.Context, c *models.Company) error {
	query := `
		UPDATE companies 
		SET name=$1, domain=$2, status=$3 
		WHERE id=$4
	`
	_, err := r.DB.ExecContext(ctx, query, c.Name, c.Domain, c.Status, c.ID)
	return err
}
