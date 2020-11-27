package component

import (
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/ecs"
        "strings"
)

func CreateEcsCluster(sess *session.Session, cluster_name string) *string {
        if sess == nil || cluster_name == "" {
                Error.Println("Not create ECS Cluster. parameter not enough.")
                return nil
        }
        svc := ecs.New(sess)
        input := &ecs.DescribeClustersInput{
                Clusters: []*string{
                        aws.String(cluster_name),
                },
        }
        result, err := svc.DescribeClusters(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecs.ErrCodeServerException:
                                Error.Println(ecs.ErrCodeServerException, aerr.Error())
                        case ecs.ErrCodeClientException:
                                Error.Println(ecs.ErrCodeClientException, aerr.Error())
                        case ecs.ErrCodeInvalidParameterException:
                                Error.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        Error.Println(err.Error())
                }
                return nil
        } else {
                Info.Println("Exists ECS cluster: ", result.Clusters[0].ClusterArn)
                return result.Clusters[0].ClusterArn
        }

        input := &ecs.CreateClusterInput{
                ClusterName: aws.String(cluster_name),
        }

        result, err := svc.CreateCluster(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecs.ErrCodeServerException:
                                Error.Println(ecs.ErrCodeServerException, aerr.Error())
                        case ecs.ErrCodeClientException:
                                Error.Println(ecs.ErrCodeClientException, aerr.Error())
                        case ecs.ErrCodeInvalidParameterException:
                                Error.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
        }
        Info.Println("Create ECS Cluster: ", result.Cluster.ClusterArn)
        return result.Cluster.ClusterArn
}

func CreateEcsService(sess *session.Session, desire_num, container_port int64, cluster_name, service_type, container_name, lbname, service_name, task_arn string) *string {
        if sess == nil || service_name == "" {
                Error.Println("Not create ECS Service. not enough parameter.")
                return nil
        }
        svc := ecs.New(sess)
        input := &ecs.DescribeServicesInput{
                Services: []*string{
                        aws.String(service_name),
                },
        }
        result, err := svc.DescribeServices(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecs.ErrCodeServerException:
                                Error.Println(ecs.ErrCodeServerException, aerr.Error())
                        case ecs.ErrCodeClientException:
                                Error.Println(ecs.ErrCodeClientException, aerr.Error())
                        case ecs.ErrCodeInvalidParameterException:
                                Error.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
                        case ecs.ErrCodeClusterNotFoundException:
                                Error.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        Error.Println(err.Error())
                }
                return nil
        } else {
                Info.Println("Exists ECS Service: ", result.Services[0].ServiceArn)
                return result.Services[0].ServiceArn
        }

        //
        input := &ecs.CreateServiceInput{
                DesiredCount: aws.Int64(desire_num),
                LoadBalancers: []*ecs.LoadBalancer{
                        {
                                ContainerName:    aws.String(container_name),
                                ContainerPort:    aws.Int64(container_port),
                                LoadBalancerName: aws.String(lbname),
                        },
                },
                //              Role:           aws.String("ecsServiceRole"),
                ServiceName:    aws.String(service_name),
                TaskDefinition: aws.String(task_arn),
                Cluster:        aws.String(cluster_name),
                LaunchType:     aws.String(service_type),
        }

        result, err := svc.CreateService(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecs.ErrCodeServerException:
                                Error.Println(ecs.ErrCodeServerException, aerr.Error())
                        case ecs.ErrCodeClientException:
                                Error.Println(ecs.ErrCodeClientException, aerr.Error())
                        case ecs.ErrCodeInvalidParameterException:
                                Error.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
                        case ecs.ErrCodeClusterNotFoundException:
                                Error.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
                        case ecs.ErrCodeUnsupportedFeatureException:
                                Error.Println(ecs.ErrCodeUnsupportedFeatureException, aerr.Error())
                        case ecs.ErrCodePlatformUnknownException:
                                Error.Println(ecs.ErrCodePlatformUnknownException, aerr.Error())
                        case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
                                Error.Println(ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())
                        case ecs.ErrCodeAccessDeniedException:
                                Error.Println(ecs.ErrCodeAccessDeniedException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
                return nil
        }

        Info.Println("Create ECS Service: ", result.Service.ServiceArn)
        return result.Service.ServiceArn
}

func CreateEcsTask(sess *session.Session, cpu, mem int64, container_name, image, task_name, role_arn, task_mode string) *string {
        if sess == nil || task_name == "" {
                Error.Println("Not create task definition. not enough parameter.")
                return nil
        }
        svc := ecs.New(sess)
        input := &ecs.DescribeTaskDefinitionInput{
                TaskDefinition: aws.String(task_name),
        }

        result, err := svc.DescribeTaskDefinition(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecs.ErrCodeServerException:
                                Error.Println(ecs.ErrCodeServerException, aerr.Error())
                        case ecs.ErrCodeClientException:
                                Error.Println(ecs.ErrCodeClientException, aerr.Error())
                        case ecs.ErrCodeInvalidParameterException:
                                Error.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
                return nil
        } else {
                Info.Println("Exists Task definition: ", result.TaskDefinition.TaskDefinitionArn)
                return result.TaskDefinition.TaskDefinitionArn
        }

        //
        input := &ecs.RegisterTaskDefinitionInput{
                ContainerDefinitions: []*ecs.ContainerDefinition{
                        {
                                Cpu:       aws.Int64(cpu),
                                Essential: aws.Bool(true),
                                Image:     aws.String(image),
                                Memory:    aws.Int64(mem),
                                Name:      aws.String(container_name),
                        },
                },
                Family:                  aws.String(task_name),
                TaskRoleArn:             aws.String(role_arn),
                ExecutionRoleArn:        aws.String(role_arn),
                RequiresCompatibilities: []*string{aws.String(task_mode)},
                Tags: []*ecs.Tag{
                        {Key: aws.String("GBL_CLASS_0"),
                                Value: aws.String("SERVICE"),
                        },
                        {Key: aws.String("GBL_CLASS_1"),
                                Value: aws.String("SA-STG-Cluster"),
                        },
                        {Key: aws.String("GBL_CLASS_2"),
                                Value: aws.String("ECS"),
                        },
                },
        }

        result, err := svc.RegisterTaskDefinition(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecs.ErrCodeServerException:
                                Error.Println(ecs.ErrCodeServerException, aerr.Error())
                        case ecs.ErrCodeClientException:
                                Error.Println(ecs.ErrCodeClientException, aerr.Error())
                        case ecs.ErrCodeInvalidParameterException:
                                Error.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
        }

        Info.Println("Create Task definition: ", result.TaskDefinition.TaskDefinitionArn)
        return result.TaskDefinition.TaskDefinitionArn
}
