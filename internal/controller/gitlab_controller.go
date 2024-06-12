// controllers/gitlab_controller.go

package controllers

import (
    "context"
    "fmt"
    "time"

    "github.com/go-logr/logr"
    "gitlab-operator/api/v1alpha1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
    "sigs.k8s.io/controller-runtime/pkg/log"
)

// GitlabReconciler reconciles a Gitlab object
type GitlabReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.example.com,resources=gitlabs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.example.com,resources=gitlabs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.example.com,resources=gitlabs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *GitlabReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    _ = log.FromContext(ctx)
    log := log.Log.WithValues("gitlab", req.NamespacedName)

    // Fetch the Gitlab instance
    gitlab := &v1alpha1.Gitlab{}
    err := r.Get(ctx, req.NamespacedName, gitlab)
    if err != nil {
        if errors.IsNotFound(err) {
            // Object not found, return. Created objects are automatically garbage collected.
            // For additional cleanup logic use finalizers.
            return ctrl.Result{}, nil
        }
        // Error reading the object - requeue the request.
        return ctrl.Result{}, err
    }

    // Check if the Gitlab project already exists, if not create it
    if gitlab.Status.Conditions == nil {
        gitlab.Status.Conditions = []metav1.Condition{}
    }

    if !r.isGitlabProjectCreated(gitlab) {
        err = r.createGitlabProject(ctx, gitlab)
        if err != nil {
            log.Error(err, "Failed to create Gitlab project")
            return ctrl.Result{}, err
        }

        // Update the status to indicate the project was created
        gitlab.Status.Conditions = append(gitlab.Status.Conditions, metav1.Condition{
            Type:               "ProjectCreated",
            Status:             metav1.ConditionTrue,
            LastTransitionTime: metav1.Now(),
            Reason:             "ProjectCreationSuccessful",
            Message:            fmt.Sprintf("Project %s created successfully", gitlab.Spec.ProjectName),
        })

        err = r.Status().Update(ctx, gitlab)
        if err != nil {
            log.Error(err, "Failed to update Gitlab status")
            return ctrl.Result{}, err
        }
    }

    // Requeue the request to ensure the status is updated periodically
    return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *GitlabReconciler) isGitlabProjectCreated(gitlab *v1alpha1.Gitlab) bool {
    for _, condition := range gitlab.Status.Conditions {
        if condition.Type == "ProjectCreated" && condition.Status == metav1.ConditionTrue {
            return true
        }
    }
    return false
}

func (r *GitlabReconciler) createGitlabProject(ctx context.Context, gitlab *v1alpha1.Gitlab) error {
    // Placeholder for the logic to create a project on Gitlab.
    // You can use the Gitlab API client to create a project and repository.
    // Example:
    // gitlabClient := gitlab.NewClient("your-token")
    // project, err := gitlabClient.CreateProject(gitlab.Spec.ProjectName)
    // if err != nil {
    //     return err
    // }

    // For simplicity, let's assume the project creation is always successful.
    return nil
}

func (r *GitlabReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&v1alpha1.Gitlab{}).
        Complete(r)
}

