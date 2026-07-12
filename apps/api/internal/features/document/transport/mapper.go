package transport

import "github.com/nhanvnguyen8x/palisade/internal/features/document/application"

func ToUploadResponse(result *application.UploadResult) *UploadResponse {

	return &UploadResponse{
		ID:        result.ID.String(),
		FileName:  result.FileName,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
	}
}
