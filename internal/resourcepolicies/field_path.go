package resourcepolicies

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ExtractFieldPathAsString extracts the field from the given object and returns it as a string.
func ExtractFieldPathAsString(obj unstructured.Unstructured, fieldPath string) (string, error) {
	value, err := extractFieldPathAsString(obj.Object, fieldPath)
	if err != nil {
		return "", fmt.Errorf("failed to extract field %s as string: %v", fieldPath, err)
	}
	return value, nil
}

func extractFieldPathAsString(obj map[string]any, fieldPath string) (string, error) {
	parts := strings.Split(fieldPath, ".")
	for i, part := range parts {
		var (
			value any
			ok    bool
			str   string
		)
		if value, ok = obj[part]; !ok {
			return "", fmt.Errorf("field %s not found in object", part)
		}
		if i == len(parts)-1 {
			if str, ok = value.(string); !ok {
				return "", fmt.Errorf("field %s is not a string", part)
			}
			return str, nil
		}
		if obj, ok = value.(map[string]any); !ok {
			return "", fmt.Errorf("field %s is not a map", part)
		}
	}
	return "", fmt.Errorf("field path %s is empty", fieldPath)
}
