package validation

import (
	"errors"
	"reflect"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	kle "knative.dev/pkg/leaderelection"
)

func okConfig() *kle.Config {
	return &kle.Config{
		ResourceLock:      "leases",
		LeaseDuration:     15 * time.Second,
		RenewDeadline:     10 * time.Second,
		RetryPeriod:       2 * time.Second,
		EnabledComponents: sets.NewString("controller"),
	}
}

func okData() map[string]string {
	return map[string]string{
		"resourceLock": "leases",
		// values in this data come from the defaults suggested in the
		// code:
		// https://github.com/kubernetes/client-go/blob/kubernetes-1.16.0/tools/leaderelection/leaderelection.go
		"leaseDuration":     "15s",
		"renewDeadline":     "10s",
		"retryPeriod":       "2s",
		"enabledComponents": "controller",
	}
}

func TestValidateLeaderElectionConfig(t *testing.T) {
	cases := []struct {
		name     string
		data     map[string]string
		expected *kle.Config
		err      error
	}{
		{
			name:     "OK",
			data:     okData(),
			expected: okConfig(),
		},
		{
			name: "invalid component",
			data: func() map[string]string {
				data := okData()
				data["enabledComponents"] = "controller,frobulator"
				return data
			}(),
			err: errors.New(`invalid enabledComponent "frobulator": valid values are ["certcontroller" "controller" "hpaautoscaler" "istiocontroller" "nscontroller"]`),
		},
	}

	for i := range cases {
		tc := cases[i]
		actualConfig, actualErr := ValidateLeaderElectionConfig(&corev1.ConfigMap{Data: tc.data})
		if !reflect.DeepEqual(tc.err, actualErr) {
			t.Errorf("%v: expected error %v, got %v", tc.name, tc.err, actualErr)
			continue
		}

		if !reflect.DeepEqual(tc.expected, actualConfig) {
			t.Errorf("%v: expected config:\n%+v\ngot:\n%+v", tc.name, tc.expected, actualConfig)
			continue
		}
	}

}
