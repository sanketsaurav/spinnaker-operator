package test

import (
	"fmt"
	"github.com/armory/spinnaker-operator/pkg/apis/spinnaker/v1alpha2"
	"github.com/armory/spinnaker-operator/pkg/generated"
	"io/ioutil"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"
	"testing"
)

func ManifestToSpinService(manifestYaml string, t *testing.T) *v1alpha2.SpinnakerService {
	svc := &v1alpha2.SpinnakerService{}
	ReadYamlFile(manifestYaml, svc, t)
	return svc
}

func ReadYamlFile(path string, target interface{}, t *testing.T) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(bytes, target)
	if err != nil {
		t.Fatal(err)
	}
}

func FakeClient(t *testing.T, objs ...runtime.Object) client.Client {
	return fake.NewFakeClientWithScheme(scheme.Scheme, objs...)
}

func BuildSvc(name string, svcType string, publicPort int32, t *testing.T) *corev1.Service {
	svc := &corev1.Service{}
	ReadYamlFile("testdata/service.yml", svc, t)
	svc.Name = name
	svc.Spec.Selector["cluster"] = name
	svc.Spec.Type = corev1.ServiceType(svcType)
	svc.Spec.Ports[0].Port = publicPort
	portName := fmt.Sprintf("%s-tcp", name[len("spin-"):])
	svc.Spec.Ports[0].Name = portName
	return svc
}

func AddDeploymentToGenConfig(gen *generated.SpinnakerGeneratedConfig, depName string, fileName string, t *testing.T) {
	if gen.Config == nil {
		gen.Config = make(map[string]generated.ServiceConfig)
	}
	dep := &v1.Deployment{}
	ReadYamlFile(fileName, dep, t)
	gen.Config[depName] = generated.ServiceConfig{
		Deployment: dep,
	}
}

func AddServiceToGenConfig(gen *generated.SpinnakerGeneratedConfig, svcName string, fileName string, t *testing.T) {
	if gen.Config == nil {
		gen.Config = make(map[string]generated.ServiceConfig)
	}
	svc := &corev1.Service{}
	ReadYamlFile(fileName, svc, t)
	gen.Config[svcName] = generated.ServiceConfig{
		Service: svc,
	}
}
