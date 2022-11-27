package sentinel

import (
	"testing"
)

// ISentinelSearcher
type mockSentinelSearcher struct {
}

func (m mockSentinelSearcher) Query(params SearchParameters) (QueryResponse, error) {
	return QueryResponse{}, nil
}

type mockDlEngine struct {
	path     string
	isOnline bool
}

func (m mockDlEngine) Download(productID string, dst string) (string, error) {
	return m.path, nil
}

func (m mockDlEngine) IsOnline(productID string) (bool, error) {
	return m.isOnline, nil
}
func TestNewClient(t *testing.T) {
	client := NewClient(mockSentinelSearcher{}, mockDlEngine{})
	if client == nil {
		t.Error("client is nil but should not be")
	}
}
