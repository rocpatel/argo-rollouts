package traefik

import (
	"fmt"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/argoproj/argo-rollouts/utils/diff"
	ingressutil "github.com/argoproj/argo-rollouts/utils/ingress"
	logutil "github.com/argoproj/argo-rollouts/utils/log"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	extensionslisters "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/tools/record"
)

// Type holds this controller type
const Type = "Traefik"

// ReconcilerConfig describes static configuration data for the traefik reconciler
type ReconcilerConfig struct {
	Rollout        *v1alpha1.Rollout
	Client         kubernetes.Interface
	Recorder       record.EventRecorder
	ControllerKind schema.GroupVersionKind
	IngressLister  extensionslisters.IngressLister
}

// Reconciler holds required fields to reconcil Traefik resources
type Reconciler struct {
	cfg ReconcilerConfig
	log *logrus.Entry
}

// NewReconciler returns a reconciler struct that brings canary Ingress into the desired state
func NewReconciler(cfg ReconcilerConfig) *Reconciler {
	return &Reconciler{
		cfg: cfg,
		log: logutil.WithRollout(cfg.Rollout),
	}
}

// Type indicates this reconciler is a Traefik reconciler
func (r *Reconciler) Type() string {
	return Type
}

// Reconcile modifies Traefik Ingress resources to reach desired state
func (r *Reconciler) Reconcile(desiredWeight int32) error {
	rollout := r.cfg.Rollout
	ingressName := r.cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.Traefik.Ingress
	ingress, err := r.cfg.IngressLister.Ingresses(rollout.Namespace).Get(ingressName)
	if err != nil {
		return err
	}
	actionService := r.cfg.Rollout.Spec.Strategy.Canary.StableService
	canaryService := r.cfg.Rollout.Spec.Strategy.Canary.CanaryService
	if r.cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.Traefik.RootService != "" {
		actionService = r.cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.Traefik.RootService
	}

	port := r.cfg.Rollout.Spec.Strategy.Canary.TrafficRouting.Traefik.ServicePort
	if !ingressutil.HasRuleWithService(ingress, actionService) {
		return fmt.Errorf("ingress does not have service `%s` in rules", actionService)
	}

	desiredAnnotations, err := getDesiredAnnotations(ingress, rollout, desiredWeight)
	if err != nil {
		return err
	}
	desiredRules, err := getDesiredRules(ingress, actionService, canaryService, port)
	if err != nil {
		return err
	}

	patch, modified, err := calculatePatch(ingress, desiredAnnotations, desiredRules)
	if err != nil {
		return err
	}
	if !modified {
		r.log.Info("no changes to the Traefik Ingress")
		return nil
	}

	r.log.WithField("patch", string(patch)).Debug("applying Traefik Ingress patch")
	r.log.WithField("desiredWeight", desiredWeight).Info("updating Traefik ingress")
	r.cfg.Recorder.Event(r.cfg.Rollout, corev1.EventTypeNormal, "PatchingTraefikIngress", fmt.Sprintf("Updating Ingress `%s` to desiredWeight '%d'", ingressName, desiredWeight))
	_, err = r.cfg.Client.ExtensionsV1beta1().Ingresses(ingress.Namespace).Patch(ingress.Name, types.MergePatchType, patch)
	if err != nil {
		r.log.WithField("err", err.Error()).Error("error patching traefik ingress")
		return fmt.Errorf("error patching traefik ingress `%s`: %v", ingressName, err)
	}
	return nil
}

func calculatePatch(current *extensionsv1beta1.Ingress, desiredAnnotations map[string]string, desiredRules []extensionsv1beta1.IngressRule) ([]byte, bool, error) {
	//only Compare Annotations and Rules
	return diff.CreateTwoWayMergePatch(
		&extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: current.Annotations,
			},
			Spec: extensionsv1beta1.IngressSpec{
				Rules: current.Spec.Rules,
			},
		},
		&extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: desiredAnnotations,
			},
			Spec: extensionsv1beta1.IngressSpec{
				Rules: desiredRules,
			},
		}, extensionsv1beta1.Ingress{})
}

func getServiceWeightsString(canaryService string, desiredWeight int32) string {
	s := `|
%s: %d%`
	return fmt.Sprintf(s, canaryService, desiredWeight)
}

func getDesiredAnnotations(current *extensionsv1beta1.Ingress, r *v1alpha1.Rollout, desiredWeight int32) (map[string]string, error) {
	canaryService := r.Spec.Strategy.Canary.CanaryService
	desired := current.DeepCopy().Annotations
	key := ingressutil.TraefikServiceWeightsKey(r)
	desired[key] = getServiceWeightsString(canaryService, desiredWeight)
	return desired, nil
}

func getDesiredRules(current *extensionsv1beta1.Ingress, actionService string, canaryService string, port int32) ([]extensionsv1beta1.IngressRule, error) {
	desired := current.DeepCopy().Spec.Rules

	for ruleIndex, rule := range current.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.ServiceName == actionService {
					portStr := intstr.FromInt(int(port))
					canaryPath := path.DeepCopy()
					canaryPath.Backend.ServiceName = canaryService
					canaryPath.Backend.ServicePort = portStr
					desired[ruleIndex].HTTP.Paths = append(rule.HTTP.Paths, *canaryPath)
				}
			}
		}
	}

	return desired, nil
}
