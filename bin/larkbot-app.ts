#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { LarkbotAppStack } from '../lib/larkbot-app-stack';

const app = new cdk.App();
const stackname = app.node.tryGetContext('stackName')
new LarkbotAppStack(app, 'LarkbotAppStack', {
  stackName: stackname,
  env: {
    account: process.env.CDK_DEPLOY_ACCOUNT || process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEPLOY_REGION || process.env.CDK_DEFAULT_REGION
  }
});

