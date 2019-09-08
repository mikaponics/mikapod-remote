/**
 *  - https://github.com/grpc/grpc-go/blob/master/examples/features/interceptor/client/main.go
 *  - https://github.com/kelseyhightower/grpc-hello-service
 */

package internal

import (
	"context"
	// "flag"
	// "fmt"
	// "io"
	// "log"
	// "time"

	// "golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	// "google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/credentials/oauth"
	// ecpb "google.golang.org/grpc/examples/features/proto/echo"
	// "google.golang.org/grpc/testdata"
)

// unaryInterceptor is an example unary interceptor.
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add the current bearer token to the metadata and call the RPC
    // command
    ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+GetAuthorizationBearerToken() )
    return invoker(ctx, method, req, reply, cc, opts...)
}
