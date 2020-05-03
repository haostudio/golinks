package authtest

import (
	"context"

	"github.com/haostudio/golinks/internal/auth"
)

// CreateOrg creates org with admin.
func CreateOrg(provider auth.Provider, orgName, admin, pwd string) {
	var err error
	user, err := auth.NewUser(admin, pwd, orgName)
	if err != nil {
		panic(err)
	}
	org := auth.Organization{
		Name:       orgName,
		AdminEmail: admin,
	}
	err = provider.SetUser(context.Background(), *user)
	if err != nil {
		panic(err)
	}
	err = provider.SetOrg(context.Background(), org)
	if err != nil {
		panic(err)
	}
}
