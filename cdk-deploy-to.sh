#!/usr/bin/env bash
if [[ $# -ge 2 ]]; then
    export CDK_DEPLOY_ACCOUNT=$1
    export CDK_DEPLOY_REGION=$2
    shift; shift
    npx cdk deploy "$@"
    exit $?
else
    echo 1>&2 "Provide account and region as first two args."
    echo 1>&2 "Additional args are passed through to cdk deploy."
    exit 1
fi

## Reference code
## https://docs.aws.amazon.com/cdk/v2/guide/environments.html

## Example
## ./cdk-deploy-to.sh 123456789 us-east-1 "$@"