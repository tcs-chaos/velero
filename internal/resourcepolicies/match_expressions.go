package resourcepolicies

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type MatchExpressions []Requirement

func (m MatchExpressions) Matches(obj unstructured.Unstructured) (bool, error) {
	for _, r := range m {
		match, err := r.Matches(obj)
		if err != nil || !match {
			return false, err
		}
	}
	return true, nil
}

type Operator string

const (
	OperatorIn           Operator = "In"
	OperatorNotIn        Operator = "NotIn"
	OperatorExists       Operator = "Exists"
	OperatorDoesNotExist Operator = "DoesNotExist"
)

type Requirement struct {
	// Label key is the label key that the selector applies to.
	LabelKey string `yaml:"labelKey" protobuf:"bytes,1,opt,name=labelKey"`
	// fieldPath represents a key's path to a field in a structured object.
	FieldPath string `yaml:"fieldPath" protobuf:"bytes,2,opt,name=fieldPath"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator Operator `yaml:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	// +listType=atomic
	Values []string `yaml:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

func (r *Requirement) hasValue(value string) bool {
	for _, v := range r.Values {
		if v == value {
			return true
		}
	}
	return false
}

func (r *Requirement) Matches(obj unstructured.Unstructured) (bool, error) {
	if r.LabelKey == "" && r.FieldPath == "" {
		switch r.Operator {
		// empty operator means Equal
		case "", OperatorExists, OperatorIn:
			return true, nil
		default:
			return false, nil
		}
	}

	if r.LabelKey != "" && r.FieldPath != "" {
		return false, fmt.Errorf("key and fieldPath cannot be set at the same time")
	}

	if r.FieldPath != "" {
		return r.matchFieldPath(obj)
	}
	return r.matchLabel(obj)
}

func (r *Requirement) matchLabel(obj unstructured.Unstructured) (bool, error) {
	ls := Set(obj.GetLabels())

	switch r.Operator {
	case OperatorExists:
		return ls.Has(r.LabelKey), nil
	case OperatorDoesNotExist:
		return !ls.Has(r.LabelKey), nil
	case OperatorIn:
		if !ls.Has(r.LabelKey) {
			return false, nil
		}
		return r.hasValue(ls.Get(r.LabelKey)), nil
	case OperatorNotIn:
		if !ls.Has(r.LabelKey) {
			return true, nil
		}
		return !r.hasValue(ls.Get(r.LabelKey)), nil
	default:
		return false, fmt.Errorf("unknown operator %s", r.Operator)
	}
}

func (r *Requirement) matchFieldPath(obj unstructured.Unstructured) (bool, error) {
	var exist bool
	value, err := ExtractFieldPathAsString(obj, r.FieldPath)
	if err == nil {
		exist = true
	}

	switch r.Operator {
	case OperatorExists, OperatorDoesNotExist:
		return exist, nil
	default:
		if !exist {
			return false, fmt.Errorf("field %s not found in object", r.FieldPath)
		}
	}

	switch r.Operator {
	case OperatorIn:
		return r.hasValue(value), nil
	case OperatorNotIn:
		return !r.hasValue(value), nil
	default:
		return false, fmt.Errorf("unknown operator %s", r.Operator)
	}
}

type Set map[string]string

func (s Set) Has(key string) bool {
	_, exist := s[key]
	return exist
}

func (s Set) Get(key string) string {
	return s[key]
}
