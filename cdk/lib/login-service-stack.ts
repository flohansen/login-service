import * as ec2 from '@aws-cdk/aws-ec2';
import * as ecs from '@aws-cdk/aws-ecs';
import * as ecsp from '@aws-cdk/aws-ecs-patterns';
import * as cdk from '@aws-cdk/core';

export class LoginServiceStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const vpc = new ec2.Vpc(this, 'vpc', { maxAzs: 2 });
        const cluster = new ecs.Cluster(this, 'cluster', { vpc })
        const fargateservice = new ecsp.ApplicationLoadBalancedFargateService(
            this, 'fargateservice', {
                cluster: cluster,
                taskImageOptions: {
                    image: ecs.ContainerImage.fromAsset(`${__dirname}/../..`),
                    containerPort: 8080,
                },
                desiredCount: 1,
                publicLoadBalancer: true,
            }
        );

        new cdk.CfnOutput(this, 'LoadBalancerDns', { value: fargateservice.loadBalancer.loadBalancerDnsName });
    }
}