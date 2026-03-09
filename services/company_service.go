package services

import (
	"company_mgmt_api/models"
	"company_mgmt_api/repositories"
	"context"
)

type CompanyService struct {
	repo      *repositories.CompanyRepository
	auditRepo *repositories.AuditRepository
}

func NewCompanyService(repo *repositories.CompanyRepository, auditRepo *repositories.AuditRepository) *CompanyService {
	return &CompanyService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *CompanyService) Get(ctx context.Context, companyID string) (*models.Company, error) {
	return s.repo.FindByID(ctx, companyID)
}

func (s *CompanyService) Update(ctx context.Context, c *models.Company, actorID string) error {
	err := s.repo.Update(ctx, c)
	if err != nil {
		return err
	}

	_ = s.auditRepo.Log(ctx, "company_update", actorID, c.ID, map[string]string{
		"name":   c.Name,
		"domain": c.Domain,
		"status": c.Status,
	})

	return nil
}
