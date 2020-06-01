# Usage

## Authorization

- [http://go/auth/login](http://go/auth/login): Login / Register
- [http://go/auth/org/register](http://go/auth/org/register): Create new organization
- [http://go/auth/org/manage](http://go/auth/org/manage): Add user to org

`golinks` supports multiple organizations with JWT authentication.
So first, we have to register an organization.

!!! TIP
    Skip authorization setup if run in NoAuth mode (`AUTHPROVIDER_NOAUTH_ENABLED=true`)

## Edit link

`golinks` automatically redirect to the edit page if the link doesn't exist. For
now, `golinks` supports 3 versions of links.

- [http://go/my.link](http://go/my.link) / [http://go/links/edit/my.link](http://go/links/edit/my.link)

!!! TIP
    - **v0 / Basic Mode**
      Simple URL redirect, e.g. https://go/xxx (`xxx -> https://github.com`) -> https://github.com
    - **v1 / Single-Parameter Mode**
      Redirect with `{}` in value replaced by text after `/` in URL path,
      e.g. https://go/xxx/haostudio/golinks (`xxx -> https://github.com/{}/issues`) -> https://github.com/haostudio/golinks/issues
    - **v2 / Multi-Parameter Mode**
      The text after `/` in URL path will be separated by `/` into a list of
      parameters and replace `{0}`, `{1}`, `{2}` ... in value, e.g.
      https://go/xxx/haostudio/golinks (`xxx -> https://github.com/{0}/{1}`) -> https://github.com/haostudio/golinks

![edit_link](img/edit_link.png)

## Show all links

- [http://go/links](http://go/links)

![links](img/links.png)
