// Package testdata contains data for testing.
package testdata

// https://github.com/alex-leonhardt/k8s-mutate-webhook/blob/master/pkg/mutate/mutate_test.go

// ReviewWithOneImage is a test AdmissionReview used for templating
// a single container image.
const ReviewWithOneImage = `{
	"kind": "AdmissionReview",
	"apiVersion": "admission.k8s.io/v1beta1",
	"request": {
		"uid": "c4e1ae3b-7a3c-496a-99cb-a1f950e32edf",
		"kind": {
			"group": "",
			"version": "v1",
			"kind": "Pod"
		},
		"resource": {
			"group": "",
			"version": "v1",
			"resource": "pods"
		},
		"requestKind": {
			"group": "",
			"version": "v1",
			"kind": "Pod"
		},
		"requestResource": {
			"group": "",
			"version": "v1",
			"resource": "pods"
		},
		"namespace": "default",
		"operation": "CREATE",
		"userInfo": {
			"username": "kubernetes-admin",
			"groups": [
				"system:masters",
				"system:authenticated"
			]
		},
		"object": {
			"kind": "Pod",
			"apiVersion": "v1",
			"metadata": {
				"name": "test-ecr-containers",
				"namespace": "default",
				"creationTimestamp": null,
				"labels": {
					"cloud": "aws"
				}
			},
			"spec": {
				"containers": [{
						"name": "generic",
						"image": "%s",
						"terminationMessagePath": "/dev/termination-log",
						"terminationMessagePolicy": "FallbackToLogsOnError",
						"imagePullPolicy": "IfNotPresent"
					}
				],
				"restartPolicy": "Always",
				"terminationGracePeriodSeconds": 30,
				"dnsPolicy": "ClusterFirst",
				"serviceAccountName": "default",
				"serviceAccount": "default",
				"securityContext": {},
				"schedulerName": "default-scheduler",
				"priority": 0,
				"enableServiceLinks": true
			},
			"status": {}
		},
		"oldObject": null,
		"dryRun": false,
		"options": {
			"kind": "CreateOptions",
			"apiVersion": "meta.k8s.io/v1"
		}
	}
}`
