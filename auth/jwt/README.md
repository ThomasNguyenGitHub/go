# package auth/jwt

`package auth/jwt` provides a set of interfaces for service authorization
through [JSON Web Tokens](https://jwt.io/).

## Usage

NewParser takes a key function and an expected signing method and returns an
`endpoint.Middleware`. The middleware will parse a token passed into the
context via the `jwt.JWTTokenContextKey`. If the token is valid, any claims
will be added to the context via the `jwt.JWTClaimsContextKey`.

```go
import (
	stdjwt "github.com/dgrijalva/jwt-go"

	"github.com/ThomasNguyenGitHub/go/auth/jwt"
	"github.com/ThomasNguyenGitHub/go/endpoint"
)

func main() {
	var exampleEndpoint endpoint.Endpoint
	{
		kf := func(token *stdjwt.Token) (interface{}, error) { return []byte("SigningString"), nil }
		exampleEndpoint = MakeExampleEndpoint(service)
		exampleEndpoint = jwt.NewParser(kf, stdjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(exampleEndpoint)
	}
}
```

NewSigner takes a JWT key ID header, the signing key, signing method, and a
claims object. It returns an `endpoint.Middleware`. The middleware will build
the token string and add it to the context via the `jwt.JWTTokenContextKey`.

```go
import (
	stdjwt "github.com/dgrijalva/jwt-go"

	"github.com/ThomasNguyenGitHub/go/auth/jwt"
	"github.com/ThomasNguyenGitHub/go/endpoint"
)

func main() {
	var exampleEndpoint endpoint.Endpoint
	{
		exampleEndpoint = grpctransport.NewClient(...).Endpoint()
		exampleEndpoint = jwt.NewSigner(
			"kid-header",
			[]byte("SigningString"),
			stdjwt.SigningMethodHS256,
			jwt.Claims{},
		)(exampleEndpoint)
	}
}
```

In order for the parser and the signer to work, the authorization headers need
to be passed between the request and the context. `HTTPToContext()`,
`ContextToHTTP()`, `GRPCToContext()`, and `ContextToGRPC()` are given as
helpers to do this. These functions implement the correlating transport's
RequestFunc interface and can be passed as ClientBefore or ServerBefore
options.

Example of use in a client:

```go
import (
	stdjwt "github.com/dgrijalva/jwt-go"

	grpctransport "github.com/ThomasNguyenGitHub/go/transport/grpc"
	"github.com/ThomasNguyenGitHub/go/auth/jwt"
	"github.com/ThomasNguyenGitHub/go/endpoint"
)

func main() {

	options := []gotransport.ClientOption{}
	var exampleEndpoint endpoint.Endpoint
	{
		exampleEndpoint = grpctransport.NewClient(..., grpctransport.ClientBefore(jwt.ContextToGRPC())).Endpoint()
		exampleEndpoint = jwt.NewSigner(
			"kid-header",
			[]byte("SigningString"),
			stdjwt.SigningMethodHS256,
			jwt.Claims{},
		)(exampleEndpoint)
	}
}
```

Example of use in a server:

```go
import (
	"context"

	"github.com/ThomasNguyenGitHub/go/auth/jwt"
	"github.com/ThomasNguyenGitHub/go/log"
	grpctransport "github.com/ThomasNguyenGitHub/go/transport/grpc"
)

func MakeGRPCServer(ctx context.Context, endpoints Endpoints, logger Logger) pb.ExampleServer {
	options := []grpctransport.ServerOption{grpctransport.ServerErrorLogger(logger)}

	return &grpcServer{
		createUser: grpctransport.NewServer(
			ctx,
			endpoints.CreateUserEndpoint,
			DecodeGRPCCreateUserRequest,
			EncodeGRPCCreateUserResponse,
			append(options, grpctransport.ServerBefore(jwt.GRPCToContext()))...,
		),
		getUser: grpctransport.NewServer(
			ctx,
			endpoints.GetUserEndpoint,
			DecodeGRPCGetUserRequest,
			EncodeGRPCGetUserResponse,
			options...,
		),
	}
}
```
