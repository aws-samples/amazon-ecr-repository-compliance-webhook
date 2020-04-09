#!/bin/bash

BUCKET_PREFIX=swoldemi
APPLICATION=amazon-ecr-repository-compliance-webhook
ACCOUNT=273450712882
declare -a REGIONS=(
    "us-east-2"
    "us-east-1"
    "us-west-1"
    "us-west-2"
    "ap-east-1"
    "ap-south-1"
    "ap-northeast-2"
    "ap-southeast-1"
    "ap-southeast-2"
    "ap-northeast-1"
    "ca-central-1"
    "eu-central-1"
    "eu-west-1"
    "eu-west-2"
    "eu-west-3"
    "eu-north-1"
    "sa-east-1"
)

publish_all_regions()
{   
    for REGION in "${REGIONS[@]}"
    do
        echo Deploying to region $REGION
        sam package --template-file template.yaml --s3-bucket $BUCKET_PREFIX-$REGION --output-template-file packaged.yaml
        sam publish --region $REGION --template packaged.yaml
    done

    aws serverlessrepo put-application-policy \
        --region us-east-1 \
        --application-id arn:aws:serverlessrepo:us-east-1:$ACCOUNT:applications/$APPLICATION \
        --statements Principals=*,Actions=Deploy
}

echo "On branch `basename $CODEBUILD_WEBHOOK_HEAD_REF`"
if [ "`basename $CODEBUILD_WEBHOOK_HEAD_REF`" = "master" ]
then 
    publish_all_regions
else
    echo Skipping publish
fi
