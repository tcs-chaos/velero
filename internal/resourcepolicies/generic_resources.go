package resourcepolicies

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GenericPolicy struct {
	Condition GenericPolicyCondition `yaml:"condition"`
	Action    Action                 `yaml:"action"`
}

type GenericPolicyCondition struct {
	Resource         string           `yaml:"resource"` // regex expression
	MatchExpressions MatchExpressions `yaml:"matchExpressions"`
}

func (p *GenericPolicyCondition) Match(obj runtime.Unstructured, resource schema.GroupResource) (bool, error) {
	unstructuredObj := unstructured.Unstructured{Object: obj.UnstructuredContent()}

	// TODO: just consider the resource kind for now
	// resource can be core/v1/pods or core/v1/pods/* or v1/pods
	var matched bool
	var err error
	// if resource pattern is not specific, just match expression.
	if len(p.Resource) == 0 || p.Resource == "*" {
		goto expression
	}
	matched = p.Resource == resource.String()
	if !matched {
		return false, err
	}

expression:
	for _, expr := range p.MatchExpressions {
		matched, err = expr.Matches(unstructuredObj)
		if !matched || err != nil {
			return false, err
		}
	}
	return true, nil
}

type GenericPolicyList []GenericPolicy

func (pl GenericPolicyList) Match(obj runtime.Unstructured, resource schema.GroupResource) (*Action, error) {
	for _, p := range pl {
		if ok, err := p.Condition.Match(obj, resource); ok || err != nil {
			return &p.Action, err
		}
	}
	return nil, nil
}
