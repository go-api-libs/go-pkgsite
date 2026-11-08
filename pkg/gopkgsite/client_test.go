package gopkgsite

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/go-api-libs/api"
)

type recordedRequest struct {
	method, path, rawQuery, userAgent string
	query                             url.Values
}

// serve starts a server answering every request with status, contentType and body, and returns a client for it
// together with a channel receiving each request it saw.
func serve(t *testing.T, status int, contentType, body string) (*Client, <-chan recordedRequest) {
	t.Helper()

	reqs := make(chan recordedRequest, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs <- recordedRequest{
			method:    r.Method,
			path:      r.URL.Path,
			rawQuery:  r.URL.RawQuery,
			userAgent: r.UserAgent(),
			query:     r.URL.Query(),
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	baseURL, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Fatal(err)
	}

	return c, reqs
}

// version is a caller's own result type, declaring only one member of any response.
type version struct {
	Version string `json:"version"`
}

// endpoints calls each operation with no query parameters, once into its own result type and once into a [version].
var endpoints = []struct {
	name   string
	call   func(*testing.T, *Client) error
	custom func(*testing.T, *Client) (*version, error)
}{
	{
		"GetImportedBy",
		func(t *testing.T, c *Client) error {
			_, err := c.GetImportedBy(t.Context(), "golang.org/x/time/rate", nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetImportedByWithResult[version](t.Context(), "golang.org/x/time/rate", nil)
		},
	},
	{
		"GetModule",
		func(t *testing.T, c *Client) error {
			_, err := c.GetModule(t.Context(), "golang.org/x/time", nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetModuleWithResult[version](t.Context(), "golang.org/x/time", nil)
		},
	},
	{
		"GetPackage",
		func(t *testing.T, c *Client) error {
			_, err := c.GetPackage(t.Context(), "golang.org/x/time/rate", nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetPackageWithResult[version](t.Context(), "golang.org/x/time/rate", nil)
		},
	},
	{
		"GetPackages",
		func(t *testing.T, c *Client) error {
			_, err := c.GetPackages(t.Context(), "golang.org/x/time", nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetPackagesWithResult[version](t.Context(), "golang.org/x/time", nil)
		},
	},
	{
		"GetSearch",
		func(t *testing.T, c *Client) error {
			_, err := c.GetSearch(t.Context(), nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetSearchWithResult[version](t.Context(), nil)
		},
	},
	{
		"GetSymbols",
		func(t *testing.T, c *Client) error {
			_, err := c.GetSymbols(t.Context(), "golang.org/x/time/rate", nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetSymbolsWithResult[version](t.Context(), "golang.org/x/time/rate", nil)
		},
	},
	{
		"GetVersions",
		func(t *testing.T, c *Client) error {
			_, err := c.GetVersions(t.Context(), "golang.org/x/time", nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetVersionsWithResult[version](t.Context(), "golang.org/x/time", nil)
		},
	},
	{
		"GetVulns",
		func(t *testing.T, c *Client) error {
			_, err := c.GetVulns(t.Context(), "golang.org/x/image", nil)
			return err
		},
		func(t *testing.T, c *Client) (*version, error) {
			return c.GetVulnsWithResult[version](t.Context(), "golang.org/x/image", nil)
		},
	},
}

func TestClient_Request(t *testing.T) {
	for _, tc := range []struct {
		name      string
		call      func(*testing.T, *Client) error
		wantPath  string
		wantQuery url.Values
	}{
		{
			name: "GetImportedBy",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetImportedBy(t.Context(), "golang.org/x/time/rate", &GetImportedByParams{
					Module: "golang.org/x/time", Version: "v0.15.0", Limit: 1, Token: "tok", Filter: "^k8s",
				})

				return err
			},
			wantPath: "/imported-by/golang.org/x/time/rate",
			wantQuery: url.Values{
				"module": {"golang.org/x/time"}, "version": {"v0.15.0"}, "limit": {"1"},
				"token": {"tok"}, "filter": {"^k8s"},
			},
		},
		{
			name: "GetModule",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetModule(t.Context(), "golang.org/x/time", &GetModuleParams{
					Version: "latest", Licenses: true, Readme: true,
				})

				return err
			},
			wantPath:  "/module/golang.org/x/time",
			wantQuery: url.Values{"version": {"latest"}, "licenses": {"true"}, "readme": {"true"}},
		},
		{
			name: "GetPackage",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetPackage(t.Context(), "golang.org/x/time/rate", &GetPackageParams{
					Module: "golang.org/x/time", Version: "master", Goos: "linux", Goarch: "arm64",
					Doc: "md", Examples: true, Imports: true, Licenses: true,
				})

				return err
			},
			wantPath: "/package/golang.org/x/time/rate",
			wantQuery: url.Values{
				"module": {"golang.org/x/time"}, "version": {"master"}, "goos": {"linux"}, "goarch": {"arm64"},
				"doc": {"md"}, "examples": {"true"}, "imports": {"true"}, "licenses": {"true"},
			},
		},
		{
			name: "GetPackages",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetPackages(t.Context(), "golang.org/x/time", &GetPackagesParams{
					Version: "v0.14.0", Limit: 5, Token: "next", Filter: "rate",
				})

				return err
			},
			wantPath:  "/packages/golang.org/x/time",
			wantQuery: url.Values{"version": {"v0.14.0"}, "limit": {"5"}, "token": {"next"}, "filter": {"rate"}},
		},
		{
			name: "GetSearch",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetSearch(t.Context(), &GetSearchParams{
					Q: "rate limiter", Symbol: "Limiter", Limit: 20, Token: "abc", Filter: "x/time",
				})

				return err
			},
			wantPath: "/search",
			wantQuery: url.Values{
				"q": {"rate limiter"}, "symbol": {"Limiter"}, "limit": {"20"}, "token": {"abc"}, "filter": {"x/time"},
			},
		},
		{
			name: "GetSymbols",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetSymbols(t.Context(), "golang.org/x/time/rate", &GetSymbolsParams{
					Module: "golang.org/x/time", Version: "main", Goos: "windows", Goarch: "386",
					Limit: 3, Token: "t", Filter: "Limit",
				})

				return err
			},
			wantPath: "/symbols/golang.org/x/time/rate",
			wantQuery: url.Values{
				"module": {"golang.org/x/time"}, "version": {"main"}, "goos": {"windows"}, "goarch": {"386"},
				"limit": {"3"}, "token": {"t"}, "filter": {"Limit"},
			},
		},
		{
			name: "GetVersions",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetVersions(t.Context(), "golang.org/x/time", &GetVersionsParams{
					Limit: 3, Token: "p2", Filter: "v0\\.1", Pseudo: true,
				})

				return err
			},
			wantPath:  "/versions/golang.org/x/time",
			wantQuery: url.Values{"limit": {"3"}, "token": {"p2"}, "filter": {"v0\\.1"}, "pseudo": {"true"}},
		},
		{
			name: "GetVulns",
			call: func(t *testing.T, c *Client) error {
				_, err := c.GetVulns(t.Context(), "golang.org/x/image", &GetVulnsParams{
					Module: "golang.org/x/image", Version: "v0.10.0", Limit: 1, Token: "v", Filter: "GO-2023",
				})

				return err
			},
			wantPath: "/vulns/golang.org/x/image",
			wantQuery: url.Values{
				"module": {"golang.org/x/image"}, "version": {"v0.10.0"}, "limit": {"1"}, "token": {"v"},
				"filter": {"GO-2023"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, reqs := serve(t, http.StatusOK, "application/json", "{}")

			if err := tc.call(t, c); err != nil {
				t.Fatal(err)
			}

			got := <-reqs
			if got.method != http.MethodGet {
				t.Errorf("method: got %q, want GET", got.method)
			}

			if got.path != tc.wantPath {
				t.Errorf("path: got %q, want %q", got.path, tc.wantPath)
			}

			if got.userAgent != defaultUserAgent {
				t.Errorf("User-Agent: got %q, want %q", got.userAgent, defaultUserAgent)
			}

			if !reflect.DeepEqual(got.query, tc.wantQuery) {
				t.Errorf("query: got %v, want %v", got.query, tc.wantQuery)
			}
		})
	}
}

func TestClient_Request_NoParams(t *testing.T) {
	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			c, reqs := serve(t, http.StatusOK, "application/json", "{}")

			if err := ep.call(t, c); err != nil {
				t.Fatal(err)
			}

			if got := <-reqs; got.rawQuery != "" {
				t.Errorf("query: got %q, want none", got.rawQuery)
			}
		})
	}

	t.Run("zero params", func(t *testing.T) {
		c, reqs := serve(t, http.StatusOK, "application/json", "{}")

		if _, err := c.GetPackage(t.Context(), "fmt", &GetPackageParams{}); err != nil {
			t.Fatal(err)
		}

		if got := <-reqs; got.rawQuery != "" {
			t.Errorf("query: got %q, want none", got.rawQuery)
		}
	})
}

func TestWithBaseURL(t *testing.T) {
	reqs := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs <- r.URL.Path

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	t.Cleanup(srv.Close)

	baseURL, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	baseURL.Path = "v1"

	c, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.GetModule(t.Context(), "golang.org/x/time", nil); err != nil {
		t.Fatal(err)
	}

	if got, want := <-reqs, "/v1/module/golang.org/x/time"; got != want {
		t.Errorf("path: got %q, want %q", got, want)
	}
}

func TestNewClient_Defaults(t *testing.T) {
	var got *http.Request

	c, err := NewClient(WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		got = r
		return nil, errors.New("stop")
	})}))
	if err != nil {
		t.Fatal(err)
	}

	_, _ = c.GetSymbols(t.Context(), "fmt", nil)

	if got == nil {
		t.Fatal("no request sent")
	}

	if want := "https://pkg.go.dev/v1/symbols/fmt"; got.URL.String() != want {
		t.Errorf("URL: got %q, want %q", got.URL, want)
	}
}

func TestClient_Responses(t *testing.T) {
	ctx := t.Context()

	c, err := NewClient(WithHTTPClient(&http.Client{Transport: replay(t)}))
	if err != nil {
		t.Fatal(err)
	}

	pkg, err := c.GetPackage(ctx, "golang.org/x/time/rate", nil)
	if err != nil {
		t.Fatal(err)
	}

	if pkg.Name != "rate" || pkg.ModulePath != "golang.org/x/time" || pkg.Version != "v0.15.0" ||
		pkg.Synopsis != "Package rate provides a rate limiter." || pkg.Goos != "all" ||
		!isTrue(pkg.IsLatest) || !isFalse(pkg.IsStandardLibrary) || !isTrue(pkg.IsRedistributable) {
		t.Errorf("GetPackage: unexpected %+v", pkg)
	}

	mod, err := c.GetModule(ctx, "golang.org/x/time", nil)
	if err != nil {
		t.Fatal(err)
	}

	if want := time.Date(2026, 2, 11, 19, 14, 29, 0, time.UTC); !mod.CommitTime.Equal(want) {
		t.Errorf("GetModule: CommitTime %v, want %v", mod.CommitTime, want)
	}

	if mod.RepoURL != "https://cs.opensource.google/go/x/time" || !isTrue(mod.HasGoMod) || mod.Readme != nil {
		t.Errorf("GetModule: unexpected %+v", mod)
	}

	versions, err := c.GetVersions(ctx, "golang.org/x/time", &GetVersionsParams{Limit: 3})
	if err != nil {
		t.Fatal(err)
	}

	if got := len(versions.Items); got != 3 {
		t.Fatalf("GetVersions: got %d items, want 3", got)
	}

	if v := versions.Items[2]; v.Version != "v0.13.0" || v.LatestVersion != "v0.15.0" || !isFalse(v.Retracted) ||
		!isFalse(v.Deprecated) {
		t.Errorf("GetVersions: unexpected third item %+v", v)
	}

	if versions.Total == nil || *versions.Total != -1 || versions.NextPageToken == "" {
		t.Errorf("GetVersions: unexpected pagination total=%v token=%q", versions.Total, versions.NextPageToken)
	}

	pkgs, err := c.GetPackages(ctx, "golang.org/x/time", nil)
	if err != nil {
		t.Fatal(err)
	}

	if pkgs.Packages == nil || len(pkgs.Packages.Items) != 1 || pkgs.Packages.Items[0].Path != "golang.org/x/time/rate" ||
		pkgs.Packages.Total == nil || *pkgs.Packages.Total != 1 {
		t.Errorf("GetPackages: unexpected %+v", pkgs.Packages)
	}

	search, err := c.GetSearch(ctx, &GetSearchParams{Q: "xyzzy"})
	if err != nil {
		t.Fatal(err)
	}

	if len(search.Items) != 13 || search.Items[0] != (SearchResult{
		ModulePath:  "github.com/docktermj/go-xyzzy-helpers",
		PackagePath: "github.com/docktermj/go-xyzzy-helpers/logger",
		Synopsis:    "Package helper ...",
		Version:     "v0.2.2",
	}) {
		t.Errorf("GetSearch: unexpected %d items, first %+v", len(search.Items), search.Items)
	}

	syms, err := c.GetSymbols(ctx, "golang.org/x/time/rate", nil)
	if err != nil {
		t.Fatal(err)
	}

	if syms.Symbols == nil || len(syms.Symbols.Items) != 31 {
		t.Fatalf("GetSymbols: unexpected %+v", syms.Symbols)
	}

	if got, want := syms.Symbols.Items[3], (Symbol{
		Kind: "Method", Name: "Limiter.Allow", Parent: "Limiter", Synopsis: "func (lim *Limiter) Allow() bool",
	}); got != want {
		t.Errorf("GetSymbols: got %+v, want %+v", got, want)
	}

	importers, err := c.GetImportedBy(ctx, "golang.org/x/time/rate", &GetImportedByParams{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}

	if importers.Version != "latest" || importers.ImportedBy == nil || len(importers.ImportedBy.Items) != 10 ||
		importers.ImportedBy.Items[0] != "0xacab.org/leap/geoip-service/internal/service" ||
		importers.ImportedBy.Total == nil || *importers.ImportedBy.Total != 14348 ||
		importers.ImportedBy.NextPageToken == "" {
		t.Errorf("GetImportedBy: unexpected %+v", importers.ImportedBy)
	}

	vulns, err := c.GetVulns(ctx, "golang.org/x/image", nil)
	if err != nil {
		t.Fatal(err)
	}

	if vulns.Items != nil || vulns.Total == nil || *vulns.Total != 0 {
		t.Errorf("GetVulns: unexpected %+v", vulns)
	}
}

func TestClient_ErrorResponse(t *testing.T) {
	const body = `{"code":404,"message":"package not found","fixes":["check the path"],` +
		`"candidates":[{"packagePath":"golang.org/x/time/rate","modulePath":"golang.org/x/time"}]}`

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			t.Run("decoded", func(t *testing.T) {
				c, _ := serve(t, http.StatusNotFound, "application/json; charset=utf-8", body)

				err := ep.call(t, c)

				apiErr, ok := errors.AsType[*api.Error](err)
				if !ok {
					t.Fatalf("got %T (%v), want *api.Error", err, err)
				}

				if !apiErr.IsCustom || apiErr.StatusCode() != http.StatusNotFound {
					t.Errorf("got IsCustom=%t status=%d, want custom 404", apiErr.IsCustom, apiErr.StatusCode())
				}

				pkgErr, ok := errors.AsType[*Error](err)
				if !ok {
					t.Fatalf("got %T, want it to wrap *Error", err)
				}

				want := &Error{
					Candidates: Candidates{{ModulePath: "golang.org/x/time", PackagePath: "golang.org/x/time/rate"}},
					Code:       404,
					Fixes:      []string{"check the path"},
					Message:    "package not found",
				}
				if !reflect.DeepEqual(pkgErr, want) {
					t.Errorf("got %+v, want %+v", pkgErr, want)
				}
			})

			t.Run("unknown content type", func(t *testing.T) {
				c, _ := serve(t, http.StatusBadRequest, "text/html", "<html>bad request</html>")

				if err := ep.call(t, c); !errors.Is(err, api.ErrUnknownContentType) {
					t.Fatalf("got %v, want %v", err, api.ErrUnknownContentType)
				}
			})

			t.Run("invalid JSON", func(t *testing.T) {
				c, _ := serve(t, http.StatusInternalServerError, "application/json", "oops")

				if _, ok := errors.AsType[*api.DecodingError](ep.call(t, c)); !ok {
					t.Fatal("want a *api.DecodingError")
				}
			})

			t.Run("unknown member", func(t *testing.T) {
				c, _ := serve(t, http.StatusBadRequest, "application/json", `{"code":400,"hint":"x"}`)

				if _, ok := errors.AsType[*api.DecodingError](ep.call(t, c)); !ok {
					t.Fatal("want a *api.DecodingError for a member the spec does not know")
				}
			})
		})
	}
}

func TestClient_SuccessDecoding(t *testing.T) {
	const body = `{"version":"v0.15.0","newField":1}`

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			t.Run("spec type rejects unknown members", func(t *testing.T) {
				c, _ := serve(t, http.StatusOK, "application/json", body)

				if _, ok := errors.AsType[*api.DecodingError](ep.call(t, c)); !ok {
					t.Fatal("want a *api.DecodingError")
				}
			})

			t.Run("own type is decoded leniently", func(t *testing.T) {
				c, _ := serve(t, http.StatusOK, "application/json", body)

				got, err := ep.custom(t, c)
				if err != nil {
					t.Fatal(err)
				}

				if got.Version != "v0.15.0" {
					t.Errorf("got %q, want %q", got.Version, "v0.15.0")
				}
			})

			t.Run("content type parameters are ignored", func(t *testing.T) {
				c, _ := serve(t, http.StatusOK, "application/json; charset=utf-8", "{}")

				if err := ep.call(t, c); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("unknown content type", func(t *testing.T) {
				c, _ := serve(t, http.StatusOK, "text/plain", "v0.15.0")

				if err := ep.call(t, c); !errors.Is(err, api.ErrUnknownContentType) {
					t.Fatalf("got %v, want %v", err, api.ErrUnknownContentType)
				}
			})
		})
	}
}

func TestError_Error(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "code and message",
			err:  &Error{Code: 400, Message: "bad request"},
			want: "400: bad request",
		},
		{
			name: "fixes",
			err:  &Error{Code: 400, Message: "bad request", Fixes: []string{"add a version", "drop the filter"}},
			want: "400: bad request\n\nFixes:\n- add a version\n- drop the filter",
		},
		{
			name: "candidates",
			err: &Error{Code: 404, Message: "not found", Candidates: Candidates{
				{ModulePath: "golang.org/x/time", PackagePath: "golang.org/x/time/rate"},
				{ModulePath: "example.com/m", PackagePath: "example.com/m/rate"},
			}},
			want: "404: not found\n\nCandidates:\n- pkg: golang.org/x/time/rate, module: golang.org/x/time" +
				"\n- pkg: example.com/m/rate, module: example.com/m",
		},
		{
			name: "fixes and candidates",
			err: &Error{
				Code: 404, Message: "not found", Fixes: []string{"try rate"},
				Candidates: Candidates{{ModulePath: "golang.org/x/time", PackagePath: "golang.org/x/time/rate"}},
			},
			want: "404: not found\n\nFixes:\n- try rate\n\nCandidates:\n- pkg: golang.org/x/time/rate, module: golang.org/x/time",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

func isTrue(b *bool) bool  { return b != nil && *b }
func isFalse(b *bool) bool { return b != nil && !*b }
