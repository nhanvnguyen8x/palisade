package transport

import "github.com/nhanvnguyen8x/palisade/internal/features/workspace/domain"

func ToWorkspaceResponse(workspace *domain.Workspace) WorkspaceResponse {
	return WorkspaceResponse{
		ID:             workspace.ID.String(),
		OrganizationID: workspace.OrganizationID.String(),
		Name:           workspace.Name,
		CreatedAt:      workspace.CreatedAt,
		UpdatedAt:      workspace.UpdatedAt,
	}
}
