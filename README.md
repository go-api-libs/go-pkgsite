# Go Pkgsite (pkg.go.dev) API 

As announced [on the blog](https://go.dev/blog/pkgsite-api), [pkg.go.dev](https://pkg.go.dev/) now has an official API.

## API Issues

[The documentation](https://pkg.go.dev/v1beta/api) shows it's currently in beta. In fact, it has several issues:

The official OpenAPI spec is linked at [pkg.go.dev/v1beta/openapi.yaml](https://pkg.go.dev/v1beta/openapi.yaml) but is actually JSON. It is an invalid spec, e.g. missing path parameter definitions, unexported field in the spec (and therefore a reference to a non-existent schema `error`). At times it is quite repetitive. A cleaned up and improved version is available [here in this repo](./api/openapi.json).

Furthermore, for some routes, the documentation only provide example requests that fail (e.g. with error code 400 Bad Request).

Nevertheless, with this library, you can already use the API as best as possible. :)