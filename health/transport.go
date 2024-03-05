package health

import (
	httptransport "github.com/ThomasNguyenGitHub/go/transport/http"
	"net/http"

	"context"
)

// MakeHandler builds a go-base http transport and returns it
func MakeHandler() *httptransport.Server {
	e := makeHealthCheckEndpoint()

	healthHandler := httptransport.NewServer(
		e,
		decodeHealthCheckRequest,
		httptransport.EncodeJSONResponse,
	)
	return healthHandler
}

// decodeHealthCheckRequest returns an empty healthCheck request because there are no params for this request
func decodeHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return healthCheckRequest{}, nil
}
