package cluster

import (
	"strings"
)

type OwnerReference struct {
	ApiVersion         string `json:"apiVersion,omitempty"`
	Kind               string `json:"kind,omitempty"`
	Name               string `json:"name,omitempty"`
	Uid                string `json:"uid,omitempty"`
	Controller         bool   `json:"controller,omitempty"`
	BlockOwnerDeletion bool   `json:"blockOwnerDeletion,omitempty"`
}

type Metadata struct {
	Name                       string            `json:"name,omitempty"`
	GenerateName               string            `json:"generateName,omitempty"`
	Namespace                  string            `json:"namespace,omitempty"`
	SelfLink                   string            `json:"selfLink,omitempty"`
	Uid                        string            `json:"uid,omitempty"`
	ResourceVersion            string            `json:"resourceVersion,omitempty"`
	Generation                 int64             `json:"generation,omitempty"`
	CreationTimeStamp          string            `json:"creationTimestamp,omitempty"`
	DeletionTimeStamp          string            `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds int64             `json:"deletionGracePeriodSeconds,omitempty"`
	Labels                     interface{}       `json:"labels,omitempty"`
	Annotation                 interface{}       `json:"annotations,omitempty"`
	OwnerReferences            []*OwnerReference `json:"ownerReferences,omitempty"`
	Finalizers                 []string          `json:"finalizers,omitempty"`
	ClusterName                string            `json:"clusterName,omitempty"`
}

type InvolvedObject struct {
	Kind            string `json:"kind,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
	Name            string `json:"name,omitempty"`
	Uid             string `json:"uid,omitempty"`
	ApiVersion      string `json:"apiVersion,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	FiledPath       string `json:"fieldPath,omitempty"`
}

type Condition struct {
	Type               string `json:"type,omitempty"`
	Status             string `json:"status,omitempty"`
	LastHeartbeatTime  string `json:"lastHeartbeatTime,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitmepty"`
	Reason             string `json:"reason,omitempty"`
	Message            string `json:"message,omitempty"`
}

type MatchExpresion struct {
	Key      string   `json:"key,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Vaules   []string `json:"values,omitempty"`
}

type Selector struct {
	MatchLabels      interface{}       `json:"matchLabels,omitempty"`
	MatchExpressions []*MatchExpresion `json:"matchExpressions,omitempty"`
}

func (meta *Metadata) CheckMeta() error {
	if len(meta.Name) == 0 || len(strings.TrimSpace(meta.Name)) == 0 {
		return ParameterNotNULL("Meta.Name")
	}

	if len(meta.Namespace) == 0 || len(strings.TrimSpace(meta.Namespace)) == 0 {
		return ParameterNotNULL("Meta.Namespace")
	}

	return nil
}

func IsSpace(args string) bool {
	if len(args) == 0 || len(strings.TrimSpace(args)) == 0 {
		return true
	}
	return false
}
