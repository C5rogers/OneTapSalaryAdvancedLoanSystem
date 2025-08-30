package graph

import (
	"context"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
)

func (gc *LMSGraphClient) GetAllowedRoles() (roles []*models.Role, err error) {
	var query struct {
		Roles []*models.Role `graphql:"account_roles"`
	}
	variables := map[string]any{}

	if err = gc.Query(context.Background(), &query, variables); err != nil {
		return nil, err
	}

	if len(query.Roles) > 0 {
		roles = query.Roles
	}
	return roles, err
}

func (gc *LMSGraphClient) GetRoleByName(role string) (r *models.Role, err error) {
	var query struct {
		Roles []*models.Role `graphql:"account_roles(where:{role:{_eq:$role}})"`
	}

	variables := map[string]interface{}{
		"role": role,
	}
	if err := gc.Query(context.Background(), &query, variables); err != nil {
		return nil, err
	}
	if len(query.Roles) > 0 {
		r = query.Roles[0]
	}
	return r, nil
}
