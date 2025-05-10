package mappers

import (
	"sample-web/dto"
	"sample-web/models"
	"time"
)

func ToUserRole(role string) models.UserRole {
	switch role {
	case string(models.LandLord):
		return models.LandLord
	case string(models.Tenant):
		return models.Tenant
	default:
		return models.Tenant
	}
}

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
		PhoneNumber: dto.PhoneNumber,
		Roles:       ToUserRoles([]string{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func ToUserResponse(model models.User) dto.UserResponse {
	return dto.UserResponse{
		Id:          model.Id.Hex(),
		Name:        model.Name,
		PhoneNumber: model.PhoneNumber,
	}
}
