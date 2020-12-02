#!/usr/bin/python env
#coding:utf8

import sys
import boto3
import argparse


__author__ = 'Bambo'
parser = argparse.ArgumentParser()

parser.add_argument('-elb_name',type=str,help='elb name')
parser.add_argument('-elb_type',type=str,default='application',help='elb type')
parser.add_argument('-elb_scheme',type=str,default='internal',help='elb scheme')

parser.add_argument('-ecr_repositoryuri',type=str,help='ecr repository address')
parser.add_argument('-ecr_name',type=str,help='ecr  name')

parser.add_argument('-ecs_cluster_name',type=str,help='ecs cluster name')

parser.add_argument('-ecs_task_name',type=str,help='ecr repository name')
parser.add_argument('-ecs_task_cpu',type=str,default='512',help='ecs task cpu.default=512 # 0.25 vCPU 的有效内存范围 = 0.5GB - 2GB')
parser.add_argument('-ecs_task_mem',type=str,default='1024',help='ecs task memory. default=1024 #1GB 内存的有效 CPU 范围 = 0.25 vCPU - 0.5 vCPU。')
parser.add_argument('-ecs_task_type',type=str,default='FARGATE',help='ecs task requiresCompatibilities type FARGATE|EC2 default:FARGATE')
parser.add_argument('-ecs_task_networkmode',type=str,default='awsvpc',help='ecs task network mode. bridge|host|awsvpc|none default:awsvpc')
parser.add_argument('-ecs_task_container_name',type=str,help='ecs task container name')
parser.add_argument('-ecs_task_container_image',type=str,help='ecs task container image.default= ecr/imagename')
parser.add_argument('-ecs_task_container_port',type=int,default=80,help='ecs task container port')
parser.add_argument('-ecs_task_container_cpu',type=int,default=512,help='ecs task container port')
parser.add_argument('-ecs_task_container_mem',type=int,default=1024,help='ecs task container port')
parser.add_argument('-ecs_task_container_protocol',type=str,default='tcp',help='ecs task container protocol.default=tcp')
parser.add_argument('-ecs_task_container_role',type=str,default='ecsTaskExecutionRole',help='ecs task container port')
parser.add_argument('-ecs_task_role',type=str,default='ecsTaskExecutionRole',help='ecs task container image')

parser.add_argument('-ecs_service_name',type=str,help='ecs service name.')
parser.add_argument('-ecs_service_pubip',type=str,default='ENABLED',help='ecs service public ip DISABLED|ENABLED.default: ENABLED')
parser.add_argument('-ecs_service_role',type=str,default='ecsTaskExecutionRole',help='ecs service name.')
parser.add_argument('-ecs_service_desier_num',type=int,default=1,help='ecs service desired count. default: 1')
parser.add_argument('-ecs_service_type',type=str,default='FARGATE',help='ecs service launch type. EC2|FARGATE default:FARGATE')

parser.add_argument('-vpcid',type=str,default='vpc-3cbdfa58',help='vpc id.')
parser.add_argument('-subnets',type=str,default='subnet-0f757a90deb6870e0,subnet-068110c72ac9bfc65',help='subnets.')
parser.add_argument('-securitygroups',type=str,default='sg-0b8d77f79f681fc25',help='securitygroups.')

parser.add_argument('-targetgroup_name',type=str,help='target group name.')
parser.add_argument('-targetgroup_protocol',type=str,default='TCP',help='target group protocol. default: TCP')
parser.add_argument('-targetgroup_port',type=int,default=80,help='target group port.default: 80')
parser.add_argument('-targetgroup_healthcheck',type=bool,default=False,help='target group helth check True|False. default: False')
parser.add_argument('-targetgroup_targettype',type=str,default='ip',help='taget group target type. default: ip')
parser.add_argument('-targetgroup_healthcheckport',type=str,help='targetgroup healthcheck port.')
parser.add_argument('-targetgroup_checkpath',type=str,help='target group health check path.')
parser.add_argument('-targetgroup_healthcheckprotocol',type=str,default='HTTP',help='target group health check path.')


args = parser.parse_args()

ecs = boto3.client('ecs')
ecr = boto3.client('ecr')
elb = boto3.client('elbv2')
iam = boto3.client('iam')

def get_role(rolename):
    response = iam.get_role(
            RoleName = rolename
            )
    return response['Role']['Arn']

def CreateListener():
    if args.elb_type == "application" and  "HTTP" not in args.targetgroup_protocol:
        print("Parameter set error.")
        sys.exit(1)
    elif args.elb_type == "network" and  "HTTP" in args.targetgroup_protocol:
        print("Parameter set error.")
        sys.exit(1)
    try:
        response = elb.create_listener(
          LoadBalancerArn= CreateElb(),
          Protocol=  args.targetgroup_protocol ,#'HTTP'|'HTTPS'|'TCP'|'TLS'|'UDP'|'TCP_UDP',
          Port=  args.targetgroup_port,
          DefaultActions=[
              {
                  'Type': 'forward',#'forward'|'authenticate-oidc'|'authenticate-cognito'|'redirect'|'fixed-response'
                  'TargetGroupArn': create_target_group(),
              },
          ]
         )
    except Exception as e:
        print("Not create Listener: ",e)
        return
    return response


def CreateElb():
    if args.elb_name ==None:
        print('Not create ELB ')
        return
    try:
        response = elb.describe_load_balancers(Names=[args.elb_name,])
        print("ELB {0} Exists. ".format(args.elb_name))
        return response['LoadBalancers'][0]['LoadBalancerArn']
    except:
        response = elb.create_load_balancer(
            Name=args.elb_name,
            Subnets=args.subnets.split(' '),
            SecurityGroups=args.securitygroups.split(' '),
            Scheme = args.elb_scheme,  #'internet-facing'|'internal',
            Type=args.elb_type,  #'application'|'network'|'gateway',
            IpAddressType='ipv4',
            Tags = [
                {
                    'Key': 'test',
                    'Value': 'test1'
                    },
                ]
            )
        print('ELB {0 }Create.'.format(args.elb_name))
        return response['LoadBalancers'][0]['LoadBalancerArn']

def create_target_group():
    if args.targetgroup_name == None:
        print('Not create target group')
        return
    try:
        response = elb.describe_target_groups(Names=[args.targetgroup_name,])
        print('{0} Target Group Exists.'.format(args.targetgroup_name))
        return response['TargetGroups'][0]['TargetGroupArn']
    except:
        response = elb.create_target_group(
            Name = args.targetgroup_name,
            Protocol = args.targetgroup_protocol,  #'HTTP'|'HTTPS'|'TCP'|'TLS'|'UDP'|'TCP_UDP'|'GENEVE',
            Port = args.targetgroup_port,
            VpcId = args.vpcid,
#            HealthCheckProtocol= args.targetgroup_healthcheckprotocol ,#'HTTP'|'HTTPS'|'TCP'|'TLS'|'UDP'|'TCP_UDP'|'GENEVE',
#            HealthCheckPort= args.targetgroup_healthcheckport,
            HealthCheckEnabled= True if args.targetgroup_targettype == 'ip' else  args.targetgroup_healthcheck,#True|False,
#            HealthCheckPath=args.targetgroup_checkpath,
#            HealthCheckIntervalSeconds=30,
#            HealthCheckTimeoutSeconds=5,
#            HealthyThresholdCount=3,
#            UnhealthyThresholdCount=5,
            TargetType= args.targetgroup_targettype ,#'instance'|'ip'|'lambda',
#            Matcher = {
#                'HttpCode': '200',
#                },
            )
        print('{0} Target group Create. '.format(args.targetgroup_name))
        return response['TargetGroups'][0]['TargetGroupArn']

def task_define():
    if args.ecs_task_container_name == None or args.ecs_task_name == None or args.ecs_task_container_image == None:
        print("ECS task definition args not enough.")
        return
    try:
        response = ecs.describe_task_definition(taskDefinition=args.ecs_task_name,)
        print('Task {0} Exists: '.format(args.ecs_task_name))
        return response['taskDefinition']['taskDefinitionArn']
    except:
        response = ecs.register_task_definition(
            family = args.ecs_task_name,
            taskRoleArn = args.ecs_task_role,
            executionRoleArn= args.ecs_task_container_role,
            networkMode = args.ecs_task_networkmode,
            containerDefinitions = [
                {
                    'name': args.ecs_task_container_name,
                    'image': create_ecr(),
                    'cpu': args.ecs_task_container_cpu,
                    'memory': args.ecs_task_container_mem,
                    'portMappings':[{
                        'containerPort': args.ecs_task_container_port,
                        'protocol': args.ecs_task_container_protocol,
                   }],
                    },],
            cpu = args.ecs_task_cpu,
            memory = args.ecs_task_mem,
            requiresCompatibilities = [args.ecs_task_type,],
            )
        print("Task{0} Create. ".format(args.ecs_task_name))
        return response['taskDefinition']['taskDefinitionArn']

def create_cluster():
    if args.ecs_cluster_name != None:
        try:
            response = ecs.describe_clusters(clusters=[args.ecs_cluster_name,])
            if response['clusters'][0]['status'] == 'INACTIVE':
                raise  Exception("inactive")
            if response['clusters'] == [] and response['failures'] != []:
                raise Exception("failures")
            else:
                print('Cluster {0} exists.'.format(args.ecs_cluster_name))
                #print(response)
                return response['clusters'][0]['clusterArn']
        except Exception as e:
            print("Exception is: ",e)
            response = ecs.create_cluster(
                clusterName = args.ecs_cluster_name,
            )
            print('Create ECS Cluster: {0}'.format(args.ecs_cluster_name))
#            print(response)
            return response['cluster']['clusterArn']

def create_ecr():
    if args.ecr_name != '' and args.ecr_name != None:
        try:
            response = ecr.create_repository(
               repositoryName = args.ecr_name,
            )
            token = ecr.get_authorization_token(registryIds=[response['repository']['registryId'],])
            print('Create ECR Repository {0}\'s  authorization token: {1}'.format(args.ecr_name,token['authorizationData']))
            return response['repository']['repositoryUri']
        except:
            response = ecr.describe_repositories(repositoryNames=[args.ecr_name,])
            token = ecr.get_authorization_token(registryIds=[response['repositories'][0]['registryId'],])
            print('ECR {0} exists.'.format(args.ecr_name))
            print('ECR Repository {0}\'s  authorization token: {1}'.format(args.ecr_name,token['authorizationData']))
            return response['repositories'][0]['repositoryUri']
    else:
        print("Not create ECR Repository.")
        sys.exit(1)
        return


def create_service():
    if args.ecs_cluster_name == None or args.ecs_service_name == None or args.ecs_task_name == None:
        print("Not create ecs service. Parameter not enough.")
        return
    try:
        response = ecs.describe_services(cluster=args.ecs_cluster_name,services = [args.ecs_service_name,])
        print('Service {0} Exists.'.format(args.ecs_service_name))
        return response['services'][0]['services']
    except:
        response = ecs.create_service(
#            role = args.ecs_service_role,
            cluster = args.ecs_cluster_name,
            serviceName = args.ecs_service_name,
            taskDefinition = task_define(),
            loadBalancers =[{
                'targetGroupArn': args.targetgroup_arn if create_target_group() == None else create_target_group(),
#                'loadBalancerName': args.elb_name,
                'containerName': args.ecs_task_container_name,
                'containerPort': args.ecs_task_container_port
                },
                ],
            desiredCount = args.ecs_service_desier_num,
            launchType = args.ecs_service_type,
            networkConfiguration = {
                'awsvpcConfiguration':{
                    'subnets': args.subnets.split(' '),
                    'securityGroups': args.securitygroups.split(' '),
                    'assignPublicIp': args.ecs_service_pubip
                    }
                },
            healthCheckGracePeriodSeconds=30,
            schedulingStrategy='REPLICA', #|'DAEMON',
            enableECSManagedTags=True
            )
        print('Create Service: {0}'.format(args.ecs_service_name))
        return response['service']['serviceArn']

if __name__ == '__main__':
    create_cluster()
    create_ecr()
    CreateListener()
    task_define()
    create_service()
