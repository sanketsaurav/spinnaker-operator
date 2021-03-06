package kubernetes

import (
	"context"
	"fmt"
	"github.com/armory/spinnaker-operator/pkg/apis/spinnaker/v1alpha2"
	"github.com/armory/spinnaker-operator/pkg/secrets"
	testing2 "github.com/go-logr/logr/testing"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestMakeClient(t *testing.T) {
	s := `
apiVersion: v1
kind: Config
current-context: test-context
clusters:
- cluster:
    api-version: v1
    server: http://mycluster.com
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
users:
- name: test-user
  user:
    token: test-token
`
	a := &Account{
		Name: "test",
		Auth: &v1alpha2.KubernetesAuth{
			KubeconfigFile: fmt.Sprintf("encryptedFile:noop!%s", s),
		},
	}
	kv := &kubernetesAccountValidator{account: a}
	ctx := secrets.NewContext(context.TODO(), nil, "ns1")
	defer secrets.Cleanup(ctx)
	spinCfg := &v1alpha2.SpinnakerService{}
	c, err := kv.makeClient(ctx, spinCfg, nil)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "http://mycluster.com", c.Host)
}

func TestMakeClientWithFileFromConfig(t *testing.T) {
	y := `
apiVersion: spinnaker.io/v1alpha2
kind: SpinnakerService
metadata:
  name: spinnaker
spec:
  spinnakerConfig:
    files:
      kubecfg: | 
        apiVersion: v1
        kind: Config
        current-context: test-context
        clusters:
        - cluster:
            api-version: v1
            server: http://mycluster.com
          name: test-cluster
        contexts:
        - context:
            cluster: test-cluster
            user: test-user
          name: test-context
        users:
        - name: test-user
          user:
            token: test-token
`
	spinSvc := &v1alpha2.SpinnakerService{}
	if !assert.Nil(t, yaml.Unmarshal([]byte(y), spinSvc)) {
		return
	}
	a := &Account{
		Name: "test",
		Auth: &v1alpha2.KubernetesAuth{
			KubeconfigFile: "kubecfg",
		},
	}
	kv := &kubernetesAccountValidator{account: a}
	c, err := kv.makeClient(context.TODO(), spinSvc, nil)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "http://mycluster.com", c.Host)
}

func TestSettingsTest(t *testing.T) {
	cases := []struct {
		name        string
		account     Account
		errExpected bool
	}{
		{
			"namespaces are ok",
			Account{
				Name: "test",
				Settings: map[string]interface{}{
					"namespaces": []string{"ns1", "ns2"},
				},
			},
			false,
		},
		{
			"omitNamespaces are ok",
			Account{
				Name: "test",
				Settings: map[string]interface{}{
					"omitNamespaces": []string{"ns1", "ns2"},
				},
			},
			false,
		},
		{
			"not omitNamespaces and namespaces at the same time",
			Account{
				Name: "test",
				Settings: map[string]interface{}{
					"omitNamespaces": []string{"ns1", "ns2"},
					"namespaces":     []string{"ns2", "ns3"},
				},
			},
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := &kubernetesAccountValidator{account: &c.account}
			err := v.validateSettings(context.TODO(), testing2.NullLogger{})
			assert.Equal(t, c.errExpected, err != nil)
		})
	}
}
