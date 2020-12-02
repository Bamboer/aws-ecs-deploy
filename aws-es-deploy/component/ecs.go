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
	} else {
		Info.Println("Exists ECS cluster: ", result.Clusters[0].ClusterArn)
		return result.Clusters[0].ClusterArn
	}

	inputc := &ecs.CreateClusterInput{
		ClusterName: aws.String(cluster_name),
	}

	resultc, errc := svc.CreateCluster(inputc)
	if err != nil {
		if aerr, ok := errc.(awserr.Error); ok {
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
			Error.Println(errc.Error())
		}
	}
	Info.Println("Create ECS Cluster: ", resultc.Cluster.ClusterArn)
	return resultc.Cluster.ClusterArn
}

func CreateEcsService(sess *session.Session, desire_num, container_port int64, cluster_name, service_name, container_name, security_group, service_type, lbname, task_arn, target_group_arn, subnets, ecs_pubip string) *string {
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
	} else {
		Info.Println("Exists ECS Service: ", result.Services[0].ServiceArn)
		return result.Services[0].ServiceArn
	}

	//
	inputc := &ecs.CreateServiceInput{
		DesiredCount: aws.Int64(desire_num),
		LoadBalancers: []*ecs.LoadBalancer{
			{
				ContainerName:  aws.String(container_name),
				ContainerPort:  aws.Int64(container_port),
				TargetGroupArn: aws.String(target_group_arn),
			},
		},
		//              Role:           aws.String("ecsServiceRole"),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        aws.StringSlice(strings.Split(subnets, " ")),
				SecurityGroups: aws.StringSlice(strings.Split(security_group, " ")),
				AssignPublicIp: aws.String(ecs_pubip),
			}},
		ServiceName:    aws.String(service_name),
		TaskDefinition: aws.String(task_arn),
		Cluster:        aws.String(cluster_name),
		LaunchType:     aws.String(service_type),
	}

	resultc, errc := svc.CreateService(inputc)
	if err != nil {
		if aerr, ok := errc.(awserr.Error); ok {
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
			Error.Println(errc.Error())
		}
		return nil
	}

	Info.Println("Create ECS Service: ", resultc.Service.ServiceArn)
	return resultc.Service.ServiceArn
}

func CreateEcsTask(sess *session.Session, cpu, mem int64, container_name, image, task_name, role_arn, execution_role_arn, task_mode, networkmode string) *string {
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
	} else {
		Info.Println("Exists Task definition: ", result.TaskDefinition.TaskDefinitionArn)
		return result.TaskDefinition.TaskDefinitionArn
	}

	//
	inputc := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Cpu:       aws.Int64(cpu),
				Essential: aws.Bool(true),
				Image:     aws.String(image),
				Memory:    aws.Int64(mem),
				Name:      aws.String(container_name),
			},
		},
		NetworkMode:             aws.String(networkmode),
		Family:                  aws.String(task_name),
		TaskRoleArn:             aws.String(role_arn),
		ExecutionRoleArn:        aws.String(execution_role_arn),
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

	resultc, errc := svc.RegisterTaskDefinition(inputc)
	if errc != nil {
		if aerr, ok := errc.(awserr.Error); ok {
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
			Error.Println(errc.Error())
		}
		return nil
	}

	Info.Println("Create Task definition: ", resultc.TaskDefinition.TaskDefinitionArn)
	return resultc.TaskDefinition.TaskDefinitionArn
}
