# Header Authorization
A simple middleware plugin for authorization using the request header by allowed values and methods to multiple endpoints.

The existing plugins can be browsed into the [Plugin Catalog](https://plugins.traefik.io).

# Configuration

### Static

```yaml
experimental:
  plugins:
    headauth:
      moduleName: "github.com/poloyacero/headauth"
      version: "v0.0.1"
```

### Dynamic

```yaml
http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin
  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  middlewares:
    my-plugin:
      plugin:
        headauth:
          header:
            name: X-Forward-Role
          allowed:
            - "admin"
          methods:
            - "GET"
            - "PATCH"
    another-plugin:
      plugin:
        example:
          header:
            name: X-Forward-Role
          allowed:
            - "admin"
            - "merchant"
          methods:
            - "POST"
```

### Local Mode

Traefik also offers a developer mode that can be used for temporary testing of plugins not hosted on GitHub.
To use a plugin in local mode, the Traefik static configuration must define the module name (as is usual for Go packages) and a path to a [Go workspace](https://golang.org/doc/gopath_code.html#Workspaces), which can be the local GOPATH or any directory.

The plugins must be placed in `./plugins-local` directory,
which should be in the working directory of the process running the Traefik binary.
The source code of the plugin should be organized as follows:

```
./plugins-local/
    └── src
        └── github.com
            └── traefik
                └── plugindemo
                    ├── demo.go
                    ├── demo_test.go
                    ├── go.mod
                    ├── LICENSE
                    ├── Makefile
                    └── readme.md
traefik.exe
```
