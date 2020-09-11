package traefik

import (
	"fmt"
	"testing"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	ingressutil "github.com/argoproj/argo-rollouts/utils/ingress"
	"github.com/stretchr/testify/assert"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/record"
)

func fakeRollout(stableSvc, canarySvc, stableIng string, port int32) *v1alpha1.Rollout {
	return &v1alpha1.Rollout{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rollout",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: v1alpha1.RolloutSpec{
			Strategy: v1alpha1.RolloutStrategy{
				Canary: &v1alpha1.CanaryStrategy{
					StableService: stableSvc,
					CanaryService: canarySvc,
					TrafficRouting: &v1alpha1.RolloutTrafficRouting{
						Traefik: &v1alpha1.TraefikTrafficRouting{
							Ingress:     stableIng,
							ServicePort: port,
						},
					},
				},
			},
		},
	}
}

func traefikWeightAnnotation() string {
	return fmt.Sprintf("%s%s", ingressutil.TraefikIngressAnnotation, ingressutil.TraefikServiceWeightsPrefix)
}

func ingress(name, stableSvc, canarySvc string, port, weight int32, managedBy string) *extensionsv1beta1.Ingress {
	return &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
			Annotations: map[string]string{
				traefikWeightAnnotation(): getServiceWeightsString(canarySvc, weight),
			},
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: stableSvc,
										ServicePort: intstr.FromInt(int(port)),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestType(t *testing.T) {
	client := fake.NewSimpleClientset()
	rollout := fakeRollout("stable-service", "canary-service", "stable-ingress", 443)
	r := NewReconciler(ReconcilerConfig{
		Rollout:        rollout,
		Client:         client,
		Recorder:       &record.FakeRecorder{},
		ControllerKind: schema.GroupVersionKind{Group: "foo", Version: "v1", Kind: "Bar"},
	})
	assert.Equal(t, Type, r.Type())
}

func TestIngressNotFound(t *testing.T) {
	ro := fakeRollout("stable-service", "canary-service", "stable-ingress", 443)
	client := fake.NewSimpleClientset()
	k8sI := kubeinformers.NewSharedInformerFactory(client, 0)
	r := NewReconciler(ReconcilerConfig{
		Rollout:        ro,
		Client:         client,
		Recorder:       &record.FakeRecorder{},
		ControllerKind: schema.GroupVersionKind{Group: "foo", Version: "v1", Kind: "Bar"},
		IngressLister:  k8sI.Extensions().V1beta1().Ingresses().Lister(),
	})
	err := r.Reconcile(10)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestServiceNotFoundInIngress(t *testing.T) {
	ro := fakeRollout("stable-service", "canary-service", "ingress", 443)
	ro.Spec.Strategy.Canary.TrafficRouting.Traefik.RootService = "invalid-svc"
	i := ingress("ingress", "stable-service", "canary-svc", 443, 50, ro.Name)
	client := fake.NewSimpleClientset()
	k8sI := kubeinformers.NewSharedInformerFactory(client, 0)
	k8sI.Extensions().V1beta1().Ingresses().Informer().GetIndexer().Add(i)
	r := NewReconciler(ReconcilerConfig{
		Rollout:        ro,
		Client:         client,
		Recorder:       &record.FakeRecorder{},
		ControllerKind: schema.GroupVersionKind{Group: "foo", Version: "v1", Kind: "Bar"},
		IngressLister:  k8sI.Extensions().V1beta1().Ingresses().Lister(),
	})
	err := r.Reconcile(10)
	assert.Errorf(t, err, "ingress does not use the stable service")
}

func TestNoChanges(t *testing.T) {
	ro := fakeRollout("stable-svc", "canary-svc", "ingress", 443)
	i := ingress("ingress", "stable-svc", "canary-svc", 443, 10, ro.Name)
	client := fake.NewSimpleClientset()
	k8sI := kubeinformers.NewSharedInformerFactory(client, 0)
	k8sI.Extensions().V1beta1().Ingresses().Informer().GetIndexer().Add(i)
	r := NewReconciler(ReconcilerConfig{
		Rollout:        ro,
		Client:         client,
		Recorder:       &record.FakeRecorder{},
		ControllerKind: schema.GroupVersionKind{Group: "foo", Version: "v2", Kind: "Bar"},
		IngressLister:  k8sI.Extensions().V1beta1().Ingresses().Lister(),
	})
	err := r.Reconcile(10)
	assert.Nil(t, err)
	assert.Len(t, client.Actions(), 0)
}
