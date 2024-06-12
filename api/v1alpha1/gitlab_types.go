// api/v1alpha1/gitlab_types.go

package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GitlabSpec defines the desired state of Gitlab
type GitlabSpec struct {
    ProjectName string `json:"projectName"`
    Repository  string `json:"repository"`
}

// GitlabStatus defines the observed state of Gitlab
type GitlabStatus struct {
    // Insert additional status fields - define observed state of cluster
    // Important: Run "make" to regenerate code after modifying this file
    Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gitlab is the Schema for the gitlabs API
type Gitlab struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   GitlabSpec   `json:"spec,omitempty"`
    Status GitlabStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GitlabList contains a list of Gitlab
type GitlabList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Gitlab `json:"items"`
}

func init() {
    SchemeBuilder.Register(&Gitlab{}, &GitlabList{})
}

