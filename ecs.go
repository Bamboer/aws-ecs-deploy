package main

import (
	"ecs/component"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger

	protocol        = flag.String("protocol", "HTTP", "ELB and Target group comunication protocol. HTTP|TCP|UDP|HTTPS|TLS|TCP_UDP")
	protocolversion = flag.String("protocolversion", "HTTP1", "Protocol version GRPC|HTTP2|HTTP1")
	elb_name        = flag.String("elb.name", "rewards-campaignstg-elb", "The Elb name.")
	elb_type        = flag.String("elb.type", "application", "The Elb type.application|network|gateway ")

	elb_listener_port = flag.Int64("elb.listener_port", 80, "elb listener port for ecs .default:80")
	elb_scheme        = flag.String("elb.scheme", "internet-facing", "The Elb type. internet-facing|internal")

	target_group_name = flag.String("target_group.name", "", "Target group name for ecs .")
	target_group_type = flag.String("target_group.type", "IP", "Target group target type for ecs. INSTANCE|IP|Lambda ")
	target_group_port = flag.Int64("target_group.port", 80, "Target group port for ecs.")

	security_group = flag.String("ecs.security_group", "sg-017e90ce7d6cd9b97", "Security groups name like: xxxx xxxx xxxx")
	subnets        = flag.String("subnet", "subnet-0e43caddfe68f606a subnet-00e19262c15c1de8b", "VPC's Subnets name like: xxxx xxxx xxxx")
	vpcid          = flag.String("vpcid", "vpc-3cbdfa58", "VPC id.")

	ecr_name = flag.String("ecr.name", "", "ECR name.")

	ecs_cluster = flag.String("ecs.ecs_cluster.name", "", "The ecs cluster name.")

	ecs_task_name                  = flag.String("ecs.ecs_task.name", "", "The ecs task info.")
	ecs_task_role                  = flag.String("ecs.ecs_task.role", "ecsTaskExecutionRole", "The ecs task role name .")
	ecs_task_exerole               = flag.String("ecs.ecs_task.executionrole", "ecsTaskExecutionRole", "The ecs task role name .")
	ecs_task_mode                  = flag.String("ecs.ecs_task.mode", "Fargate", "The ecs task RequiresCompatibilities.Fargate|EC2")
	ecs_task_mem                   = flag.Int64("ecs.ecs_task.mem", 1024, "The ecs task memory.")
	ecs_task_cpu                   = flag.Int64("ecs.ecs_task.cpu", 512, "The ecs task cpu.")
	ecs_task_container_name        = flag.String("ecs.ecs_task.container_name", "", "The ecs task container name.")
	ecs_task_container_port        = flag.Int64("ecs.ecs_task.container_port", 80, "The ecs task container port.")
	ecs_task_container_image       = flag.String("ecs.ecs_task.container_image", "", "The ecs task container image.")
	ecs_task_container_networkmode = flag.String("ecs.ecs_task.container_networkmode", "awsvpc", "The ecs task container network mode.bridge|host|awsvpc|none")
	ecs_task_loger                 = flag.String("ecs.ecs_task.loger", "", "The ecs task info.")

	ecs_service_name     = flag.String("ecs.ecs_service.name", "", "The ecs service info.")
	ecs_service_task_num = flag.Int64("ecs.ecs_service.task_num", 1, "The ecs service num.")
	ecs_service_pubip    = flag.String("ecs.ecs_service.pubip", "ENABLED", "The ecs service wether public ip DISABLED|ENABLED.")
	ecs_service_type     = flag.String("ecs.ecs_service.type", "FARGATE", "The ecs service launch type  EC2|FARGATE.")

	version = flag.Bool("v", false, "v1.0")
)

func init() {
	Info = log.New(os.Stdout,
		"Info: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(os.Stderr),
		"Error: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	flag.Parse()
	if *elb_type == "application" && *protocol != "HTTP" || protocol == nil {
		Error.Println("Parameter set error. elb type confict with protocol.")
		return
	} else if *elb_type == "network" && *protocol == "HTTP" {
		Error.Println("Parameter set error. elb type confict with protocol.")
		return
	} else if *ecr_name == "" {
		Error.Println("Parameter set error.")
		return
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ecr_arn := component.EcrCreator(sess, *ecr_name)
	fmt.Println(*ecr_arn)

	targetgroup_arn := component.CreateTargetGroup(sess, *elb_name, *target_group_name, *target_group_type, *protocol, *protocolversion, *vpcid, *target_group_port)
	fmt.Println(*targetgroup_arn)

	elb_name_arn := component.CreateLB(sess, *elb_name, *elb_type, *elb_scheme, *security_group, *subnets)
	fmt.Println(*elb_name_arn)

	listener_arn := component.CreateListener(sess, *elb_name, *protocol, *elb_name_arn, *targetgroup_arn, *elb_listener_port)
	fmt.Println(*listener_arn)

	ecscluster_arn := component.CreateEcsCluster(sess, *ecs_cluster)
	fmt.Println(*ecscluster_arn)

	ecstask_arn := component.CreateEcsTask(sess, *ecs_task_cpu, *ecs_task_mem, *ecs_task_container_name, *ecs_task_container_image, *ecs_task_name, *ecs_task_role, *ecs_task_exerole, *ecs_task_mode, *ecs_task_container_networkmode)
	fmt.Println(*ecstask_arn)

	ecsservice_arn := component.CreateEcsService(sess, *ecs_service_task_num, *ecs_task_container_port, *ecs_cluster, *ecs_service_name, *ecs_task_container_name, *security_group, *ecs_service_type, *elb_name, *ecstask_arn, *targetgroup_arn, *subnets, *ecs_service_pubip)
	fmt.Println(*ecsservice_arn)
}
