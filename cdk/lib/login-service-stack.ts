import * as ec2 from '@aws-cdk/aws-ec2';
import * as ecs from '@aws-cdk/aws-ecs';
import * as ecsp from '@aws-cdk/aws-ecs-patterns';
import * as cdk from '@aws-cdk/core';
import * as autoscaling from '@aws-cdk/aws-autoscaling';

export class LoginServiceStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const vpc = new ec2.Vpc(this, 'vpc', { maxAzs: 2 });

        const autoScalingGroup = new autoscaling.AutoScalingGroup(this, 'login-service-scaling-group', {
            vpc: vpc,
            instanceType: new ec2.InstanceType('t2.micro'),
            machineImage: ecs.EcsOptimizedImage.amazonLinux2(),
            minCapacity: 0,
            maxCapacity: 1,
        });

        const capacityProvider = new ecs.AsgCapacityProvider(this, 'login-service-capacity-provider', {
            autoScalingGroup: autoScalingGroup,
        });

        const cluster = new ecs.Cluster(this, 'login-service-cluster', { 
            vpc: vpc,
        });
        cluster.addAsgCapacityProvider(capacityProvider);

        const taskDefinition = new ecs.Ec2TaskDefinition(this, 'login-service-task-definition');
        taskDefinition.addContainer('web', {
            image: ecs.ContainerImage.fromAsset(`${__dirname}/../..`),
            portMappings: [
                { hostPort: 8080, containerPort: 8080 },
            ],
            memoryReservationMiB: 128,
        });

        const loadBalancedService = new ecsp.ApplicationLoadBalancedEc2Service(
            this, 'login-service-balancer', {
                cluster: cluster,
                taskDefinition: taskDefinition,
                desiredCount: 1,
                publicLoadBalancer: true,
            }
        );

        new cdk.CfnOutput(this, 'LoadBalancerDns', { value: loadBalancedService.loadBalancer.loadBalancerDnsName });
    }
}