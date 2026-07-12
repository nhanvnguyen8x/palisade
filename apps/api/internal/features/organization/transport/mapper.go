package transport

import "github.com/nhanvnguyen8x/palisade/internal/features/organization/domain"

func ToOrganizationResponse(org *domain.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
}
