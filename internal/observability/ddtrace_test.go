package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDDTraceEnabledAndAddr(t *testing.T) {
	t.Setenv("TRACING_PROVIDER", "")
	t.Setenv("DD_TRACE_ENABLED", "")
	t.Setenv("DD_AGENT_HOST", "")
	if DDTraceEnabled() {
		t.Fatal("expected DDTraceEnabled to be false when DD_AGENT_HOST is unset")
	}

	t.Setenv("DD_AGENT_HOST", "datakit-service.observability-prod.svc")
	t.Setenv("DD_TRACE_AGENT_PORT", "")
	if !DDTraceEnabled() {
		t.Fatal("expected DDTraceEnabled to be true when DD_AGENT_HOST is set")
	}
	if got := DDTraceAgentAddr(); got != "datakit-service.observability-prod.svc:9529" {
		t.Fatalf("got %q", got)
	}

	t.Setenv("DD_TRACE_ENABLED", "false")
	if DDTraceEnabled() {
		t.Fatal("expected DD_TRACE_ENABLED=false to disable ddtrace")
	}
}

func TestWrapHTTPHandler(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	for _, enabled := range []bool{false, true} {
		h := WrapHTTPHandler(base, "web-saas-content", "content.request", enabled)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("enabled=%v: got %d", enabled, rec.Code)
		}
	}
}
