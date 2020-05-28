package traced

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"

	"github.com/haostudio/golinks/internal/auth"
)

// New returns a traced auth provider.
func New(p auth.Provider) auth.Provider {
	return &provider{p}
}

type provider struct {
	provider auth.Provider
}

func (p *provider) getSpan(ctx context.Context, name string) (
	context.Context, *trace.Span) {
	ctx, span := trace.StartSpan(ctx, name)
	defer span.End()
	span.AddAttributes(trace.StringAttribute("type", "auth_provider"))
	span.AddAttributes(trace.StringAttribute("provider", p.provider.String()))
	return ctx, span
}

// users
func (p *provider) GetUser(ctx context.Context, email string) (
	auth.User, error) {
	ctx, span := p.getSpan(ctx, "provider.GetUser")
	defer span.End()
	return p.provider.GetUser(ctx, email)
}

func (p *provider) GetUsers(ctx context.Context) ([]string, error) {
	ctx, span := p.getSpan(ctx, "provider.GetUsers")
	defer span.End()
	return p.provider.GetUsers(ctx)
}

func (p *provider) SetUser(ctx context.Context, user auth.User) error {
	ctx, span := p.getSpan(ctx, "provider.SetUser")
	defer span.End()
	return p.provider.SetUser(ctx, user)
}

func (p *provider) DeleteUser(ctx context.Context, email string) error {
	ctx, span := p.getSpan(ctx, "provider.DeleteUser")
	defer span.End()
	return p.provider.DeleteUser(ctx, email)
}

// organization
func (p *provider) GetOrg(ctx context.Context, name string) (
	auth.Organization, error) {
	ctx, span := p.getSpan(ctx, "provider.GetOrg")
	defer span.End()
	return p.provider.GetOrg(ctx, name)
}

func (p *provider) GetOrgUsers(ctx context.Context, name string) (
	[]string, error) {
	ctx, span := p.getSpan(ctx, "provider.GetOrgUsers")
	defer span.End()
	return p.provider.GetOrgUsers(ctx, name)
}

func (p *provider) SetOrg(ctx context.Context, org auth.Organization) error {
	ctx, span := p.getSpan(ctx, "provider.SetOrg")
	defer span.End()
	return p.provider.SetOrg(ctx, org)
}

func (p *provider) DeleteOrg(ctx context.Context, name string) error {
	ctx, span := p.getSpan(ctx, "provider.DeleteOrg")
	defer span.End()
	return p.provider.DeleteOrg(ctx, name)
}

// tokens
func (p *provider) GetToken(ctx context.Context, tokenStr string) (
	token auth.Token, err error) {
	ctx, span := p.getSpan(ctx, "provider.GetToken")
	defer span.End()
	return p.provider.GetToken(ctx, tokenStr)
}

func (p *provider) SetToken(ctx context.Context, token auth.Token) error {
	ctx, span := p.getSpan(ctx, "provider.SetToken")
	defer span.End()
	return p.provider.SetToken(ctx, token)
}

func (p *provider) DeleteToken(ctx context.Context, token string) error {
	ctx, span := p.getSpan(ctx, "provider.DeleteToken")
	defer span.End()
	return p.provider.DeleteToken(ctx, token)
}

func (p *provider) String() string {
	return fmt.Sprintf("traced(%s)", p.provider)
}
