#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { LoginServiceStack } from '../lib/login-service-stack';

const app = new cdk.App();

new LoginServiceStack(app, 'LoginService', {
  env: {
    account: process.env.AWS_ACCOUNT_ID,
    region: process.env.AWS_REGION,
  }
});