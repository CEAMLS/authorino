package service

import "github.com/kuadrant/authorino/pkg/expressions/cel"

// The metric label expressions come with each request, in the filter metadata.
// They are the same few strings for every request of a gateway, so compile
// each one once and keep it.
var metricLabelExpressions = cel.NewExpressionCache(256)
