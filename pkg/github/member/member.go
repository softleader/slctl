package member

import (
	"context"

	"github.com/google/go-github/v21/github"
)

// PublicizeOrganization Publicize the organization
func PublicizeOrganization(ctx context.Context, client *github.Client, org string) (err error) {
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return
	}
	_, err = client.Organizations.PublicizeMembership(ctx, org, user.GetLogin())
	return
}
