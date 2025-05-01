package mappers

import (
	"sample-web/dto"
	"sample-web/models"
	"time"
)

func ToUserRoles(roles []string) []models.UserRole {
	if len(roles) == 0 {
		roles = []string{string(models.LandLord), string(models.Tenant)}
	}
	userRoles := make([]models.UserRole, len(roles))
	for i, role := range roles {
		userRoles[i] = models.UserRole(role)
	}
	return userRoles
}

func ToUserRolesString(roles []models.UserRole) []string {
	userRoles := make([]string, len(roles))
	for i, role := range roles {
		userRoles[i] = string(role)
	}
	return userRoles
}

func ToUserModel(dto dto.UserRequest) models.User {
	now := time.Now()
	return models.User{
		Email:       dto.Email,
		Password:    dto.Password,
		PhoneNumber: dto.PhoneNumber,
		Roles:       ToUserRoles(dto.Roles),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func ToUserResponse(model models.User) dto.UserResponse {
	return dto.UserResponse{
		Id:          model.Id.Hex(),
		Name:        model.Name,
		Email:       model.Email,
		PhoneNumber: model.PhoneNumber,
		Roles:       ToUserRolesString(model.Roles),
		CreatedAt:   model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   model.UpdatedAt.Format(time.RFC3339),
	}
}
