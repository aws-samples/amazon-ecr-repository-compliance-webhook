// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"reflect"
	"testing"

	"github.com/swoldemi/amazon-ecr-repository-compliance-webhook/testdata"
	corev1 "k8s.io/api/core/v1"
)

func TestParseRepositories(t *testing.T) {
	var (
		untaggedImagePod   = newPodWithImage(testdata.UntaggedImage)
		taggedImagePod     = newPodWithImage(testdata.TaggedImage)
		cnImagePod         = newPodWithImage(testdata.CNImage)
		fipsImagePod       = newPodWithImage(testdata.FIPSImage)
		duplicateImagesPod = newPodWithImage(testdata.TaggedImage)
		twoImagesPod       = newPodWithImage(testdata.TaggedImage)
		noNamespacePod     = newPodWithImage(testdata.NoNamespace)
		aliasedImagePod    = newPodWithImage(testdata.AliasedImage)
		noImages           = newPodWithImage("")
		badImage           = newPodWithImage("elgoog/sselortsid")
	)
	duplicateImagesPod.Spec.Containers = append(duplicateImagesPod.Spec.Containers, duplicateImagesPod.Spec.Containers...)
	twoImagesPod.Spec.Containers = append(twoImagesPod.Spec.Containers, untaggedImagePod.Spec.Containers...)
	tests := []struct {
		name string
		pod  *corev1.Pod
		want []string
	}{
		{"UntaggedImage", untaggedImagePod, []string{"namespace/repo@sha256:e5e2a3236e64483c50dd2811e46e9cd49c67e82271e60d112ca69a075fc23005"}},
		{"TaggedImage", taggedImagePod, []string{"namespace/repo:40d6072"}},
		{"CNImage", cnImagePod, []string{"namespace/repo:40d6072"}},
		{"FIPSImage", fipsImagePod, []string{"namespace/repo:40d6072"}},
		{"Duplicates", duplicateImagesPod, []string{"namespace/repo:40d6072"}},
		{"TwoImages", twoImagesPod, []string{"namespace/repo:40d6072", "namespace/repo@sha256:e5e2a3236e64483c50dd2811e46e9cd49c67e82271e60d112ca69a075fc23005"}},
		{"NoNamespace", noNamespacePod, []string{"repo:40d6072"}},
		{"Aliased", aliasedImagePod, []string{"namespace/repo:40d6072"}},
		{"NoImages", noImages, nil},
		{"BadImage", badImage, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseImages(tt.pod); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRepositories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newPodWithImage(image string) *corev1.Pod {
	return &corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image: image,
				},
			},
		},
	}
}
