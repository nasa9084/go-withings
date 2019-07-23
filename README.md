go-withings
===

A Withings API Client for Go. This package does not support OAuth2 by itself, please use [golang.org/x/oauth2](https://golang.org/x/oauth2).

## Synopsis

``` go
tokenSrc := oauth2.StaticTokenSource(
    &oauth2.Token{
        AccessToken: "YOUR OAUTH ACCESS TOKEN",
    },
)
httpcl := oauth2.NewClient(oauth2.NoContext, tokenSrc)
c := withings.New(withings.WithHTTPClient(httpcl))
resp, _ := c.User().Getdevice()
fmt.Printf("%+v\n", resp)
```

## Features

### Google API style packge

All off the APIs in this packgeresemble that of google.golang.org/api.

### Full support for context.Context

The API is designed to use with context.Context.

## Status

- [x] User
  - [x] Getdevice
- [ ] Measure
  - [x] Getmeas
  - [x] Getactivity
  - [ ] Getintradayactivity
  - [ ] Getworkouts
- [ ] Sleep
- [ ] Notify
