![](https://codebuild.us-east-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiU1NyMHI4KytFRzhZSUVEY2R0YTlwanBJTk9EdWNYbW93TzdRU3NCbUJ0TFZYMy9jUktROXlUQktEOUVjd0dJSDBWbHNtVjVqSFpaNWxvbTJxd0o4dW53PSIsIml2UGFyYW1ldGVyU3BlYyI6ImgyNlBtRXoyU1ZSNjNWZjYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)
[![][sar-logo]](https://serverlessrepo.aws.amazon.com/applications/arn:aws:serverlessrepo:us-east-1:273450712882:applications~amazon-ecr-repository-compliance-webhook)

[sar-deploy]: https://img.shields.io/badge/Serverless%20Application%20Repository-Deploy%20Now-FF9900?logo=amazon%20aws&style=flat-square
[sar-logo]: https://img.shields.io/badge/Serverless%20Application%20Repository-View-FF9900?logo=amazon%20aws&style=flat-square

# Amazon ECR Repository Compliance Webhook for Kubernetes
>A Kubernetes ValidatingWebhookConfiguration and serverless backend: Deny Pods with container images that don't meet your compliance requirements

This AWS Serverless Application Repository app will create an Amazon API Gateway and an AWS Lambda Function that act as the backend for a Kubernetes ValidatingWebhookConfiguration. The function will deny Pods that create containers using images which:
1. Do not come from ECR
2. Come from ECR, but do not have image tag immutability enabled
3. Come from ECR, but do not have image scan on push enabled
4. Come from ECR, and have image scan on push enabled, but contain `CRITICAL` security vulnerabilities

![architecture](https://raw.githubusercontent.com/swoldemi/amazon-ecr-repository-compliance-webhook/master/screenshots/architecture.png)

## Usage
To use this SAR application:
1. Deploy the serverless application 
2. Configure and deploy the `ValidatingWebhookConfiguration` resource into your Kubernetes cluster (EKS or otherwise). The cluster must have this plugin enabled and have support for the admissionregistration.k8s.io/v1beta1 API. See the official Kubernetes documentation [here](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) for details. Amazon Elastic Kubernetes Service has [supported Dynamic Admission Controllers since October 12, 2018](https://aws.amazon.com/about-aws/whats-new/2018/10/amazon-eks-enables-support-for-kubernetes-dynamic-admission-cont/).

### 1. Deploying the Application
It is recommended that you deploy this Lambda function directly from the AWS Serverless Application Repository. It is also possible to deploy this function using:
- The [SAM CLI](https://aws.amazon.com/serverless/sam/)
- CloudFormation via the [AWS CLI](https://aws.amazon.com/cli/)
- CloudFormation via the [CloudFormation management console](https://aws.amazon.com/cloudformation/)

This function has been made available in 17 of the 18 commercial AWS regions that support AWS SAR. As of March 2020, Bahrain (me-south-1) does not yet support API Gateway. It is also possible to deploy the Lambda function in the GovCloud and China regions, if you have access to those regions.

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
|**South America (Sao Paulo) (sa-east-1)**     |[![][sar-deploy]](https://deploy.serverlessrepo.app/sa-east-1/?app=arn:aws:serverlessrepo:us-east-1:273450712882:applications/amazon-ecr-repository-compliance-webhook)     |

#### Parameters
|Name           |Default           |Description                                                       |Required |                 
|---------------|------------------|------------------------------------------------------------------|---------|
|RegistryRegion |Function's Region |What AWS region should this Lambda function interact with ECR in? |False    |
|LogLevel       |INFO              |The log level to set. ["DEBUG", "INFO", "WARN", "ERROR"]          |False    |

### 2. Configuration
After deploying the SAR application from the SAR console you need to:
1. Authenticate with your cluster. For example, for EKS you can use the AWS CLI: `aws eks update-kubeconfig --name your-clusters-name --region your-clusters-region`
2. Run `kubectl apply -f validatingwebhook.yaml` to deploy the `ValidatingWebhookConfiguration`. The YAML file is provided [here](https://github.com/swoldemi/amazon-ecr-repository-compliance-webhook/blob/master/deploy/validatingwebhook.yaml). Remember to update `webhooks.clientConfig.url` with your API Gateway endpoint. Make any necessary additions to match namespaces/labels for resources that are deployed. This webhook only validates Pods.
3. Run `kubectl create ns test-namespace && kubectl apply -f mydeployment.yaml` to create a sample `Deployment`. The sample is provided [here](https://github.com/swoldemi/amazon-ecr-repository-compliance-webhook/blob/master/deploy/mydeployment.yaml). Change the image to be any image you would like to test. Ensure your nodes have permission to pull from the ECR repository.
4. Run `kubectl get ev -n test-namespace` to see if there are any `FailedCreate` events as a result of the `Deployment`'s `ReplicaSet` triggering a failure from the `ValidatingWebhookConfiguration` when trying to create Pods. For example: `Error creating: admission webhook "admission.ecr.amazonaws.com" denied the request: webhook: no ecr images found in pod specification`

## Contributing
Have an idea for a feature to enhance this serverless application? Open an [issue](https://github.com/swoldemi/amazon-ecr-repository-compliance-webhook/issues) or [pull request](https://github.com/swoldemi/amazon-ecr-repository-compliance-webhook/pulls)!

### Development
This application has been developed, built, and tested against [Go 1.14](https://golang.org/dl/), the latest version of the [Serverless Application Model CLI](https://github.com/awslabs/aws-sam-cli), and the latest version of the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html), Kubernetes version 1.14, Kubernetes version 1.15, and [kubectl 1.17](https://kubernetes.io/docs/tasks/tools/install-kubectl/). A [Makefile](./Makefile) has been provided for convenience.

```
make install-tools # Install linting tools
make lint          # Run Go linting tools
make test          # Run Go tests
make compile       # Compile Go binary
make sam-package   # Package code and assets into S3 using SAM CLI
make sam-deploy    # Deploy application using SAM CLI
make sam-logs      # Tail the logs of the running Lambda function
make destroy-stack # Destroy the CloudFormation stack tied to the SAR app
```

### To Do
1. [Parameter.String] RegistryID - What registry should this Lambda verify container images for? Good for cross-account interactions.
2. [Parameter.CommaDelimitedList] IgnoredNamespaces - What namespaces should be ignored? It is also possible to set matchers on the [`ValidatingWebhookConfiguration`](./deploy/validatingwebhook.yaml).
3. Emit metric on deny/pass, to Amazon CloudWatch
4. Support the admissionregistration.k8s.io/v1 API

## References
- BanzaiCloud - In-depth introduction to Kubernetes admission webhooks: https://banzaicloud.com/blog/k8s-admission-webhooks/
- ValidatingWebhookConfiguration API Documentation - https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#validatingwebhookconfiguration-v1beta1-admissionregistration-k8s-io
- Dynamic Admission Control - https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
- Official Kubernetes example: https://github.com/kubernetes/kubernetes/blob/v1.15.0/test/images/webhook/

## Acknowledgements
[@jicowan](https://github.com/jicowan) for inspiration: https://github.com/jicowan/ecr-validation-webhook

## License
[Apache License 2.0](https://spdx.org/licenses/Apache-2.0.html)
