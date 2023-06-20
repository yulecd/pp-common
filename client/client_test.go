package client

import (
	"context"
	"testing"
)

func TestRequest_Get(t *testing.T) {
	resp, err := DefaultClient.Get("http://j.it603.com")
	if err != nil {
		t.Error("get error", err)
	} else {
		t.Log(resp.GetBody())
	}
}

func TestRequest_Wrapper(t *testing.T) {
	AddDefaultWrappers(func(next Wrapper) Wrapper {
		return func(ctx context.Context, req *Request) (*Response, error) {
			req.GetRequest().Header.Add("x-test-a", "test")
			t.Log("wrapper run")
			resp, err := next(ctx, req)
			if err != nil {
				t.Log("wrapper error:", err)
			} else {
				body, err := resp.GetBody()
				if err != nil {
					t.Log("wrapper body error:", err)
				} else {
					t.Log("wrapper get response", body.GetContents())
				}

			}
			return resp, err
		}
	})
	client := newHttpClient()
	resp, err := client.Get("http://j.it603.com")
	if err != nil {
		t.Error("get error", err)
	} else {
		t.Log(resp.GetBody())
	}
}
