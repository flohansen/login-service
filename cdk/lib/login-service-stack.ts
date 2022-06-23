import { Vpc } from '@aws-cdk/aws-ec2';
import { Cluster, ContainerImage } from '@aws-cdk/aws-ecs';
import { ApplicationLoadBalancedFargateService } from '@aws-cdk/aws-ecs-patterns';
import { Stack, Construct, StackProps, CfnOutput } from '@aws-cdk/core'

export class LoginServiceStack extends Stack {
    constructor(scope: Construct, id: string, props?: StackProps) {
        super(scope, id, props);

        const vpc = new Vpc(this, 'vpc', { maxAzs: 2 });
        const cluster = new Cluster(this, 'Cluster', { vpc })
        const fargateService = new ApplicationLoadBalancedFargateService(
            this, 'FargateService', {
                cluster: cluster,
                taskImageOptions: {
                    image: ContainerImage.fromAsset(`${__dirname}/../..`),
                    containerPort: 8080,
                    environment: {
                        DEPLOYED_DATE: Date.now().toLocaleString(),
                    },
                },
                desiredCount: 1,
                publicLoadBalancer: true,
            }
        );

        new CfnOutput(this, 'LoadBalancerDNS', { value: fargateService.loadBalancer.loadBalancerDnsName });
    }
}