package graph

import (
	"context"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/models"
)

type account_users_insert_input map[string]interface{}

func (gc *LMSGraphClient) GetUserByEmail(email string) (u *models.User, err error) {
	var query struct {
		Users []models.User `graphql:"account_users(where:{email:{_eq:$email}})"`
	}

	variables := map[string]interface{}{
		"email": email,
	}

	if err = gc.Query(context.Background(), &query, variables); err != nil {
		return nil, err
	}

	if len(query.Users) > 0 {
		u = &query.Users[0]
	}

	return u, err
}

func (gc *LMSGraphClient) RegisterUser(user *models.RegisteringUser) (u *models.User, err error) {

	var mutation struct {
		InsertUser struct {
			ID             string        `graphql:"id"`
			Email          string        `graphql:"email"`
			FullName       string        `graphql:"full_name"`
			PhoneNumber    string        `graphql:"phone_number"`
			ProfilePicture string        `graphql:"profile_picture"`
			Roles          []models.Role `graphql:"roles"`
		} `graphql:"insert_account_users_one(object:$object,on_conflict:{constraint:users_email_key,update_columns:[full_name]})"`
	}

	rolesData := []map[string]interface{}{
		{
			"role": user.Role, // provided by the user
		},
		{
			"role": "user:read", // extra default role
		},
	}

	variables := map[string]interface{}{
		"object": account_users_insert_input{
			"email":           user.Email,
			"full_name":       user.FullName,
			"phone_number":    user.PhoneNumber,
			"password":        user.Password,
			"profile_picture": user.ProfilePicture,
			"roles": map[string]interface{}{
				"data": rolesData,
			},
		},
	}

	if err = gc.Mutate(context.Background(), &mutation, variables); err != nil {
		return nil, err
	}
	u = &models.User{
		ID:          mutation.InsertUser.ID,
		Email:       mutation.InsertUser.Email,
		FullName:    mutation.InsertUser.FullName,
		PhoneNumber: mutation.InsertUser.PhoneNumber,
		Roles:       mutation.InsertUser.Roles,
	}
	return u, nil
}
