// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package webhook contains resources for the ValidatingWebhookConfiguration.
// Referenced: https://github.com/kubernetes/kubernetes/blob/v1.15.0/test/images/webhook
package webhook

import (
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func init() {
	// Permissive plugin-or-panic for adding runtime schemas
	utilruntime.Must(corev1.AddToScheme(runtimeScheme))
	utilruntime.Must(admissionv1beta1.AddToScheme(runtimeScheme))
	utilruntime.Must(admissionregistrationv1beta1.AddToScheme(runtimeScheme))
}
