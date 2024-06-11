package resourcepolicies

import (
	"fmt"

	"regexp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type GenericPolicy struct {
	Condition GenericPolicyCondition `yaml:"conditions"`
	Action    Action                 `yaml:"action"`
}

type GenericPolicyCondition struct {
	Resource         string           `yaml:"resource"` // regex expression
	MatchExpressions MatchExpressions `yaml:"matchExpressions"`
}

func (p *GenericPolicyCondition) Match(obj interface{}) (bool, error) {
	objc, ok := obj.(runtime.Unstructured)
	if !ok {
		return false, fmt.Errorf("failed to convert object")
	}

	unstructuredObj := unstructured.Unstructured{Object: objc.UnstructuredContent()}
	gvk := unstructuredObj.GroupVersionKind()

	// TODO: just consider the resource kind for now
	// resource can be core/v1/pods or core/v1/pods/* or v1/pods
	var matched bool
	var err error
	// if resource pattern is not specific, just match expression.
	if len(p.Resource) == 0 {
		goto expression
	}
	matched, err = regexp.Match(p.Resource, []byte(gvk.Kind))
	if !matched || err != nil {
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

func (pl GenericPolicyList) Match(obj interface{}) (*Action, error) {
	for _, p := range pl {
		if ok, err := p.Condition.Match(obj); ok || err != nil {
			return &p.Action, err
		}
	}
	return nil, nil
}
