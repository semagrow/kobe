// +build !

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1.KobeBenchmark":       schema_pkg_apis_kobebenchmark_v1alpha1_KobeBenchmark(ref),
		"github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1.KobeBenchmarkSpec":   schema_pkg_apis_kobebenchmark_v1alpha1_KobeBenchmarkSpec(ref),
		"github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1.KobeBenchmarkStatus": schema_pkg_apis_kobebenchmark_v1alpha1_KobeBenchmarkStatus(ref),
	}
}

func schema_pkg_apis_kobebenchmark_v1alpha1_KobeBenchmark(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KobeBenchmark is the Schema for the kobebenchmarks API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1.KobeBenchmarkSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1.KobeBenchmarkStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1.KobeBenchmarkSpec", "github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1.KobeBenchmarkStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_kobebenchmark_v1alpha1_KobeBenchmarkSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KobeBenchmarkSpec defines the desired state of KobeBenchmark",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_kobebenchmark_v1alpha1_KobeBenchmarkStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KobeBenchmarkStatus defines the observed state of KobeBenchmark",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
