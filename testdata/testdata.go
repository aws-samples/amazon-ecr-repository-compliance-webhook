// Package testdata contains data for testing.
package testdata

// ReviewWithOneImage is a test AdmissionReview used for templating
// a single container image.
// Referenced: https://github.com/alex-leonhardt/k8s-mutate-webhook/blob/master/pkg/mutate/mutate_test.go
const ReviewWithOneImage = `{
	"kind": "AdmissionReview",
	"apiVersion": "admission.k8s.io/v1beta1",
	"request": {
		"uid": "e77141b6-6033-11ea-8d6a-0ac25c990f4a",
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
		"namespace": "echo-namespace",
		"operation": "CREATE",
		"userInfo": {
			"username": "system:serviceaccount:kube-system:replicaset-controller",
			"uid": "064184fd-5e6c-11ea-8d6a-0ac25c990f4a",
			"groups": ["system:serviceaccounts", "system:serviceaccounts:kube-system", "system:authenticated"]
		},
		"object": {
			"kind": "Pod",
			"apiVersion": "v1",
			"metadata": {
				"name": "echo-68f4474876-dzsjm",
				"generateName": "echo-68f4474876-",
				"namespace": "echo-namespace",
				"uid": "e77136cf-6033-11ea-8d6a-0ac25c990f4a",
				"creationTimestamp": "2020-03-07T05:24:23Z",
				"labels": {
					"app": "echo",
					"pod-template-hash": "68f4474876"
				},
				"annotations": {
					"kubernetes.io/psp": "eks.privileged"
				},
				"ownerReferences": [{
					"apiVersion": "apps/v1",
					"kind": "ReplicaSet",
					"name": "echo-68f4474876",
					"uid": "21494a34-6033-11ea-8d6a-0ac25c990f4a",
					"controller": true,
					"blockOwnerDeletion": true
				}]
			},
			"spec": {
				"volumes": [{
					"name": "default-token-qrv2v",
					"secret": {
						"secretName": "default-token-qrv2v"
					}
				}],
				"containers": [{
					"name": "echo",
					"image": "%s",
					"ports": [{
						"containerPort": 80,
						"protocol": "TCP"
					}],
					"resources": {},
					"volumeMounts": [{
						"name": "default-token-qrv2v",
						"readOnly": true,
						"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
					}],
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File",
					"imagePullPolicy": "IfNotPresent"
				}],
				"restartPolicy": "Always",
				"terminationGracePeriodSeconds": 30,
				"dnsPolicy": "ClusterFirst",
				"serviceAccountName": "default",
				"serviceAccount": "default",
				"securityContext": {},
				"schedulerName": "default-scheduler",
				"tolerations": [{
					"key": "node.kubernetes.io/not-ready",
					"operator": "Exists",
					"effect": "NoExecute",
					"tolerationSeconds": 300
				}, {
					"key": "node.kubernetes.io/unreachable",
					"operator": "Exists",
					"effect": "NoExecute",
					"tolerationSeconds": 300
				}],
				"priority": 0,
				"enableServiceLinks": true
			},
			"status": {
				"phase": "Pending",
				"qosClass": "BestEffort"
			}
		},
		"oldObject": null,
		"dryRun": false
	}
}`

// ReviewWithBadRequest is an AdmissionReview with a bad request.
// TODO: Populate so test doesn't just fail at decoding step.
const ReviewWithBadRequest = `{
	"kind": "AdmissionReview",
	"apiVersion": "admission.k8s.io/v1beta1",
	"request": {
		"uid": "e77141b6-6033-11ea-8d6a-0ac25c990f4a"
	}
}`

// ReviewWithNoUID is an AdmissionReview with a bad request and no UID.
const ReviewWithNoUID = `{
	"kind": "AdmissionReview",
	"apiVersion": "admission.k8s.io/v1beta1",
	"request": {
		"uid": ""
	}
}`

// ECR container image URIs.
// Valid references here: https://docs.aws.amazon.com/general/latest/gr/ecr.html
// Repositories in China use a .com.cn API endpoint, but image references keep .com: https://docs.amazonaws.cn/general/latest/gr/ecr.html
const (
	UntaggedImage = "273450712882.dkr.ecr.us-east-2.amazonaws.com/namespace/repo@sha256:e5e2a3236e64483c50dd2811e46e9cd49c67e82271e60d112ca69a075fc23005"
	TaggedImage   = "273450712882.dkr.ecr.us-east-2.amazonaws.com/namespace/repo:40d6072"
	NoNamespace   = "273450712882.dkr.ecr.us-east-2.amazonaws.com/repo:40d6072"
	CNImage       = "273450712882.dkr.ecr.cn-north-1.amazonaws.com/namespace/repo:40d6072"
	FIPSImage     = "273450712882.dkr.ecr-fips.us-east-2.amazonaws.com/namespace/repo:40d6072"
	AliasedImage  = "myaccountalias.dkr.ecr-fips.us-east-2.amazonaws.com/namespace/repo:40d6072"
)
