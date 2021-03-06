package apiserver

import (
	"context"
	"google.golang.org/grpc"
)

type Client struct {
	Admin      AdminAPIClient
	Invocation WorkflowInvocationAPIClient
	Workflow   WorkflowAPIClient
	Consent    ConsentAPIClient
}

// Await blocks until the gRPC connection has been established
func (c *Client) Await(ctx context.Context) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return err
		default:
			_, err = c.Admin.Status(ctx, &Empty{})
			if err == nil {
				return nil
			}
		}
	}
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		Admin:      NewAdminAPIClient(conn),
		Invocation: NewWorkflowInvocationAPIClient(conn),
		Workflow:   NewWorkflowAPIClient(conn),
		Consent:    NewConsentAPIClient(conn),
	}
}

// Connect attempts to connect to the server at addr and returns a Client.
//
// addr should be of the format <hostname>:<port>, without a method. (e.g. workflows:5555)
func Connect(addr string) (*Client, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return NewClient(cc), nil
}
