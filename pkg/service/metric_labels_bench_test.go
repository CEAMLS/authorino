package service

import (
	gojson "encoding/json"
	"testing"

	"github.com/kuadrant/authorino/pkg/evaluators"

	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

// A check request as the Kuadrant wasm filter sends it for a TelemetryPolicy
// with labels that read the identity.
const metricLabelsBenchRequest = `{
	"attributes": {
		"metadata_context": {
			"filter_metadata": {
				"io.kuadrant.metrics.labels": {
					"user": { "cel_expr": "auth.identity.userid" },
					"subscription": { "cel_expr": "auth.identity.selected_subscription" },
					"cost_center": { "cel_expr": "has(auth.identity.subscription_info.costCenter) ? auth.identity.subscription_info.costCenter : \"\"" },
					"organization_id": { "cel_expr": "has(auth.identity.subscription_info.organizationId) ? auth.identity.subscription_info.organizationId : \"\"" }
				}
			}
		},
		"request": {
			"http": {
				"host": "maas.example.com",
				"path": "/models/a-model/v1/chat/completions",
				"method": "POST",
				"headers": {
					"authorization": "Bearer sk-oai-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
					"content-type": "application/json",
					"user-agent": "bench/1.0",
					"x-forwarded-for": "10.0.0.1,10.0.0.2",
					"x-request-id": "3e16787c-78e8-9ebb-8756-6fad46e2a0e3"
				}
			}
		}
	}
}`

// One call of metricLabels(). Authorino calls it two or three times for each
// evaluator of the AuthConfig, on every request.
func BenchmarkPipelineMetricLabels(b *testing.B) {
	var request envoy_auth.CheckRequest
	if err := gojson.Unmarshal([]byte(metricLabelsBenchRequest), &request); err != nil {
		b.Fatal(err)
	}
	pipeline := newTestAuthPipeline(
		evaluators.AuthConfig{Labels: map[string]string{"namespace": "authorino", "authconfig": "api-protection"}},
		&request,
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pipeline.metricLabels()
	}
}
