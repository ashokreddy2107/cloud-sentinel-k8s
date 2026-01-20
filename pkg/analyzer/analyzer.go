package analyzer

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Analyzer interface {
	Name() string
	Analyze(ctx context.Context, client client.Client, obj client.Object) ([]Anomaly, error)
}
