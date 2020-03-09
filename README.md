![](https://codebuild.us-east-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiU1NyMHI4KytFRzhZSUVEY2R0YTlwanBJTk9EdWNYbW93TzdRU3NCbUJ0TFZYMy9jUktROXlUQktEOUVjd0dJSDBWbHNtVjVqSFpaNWxvbTJxd0o4dW53PSIsIml2UGFyYW1ldGVyU3BlYyI6ImgyNlBtRXoyU1ZSNjNWZjYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![][sar-logo]](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:273450712882:applications~amazon-ecr-repository-compliance-webhook)

[sar-deploy]: https://img.shields.io/badge/Serverless%20Application%20Repository-Deploy%20Now-FF9900?logo=amazon%20aws&style=flat-square
[sar-logo]: https://img.shields.io/badge/Serverless%20Application%20Repository-View-FF9900?logo=amazon%20aws&style=flat-square

# Amazon ECR Repository Compliance Webhook for Kubernetes
>A Kubernetes ValidatingWebhookConfiguration and serverless backend: Deny Pods with container images that don't come from Amazon ECR, don't enforce image tag immutability, or don't enforce scanning on push

This AWS Serverless Application Repository app will create an Amazon API Gateway and an AWS Lambda Function that act as the backend for a Kubernetes ValidatingWebhookConfiguration. The function will deny Pods that create containers using images which come from ECR repositories that:
1. Do not have image tag immutability enabled
2. Do not have image scan on push enabled

Additionally, If the images do not come from ECR at all, they will be also be **denied from running in the cluster**.

![architecture](https://raw.githubusercontent.com/swoldemi/amazon-ecr-repository-compliance-webhook/master/screenshots/architecture.png)

## Usage
To use this SAR application you will:
1. Deploy the application
2. Configure and deploy the `ValidatingWebhookConfiguration` resource into your Kubernetes cluster (EKS or otherwise). The cluster must have this plugin enabled and be have support for the admissionregistration.k8s.io/v1beta1 API. See the official Kubernetes documentation [here](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) for details.

### 1. Deploying the Lambda
It is recommended that you deploy this Lambda function directly from the AWS Serverless Application Repository. It is also possible to deploy this function using:
- The [SAM CLI](https://aws.amazon.com/serverless/sam/)
- CloudFormation via the [AWS CLI](https://aws.amazon.com/cli/)
- CloudFormation via the [CloudFormation management console](https://aws.amazon.com/cloudformation/)

To deploy this function from AWS GovCloud or regions in China, you must have an account with access to these regions. This function is available in all regions that support API Gateway, AWS Lambda, and Amazon ECR. The function should be deployed in the same region as the ECR repository. If the table below is missing a region, please open a pull request!


|Region                                        |Click and Deploy                                                                                                                                 |
|----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
|**US East (Ohio) (us-east-2)**                |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-east-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**US East (N. Virginia) (us-east-1)**         |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**US West (N. California) (us-west-1)**       |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**US West (Oregon) (us-west-2)**              |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-west-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**Asia Pacific (Hong Kong) (ap-east-1)**      |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**Asia Pacific (Mumbai) (ap-south-1)**        |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-south-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)    |
|**Asia Pacific (Seoul) (ap-northeast-2)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-northeast-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)|
|**Asia Pacific (Singapore)	(ap-southeast-1)** |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-southeast-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)|
|**Asia Pacific (Sydney) (ap-southeast-2)**    |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-southeast-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)|
|**Asia Pacific (Tokyo) (ap-northeast-1)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/ap-northeast-1?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook) |
|**Canada (Central)	(ca-central-1)**           |[![][sar-deploy]](https://deploy.serverlessrepo.app/ca-central-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)  |
|**EU (Frankfurt) (eu-central-1)**             |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-central-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)  |
|**EU (Ireland)	(eu-west-1)**                  |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**EU (London) (eu-west-2)**                   |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-2/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**EU (Paris) (eu-west-3)**                    |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-west-3/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**EU (Stockholm) (eu-north-1)**               |[![][sar-deploy]](https://deploy.serverlessrepo.app/eu-north-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)    |
|**Middle East (Bahrain) (me-south-1)**        |[![][sar-deploy]](https://deploy.serverlessrepo.app/me-south-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)    |
|**South America (Sao Paulo) (sa-east-1)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/sa-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |
|**AWS GovCloud (US-East) (us-gov-east-1)**    |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-gov-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook) |
|**AWS GovCloud (US-West) (us-gov-west-1)**    |[![][sar-deploy]](https://deploy.serverlessrepo.app/us-gov-west-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook) |

#### Parameters
|Name           |Default   |Description                                                       |Required |                 
|---------------|----------|------------------------------------------------------------------|---------|
|RegistryRegion |us-east-1 |What AWS region should this Lambda function interact with ECR in? |False    |

### 2. Configuration
After deploying the SAR application from the SAR console you need to:
1. Authenticate with your cluster. For example, for EKS you can use the AWS CLI: `aws eks update-kubeconfig --name your-clusters-name --region your-clusters-region`
2. `kubectl apply -f validatingwebhook.yaml` (provided [here](./validatingwebhook.yaml)) to deploy the `ValidatingWebhookConfiguration`. Remember to update `webhooks.clientConfig.url` with your API Gateway endpoint. Make any necessary additions to match namespaces/labels for resources that are deployed. This webhook only validates `Pod`s.
2. `kubectl create ns test-namespace && kubectl apply -f mydeployment.yaml` provided [here](./mydeployment.yaml) to create a sample `Deployment`. Change the image to be whatever image you would like to test. Ensure your nodes have permission to pull from the ECR repository.
3. `kubectl get ev -n test-namespace` to see if there are any `FailedCreate` events as a result of the `Deployment`'s `ReplicaSet` triggering a failure from the `ValidatingWebhookConfiguration` when trying to create Pods. For example: `Error creating: admission webhook "ecrpolicies.amazonaws.com" denied the request: webhook: no ecr images found in pod specification`

## Contributing
Have an idea for a feature to enhance this serverless application? Open an [issue](https://github.com/swoldemi/amazon-ecr-repository-compliance-webhook/issues) or [pull request](https://github.com/swoldemi/amazon-ecr-repository-compliance-webhook/pulls)!

### Development
This application has been developed, built, and tested against [Go 1.14](https://golang.org/dl/), the latest version of the [Serverless Application Model CLI](https://github.com/awslabs/aws-sam-cli), and the latest version of the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html), Kubernetes version 1.14, Kubernetes version 1.15, and [kubectl 1.17](https://kubernetes.io/docs/tasks/tools/install-kubectl/). A [Makefile](./Makefile) has been provided for convenience.

```
make check # Run Go linting tools
make test # Run Go tests
make build # Build Go binary
make sam-package # Package code and assets into S3 using SAM CLI
make sam-deploy # Deploy application using SAM CLI
make sam-tail-logs # Tail the logs of the running Lambda function
make destroy # Destroy the CloudFormation stack tied to the app
```

### To Do
1. [Parameter.String] RegistryID - What registry should this Lambda verify container images for?
2. [Parameter.CommaDelimitedList] IgnoredNamespaces - What namespaces should be ignored? It is also possible to set matchers on the [`ValidatingWebhookConfiguration`](./validatingwebhook.yaml#L18).
3. [Authenticate the apiserver](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#authenticate-apiservers)
4. Emit metric on deny/pass, to CloudWatch
5. Move to the admissionregistration.k8s.io/v1 API when EKS supports k8s v1.17 and drops v1.14

## References
- BanzaiCloud - In-depth introduction to Kubernetes admission webhooks: https://banzaicloud.com/blog/k8s-admission-webhooks/
- ValidatingWebhookConfiguration API Documentation - https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#validatingwebhookconfiguration-v1beta1-admissionregistration-k8s-io
- Dynamic Admission Control - https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
- Official Kubernetes example: https://github.com/kubernetes/kubernetes/blob/v1.15.0/test/images/webhook/scheme.go

## Acknowledgements
[@jicowan](https://github.com/jicowan) for inspiration: https://github.com/jicowan/ecr-validation-webhook

## License
[Apache License 2.0](https://spdx.org/licenses/Apache-2.0.html)
