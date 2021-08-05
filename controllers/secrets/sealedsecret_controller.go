/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
	secretsv1beta1 "github.com/rustrial/k8s-gitops-secrets/apis/secrets/v1beta1"
	"github.com/rustrial/k8s-gitops-secrets/internal/providers"
	apiCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubectl/pkg/util/slice"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	fieldManager = "controller.secrets.rustrial.org"
)

// SealedSecretReconciler reconciles a SealedSecret object
type SealedSecretReconciler struct {
	client.Client
	Log                 logr.Logger
	Scheme              *runtime.Scheme
	ControllerNamespace string
	Recorder            record.EventRecorder
}

//+kubebuilder:rbac:groups=secrets.rustrial.org,resources=sealedsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=secrets.rustrial.org,resources=sealedsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=secrets.rustrial.org,resources=sealedsecrets/finalizers,verbs=update
//+kubebuilder:rbac:groups=secrets.rustrial.org,resources=keyencryptionkeypolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=secrets.rustrial.org,resources=keyencryptionkeypolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=secrets.rustrial.org,resources=keyencryptionkeypolicies/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *SealedSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()
	log := r.Log.WithValues("sealedsecret", req.NamespacedName)
	var sealedSecret = &secretsv1beta1.SealedSecret{}
	if err := r.Get(ctx, req.NamespacedName, sealedSecret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !sealedSecret.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, sealedSecret)
	}
	result, err := r.reconcile(ctx, sealedSecret)
	if updateStatusErr := client.IgnoreNotFound(r.patchStatus(ctx, sealedSecret, err)); updateStatusErr != nil {
		log.Error(updateStatusErr, "failed to update status after reconciliation")
		result = ctrl.Result{Requeue: true, RequeueAfter: time.Second * 5}
		if err == nil {
			err = updateStatusErr
		}
	}
	durationMsg := fmt.Sprintf("reconcilation of SealedSecret %s/%s finished in %s (%s)", sealedSecret.Namespace, sealedSecret.Name, time.Now().Sub(start).String(), err)
	if result.RequeueAfter > 0 {
		durationMsg = fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String())
	}
	log.Info(durationMsg)
	return result, nil
}

func (r *SealedSecretReconciler) reconcile(ctx context.Context, sealedSecret *secretsv1beta1.SealedSecret) (ctrl.Result, error) {
	secret, deleteSecret, err := r.decryptSecret(ctx, sealedSecret)
	if deleteSecret {
		return r.reconcileDelete(ctx, sealedSecret)
	} else if err == nil && secret != nil {
		err = controllerutil.SetControllerReference(sealedSecret, secret, r.Scheme)
		if err != nil {
			// This should not happen, as we are setting the controller reference on a fresh
			// in-memory object. However, let's be defensive here.
			r.Log.Error(err, fmt.Sprintf("Failed to set controller reference on Secret %s/%s", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name))
		}
		key := client.ObjectKeyFromObject(secret)
		latest := &apiCoreV1.Secret{}
		if err = r.Client.Get(ctx, key, latest); err != nil {
			if errors.IsNotFound(err) {
				err = r.Create(ctx, secret, &client.CreateOptions{FieldManager: fieldManager})
				if err == nil {
					msg := fmt.Sprintf("Created %s %s/%s", secret.Kind, secret.Namespace, secret.Name)
					r.Log.Info(msg)
					r.Recorder.Event(secret, "Normal", "Created", msg)
				}
			}
		} else {
			// Long-term we want to move-on to server side apply, but we have to introduce
			// it softly making sure we are not breaking any existing installations. Thus
			// we start with an opt-in phase during which JSON MergePatch will stay the default
			// and if there are no problems we might in a later release switch to an opt-out
			// approach where server side apply will become the default.
			patchStrategy := os.Getenv("PATCH_STRATEGY")
			switch patchStrategy {
			case "ServerSideApply":
				po := &client.PatchOptions{FieldManager: fieldManager}
				po.ApplyOptions([]client.PatchOption{client.ForceOwnership})
				err = r.Patch(ctx, secret, client.Apply, po)
				break
			default:
				// Make sure we retain (do not remove) any finalizers added to the Secret by other
				// controllers.
				secret.ObjectMeta.Finalizers = latest.ObjectMeta.Finalizers
				err = r.Patch(ctx, secret, client.MergeFrom(latest), &client.PatchOptions{FieldManager: fieldManager})
				break
			}
			if err == nil {
				msg := fmt.Sprintf("Patched %s %s/%s", secret.Kind, secret.Namespace, secret.Name)
				r.Log.Info(msg)
				r.Recorder.Event(secret, "Normal", "Updated", msg)
			}
		}
	}
	requeueAfter := time.Second * 0
	if err != nil {
		requeueAfter = time.Second * 60
	}
	return ctrl.Result{Requeue: err != nil, RequeueAfter: requeueAfter}, err
}

// DecryptSecret decrypts SealedSecret into a Kubernetes (core) Secret.
func (r *SealedSecretReconciler) decryptSecret(ctx context.Context, sealedSecret *secretsv1beta1.SealedSecret) (*apiCoreV1.Secret, bool, error) {
	name := sealedSecret.Spec.Metadata.Name
	if name == "" {
		name = sealedSecret.ObjectMeta.Name
	}
	secret := &apiCoreV1.Secret{
		TypeMeta: v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: name,
			// While the name can be change, labels and annotations might
			// be addes, we do not allow chosing a different namespace.
			// This is for security reasons.
			Namespace:   sealedSecret.ObjectMeta.Namespace,
			Labels:      sealedSecret.Spec.Metadata.Labels,
			Annotations: sealedSecret.Spec.Metadata.Annotations,
			// Note, OwnershipReference and Finalizer metadata is set in sealedsecret_controller.go (if at all).
		},
		Immutable:  sealedSecret.Spec.Immutable,
		Type:       sealedSecret.Spec.Type,
		Data:       sealedSecret.Spec.Data,
		StringData: sealedSecret.Spec.StringData,
	}
	if secret.Labels == nil {
		// If no labes are explicitly defined, inherit labels from SealedSecret.
		secret.Labels = sealedSecret.Labels
	}
	if secret.Annotations == nil {
		// We do not inherit any annotation from the SealedSecret as they might
		// contain things, which we don't want to inherit.
		secret.Annotations = make(map[string]string)
	}
	secret.Annotations["secrets.rustrial.org/source-type"] = fmt.Sprintf("%s/%s", sealedSecret.APIVersion, sealedSecret.Kind)
	secret.Annotations["secrets.rustrial.org/source"] = fmt.Sprintf("%s/%s", sealedSecret.Namespace, sealedSecret.Name)
	// KEK cache, to speedup things as it is very likely that the same KEK is used for multiple
	// secret fields.
	var kekCache map[string][]secretsv1beta1.KeyEncryptionKeyPolicy = make(map[string][]secretsv1beta1.KeyEncryptionKeyPolicy)
	// We collect all errors and do not abort on the first error. This will help DevOps staff
	// to fix their SealedSecret objects in one go, as they see all the errors at once.
	var errors error
	unauthorizedKeks := make([]string, 0)
	for key, envelopes := range sealedSecret.Spec.EncryptedData {
		for _, value := range envelopes {
			var err error
			keks, err := r.findKEK(ctx, value.KeyEncryptionKeyID, kekCache)
			if err == nil {
				kek, authorizationError := r.kekUsageAllowed(keks, sealedSecret, &value)
				if authorizationError == nil {
					if sealedSecret.Status.Authorizations == nil {
						sealedSecret.Status.Authorizations = make(map[string]secretsv1beta1.Authorization)
					}
					sealedSecret.Status.Authorizations[kek.Spec.KeyEncryptionKeyID] = secretsv1beta1.Authorization{
						Kind:      kek.Kind,
						Name:      kek.Name,
						Namespace: kek.Namespace,
					}
					var provider providers.DecryptionProvider
					provider, err = providers.GetProvider(ctx, &value.Provider)
					if err == nil {
						var binary []byte
						binary, err = provider.Decrypt(ctx, &value)
						if err == nil {
							if secret.Data == nil {
								secret.Data = make(map[string][]byte)
							}
							secret.Data[key] = binary
							if secret.StringData != nil {
								// encrypted secrets have highest prevendence, so make sure we delete any
								// StringData entry with same key.
								delete(secret.StringData, key)
							}
							break
						}
					}
				} else {
					if !slice.ContainsString(unauthorizedKeks, value.KeyEncryptionKeyID, nil) {
						unauthorizedKeks = append(unauthorizedKeks, value.KeyEncryptionKeyID)
					}
				}
			}
			if err != nil {
				// record error
				errors = multierror.Append(errors, fmt.Errorf("Failed to decrypt spec.encryptedData.%s: %w", key, err))
			}
		}
	}
	if len(unauthorizedKeks) == 0 {
		status := v1.ConditionTrue
		if errors != nil {
			// if there were other errors, we don't know whether all KEKs are authorized or not.
			status = v1.ConditionUnknown
		}
		updateCondition(sealedSecret, v1.Condition{
			Type:    "Authorized",
			Status:  status,
			Reason:  "Authorized",
			Message: fmt.Sprintf("All KeyEncryptionKeys referenced in SealedSecret are authorized for use in namespace '%s'", sealedSecret.Namespace),
		})
	} else {
		updateCondition(sealedSecret, v1.Condition{
			Type:    "Authorized",
			Status:  v1.ConditionFalse,
			Reason:  "NotAuthorized",
			Message: fmt.Sprintf("The following KeyEncryptionKeys referenced in SealedSecret are not authorized for use in namespace '%s': %s", sealedSecret.Namespace, strings.Join(unauthorizedKeks, ",")),
		})
	}
	if errors == nil && len(unauthorizedKeks) == 0 {
		updateCondition(sealedSecret, v1.Condition{
			Type:    "Ready",
			Status:  v1.ConditionTrue,
			Reason:  secretsv1beta1.Success,
			Message: "",
		})
		return secret, false, nil
	} else {
		recordError := errors
		if recordError == nil {
			recordError = fmt.Errorf("The following KeyEncryptionKeys referenced in SealedSecret are not authorized for use in namespace '%s': %s", sealedSecret.Namespace, strings.Join(unauthorizedKeks, ","))
		}
		updateCondition(sealedSecret, v1.Condition{
			Type:    "Ready",
			Status:  v1.ConditionFalse,
			Reason:  "Error",
			Message: fmt.Sprintf("%s", recordError),
		})
		if errors != nil {
			r.Log.Error(errors, fmt.Sprintf("Error while processing SealedSecret %s/%s", sealedSecret.Namespace, sealedSecret.Name))
		}
		// Only delete secret if it is no longer authorized.
		// If any other error occurred we do not delete it, to make sure we do not
		// destabilize any workloads that depend on it, while giving DevOps staff
		// a chance to fix the underlying problem.
		return nil, len(unauthorizedKeks) > 0, nil
	}
}

func (r *SealedSecretReconciler) kekUsageAllowed(keks []secretsv1beta1.KeyEncryptionKeyPolicy, secret *secretsv1beta1.SealedSecret, envelope *secretsv1beta1.Envelope) (*secretsv1beta1.KeyEncryptionKeyPolicy, error) {
	if keks == nil {
		return nil, fmt.Errorf("no KeyEncryptionKeyPolicy found in namespace '%s' for KeyEncryptionKey '%s'", r.ControllerNamespace, envelope.KeyEncryptionKeyID)
	}
	for _, kek := range keks {
		if slice.ContainsString(kek.Spec.Namespaces, secret.Namespace, nil) || slice.ContainsString(kek.Spec.Namespaces, "*", nil) {
			return &kek, nil
		}
	}
	return nil, fmt.Errorf("no KeyEncryptionKeyPolicy in namespace '%s' allows usage of KeyEncryptionKey '%s' in namespace '%s'", r.ControllerNamespace, envelope.KeyEncryptionKeyID, secret.Namespace)
}

func (r *SealedSecretReconciler) findKEK(ctx context.Context, keyEncryptionKeyID string, cache map[string][]secretsv1beta1.KeyEncryptionKeyPolicy) ([]secretsv1beta1.KeyEncryptionKeyPolicy, error) {
	keks := cache[keyEncryptionKeyID]
	if keks != nil {
		return keks, nil
	}
	var list secretsv1beta1.KeyEncryptionKeyPolicyList
	result := make([]secretsv1beta1.KeyEncryptionKeyPolicy, 0)
	var options []client.ListOption = make([]client.ListOption, 0)
	options = append(options, client.InNamespace(r.ControllerNamespace), client.MatchingFields{keyEncryptionKeyIDField: keyEncryptionKeyID})
	var err error = nil
	for err = r.Client.List(ctx, &list, options...); err == nil; {
		result = append(result, list.Items...)
		if list.Continue == "" {
			break
		}
		options = make([]client.ListOption, 0)
		options = append(options, client.InNamespace(r.ControllerNamespace), client.MatchingFields{keyEncryptionKeyIDField: keyEncryptionKeyID}, client.Continue(list.Continue))
	}
	cache[keyEncryptionKeyID] = result
	return result, err
}

func updateCondition(secret *secretsv1beta1.SealedSecret, condition v1.Condition) {
	if secret.Status.Conditions == nil {
		secret.Status.Conditions = make([]v1.Condition, 0)
	}
	condition.LastTransitionTime = v1.Time{Time: time.Now()}
	for i, c := range secret.Status.Conditions {
		if c.Type == condition.Type {
			if c.Status != condition.Status || c.Reason != condition.Reason || c.Message != condition.Message {
				secret.Status.Conditions[i] = condition
			}
			return
		}
	}
	secret.Status.Conditions = append(secret.Status.Conditions, condition)
}

func (r *SealedSecretReconciler) patchStatus(ctx context.Context, sealedSecret *secretsv1beta1.SealedSecret, err error) error {
	if err != nil {
		updateCondition(sealedSecret, v1.Condition{
			Type:    "Ready",
			Status:  v1.ConditionFalse,
			Reason:  "Error",
			Message: fmt.Sprintf("%s", err),
		})
	}
	key := client.ObjectKeyFromObject(sealedSecret)
	latest := &secretsv1beta1.SealedSecret{}
	if err := r.Client.Get(ctx, key, latest); err != nil {
		return err
	}
	return r.Client.Status().Patch(ctx, sealedSecret, client.MergeFrom(latest))
}

func (r *SealedSecretReconciler) reconcileDelete(ctx context.Context, sealedSecret *secretsv1beta1.SealedSecret) (ctrl.Result, error) {
	key := client.ObjectKeyFromObject(sealedSecret)
	latest := &apiCoreV1.Secret{}
	var err error = nil
	if err = r.Client.Get(ctx, key, latest); err == nil && latest != nil && latest.ObjectMeta.DeletionTimestamp.IsZero() {
		err = r.Delete(ctx, latest)
		if err == nil {
			msg := fmt.Sprintf("Deleted %s %s/%s", latest.Kind, latest.Namespace, latest.Name)
			r.Log.Info(msg)
			r.Recorder.Event(latest, "Normal", "Deleted", msg)
		}
	}
	return ctrl.Result{}, client.IgnoreNotFound(err)
}

// Select all SealedSecrets for reconciliation which contain a reference to
// the modified KEK.
//
// Unfortunately, I have not yet figured out whether we can index fields in arrays
// using `mgr.GetFieldIndexer().IndexField(...)`, therefore we have to rely on this
// brute force lookup.
func (r *SealedSecretReconciler) findAffectedSecrets(object client.Object) []reconcile.Request {
	ctx := context.TODO()
	var kek secretsv1beta1.KeyEncryptionKeyPolicy
	key := client.ObjectKeyFromObject(object)
	sealedSecrets := make([]reconcile.Request, 0)
	if err := r.Client.Get(ctx, key, &kek); err == nil {
		var list secretsv1beta1.SealedSecretList
		var options []client.ListOption = make([]client.ListOption, 0)
		options = append(options, client.Limit(50))
		for err := r.Client.List(context.TODO(), &list, options...); err == nil; {
			for _, secret := range list.Items {
				for _, envelopes := range secret.Spec.EncryptedData {
					for _, envelope := range envelopes {
						if envelope.KeyEncryptionKeyID == kek.Spec.KeyEncryptionKeyID {
							sealedSecrets = append(sealedSecrets, reconcile.Request{
								NamespacedName: types.NamespacedName{
									Name:      secret.Name,
									Namespace: secret.Namespace,
								},
							})
						}
					}
				}
			}
			if list.Continue == "" {
				break
			}
			options = make([]client.ListOption, 0)
			options = append(options, client.Limit(50), client.Continue(list.Continue))
		}
	}
	return sealedSecrets
}

const (
	keyEncryptionKeyIDField = ".spec.keyEncryptionKeyId"
)

// SetupWithManager sets up the controller with the Manager.
func (r *SealedSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Create index on keyEncryptionKeyId field, which will be used in the findKEK function.
	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &secretsv1beta1.KeyEncryptionKeyPolicy{}, keyEncryptionKeyIDField, func(rawObj client.Object) []string {
		kek := rawObj.(*secretsv1beta1.KeyEncryptionKeyPolicy)
		if kek == nil || kek.Namespace != r.ControllerNamespace {
			return nil
		}
		return []string{kek.Spec.KeyEncryptionKeyID}
	}); err != nil {
		return err
	}
	// Only deal with KEK from the controller's namespace.
	fromControllerNamespace := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectNew.GetNamespace() == r.ControllerNamespace
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return e.Object.GetNamespace() == r.ControllerNamespace
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return e.Object.GetNamespace() == r.ControllerNamespace
		},
	}
	mapper := func(object client.Object) []reconcile.Request {
		return r.findAffectedSecrets(object)
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretsv1beta1.SealedSecret{}).
		Owns(&apiCoreV1.Secret{}).
		Watches(&source.Kind{Type: &secretsv1beta1.KeyEncryptionKeyPolicy{}}, handler.EnqueueRequestsFromMapFunc(mapper), builder.WithPredicates(fromControllerNamespace)).
		Complete(r)
}
