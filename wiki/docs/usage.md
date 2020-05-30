# Usage

## Register organization

!!! Tip "[http://go/org](http://go/org)"

`golinks` supports multiple organizations with JWT authentication.
So first, we have to register an organization.

![create_org](img/create_org.png)

## Add user

!!! Tip "[http://go/org/manage](http://go/org/manage)"

Add more users in your organization.

![create_user](img/create_user.png)

## Edit link

!!! Tip "[http://go/my.link](http://go/my.link) / [http://go/links/edit/my.link](http://go/links/edit/my.link)"

`golinks` automatically redirect to the edit page if the link doesn't exist. For
now, `golinks` supports 3 versions of links.

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

!!! Tip "[http://go/links](http://go/links)"

![links](img/links.png)
