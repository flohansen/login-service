import * as cdk from '@aws-cdk/core';
import { LoginServiceStack } from '../lib/login-service-stack';
import { Template } from '@aws-cdk/assertions';

test('LoginServiceStack created', () => {
    const app = new cdk.App()

    const loginServiceStack = new LoginServiceStack(app, 'LoginServiceStack');

    const template = Template.fromStack(loginServiceStack);
    template.hasResource('AWS::ECS::Cluster', {});
    template.hasResource('AWS::EC2::VPC', {})
});
