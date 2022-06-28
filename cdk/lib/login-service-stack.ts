import { Vpc } from '@aws-cdk/aws-ec2';
import { Cluster, ContainerImage } from '@aws-cdk/aws-ecs';
import { ApplicationLoadBalancedFargateService } from '@aws-cdk/aws-ecs-patterns';
import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as path from 'path';
import * as apigateway from '@aws-cdk/aws-apigateway';
import { cloud_assembly_schema } from 'aws-cdk-lib';

export class LoginServiceStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // const vpc = new Vpc(this, 'vpc', { maxAzs: 2 });
        // const cluster = new Cluster(this, 'Cluster', { vpc })
        // const fargateService = new ApplicationLoadBalancedFargateService(
        //     this, 'FargateService', {
        //         cluster: cluster,
        //         taskImageOptions: {
        //             image: ContainerImage.fromAsset(`${__dirname}/../..`),
        //             containerPort: 8080,
        //         },
        //         desiredCount: 1,
        //         publicLoadBalancer: true,
        //     }
        // );

        // new CfnOutput(this, 'LoadBalancerDNS', { value: fargateService.loadBalancer.loadBalancerDnsName });
        const fn = new lambda.Function(this, 'LoginServiceLambda', {
            code: lambda.Code.fromAsset(path.join(__dirname, '../../src')),
            handler: 'main.main',
            runtime: lambda.Runtime.GO_1_X,
        });

        const api = new apigateway.LambdaRestApi(this, 'LoginServiceRestApi', {
            handler: fn,
            description: 'Login Service REST API',
        });

        new cdk.CfnOutput(this, 'apiUrl', { value: api.url });
    }
}