package transformer

import (
	"context"
	"testing"

	spinnakerv1alpha2 "github.com/armory/spinnaker-operator/pkg/apis/spinnaker/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
)

func TestSetTarget(t *testing.T) {
	svc := spinnakerv1alpha2.SpinnakerService{
		ObjectMeta: metav1.ObjectMeta{Namespace: "ns2"},
		Spec: spinnakerv1alpha2.SpinnakerServiceSpec{
			SpinnakerConfig: spinnakerv1alpha2.SpinnakerConfig{
				Config: spinnakerv1alpha2.FreeForm{
					"deploymentEnvironment": map[string]string{
						"location": "ns1",
					},
				},
			},
		},
	}
	tg := &targetTransformer{svc: &svc}
	ctx := context.TODO()
	err := tg.TransformConfig(ctx)
	if assert.Nil(t, err) {
		s, err := svc.GetSpinnakerConfig().GetHalConfigPropString(ctx, "deploymentEnvironment.location")
		assert.Nil(t, err)
		assert.Equal(t, "ns2", s)
	}
}
