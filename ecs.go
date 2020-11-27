package main

import (
        "flag"
//      "github.com/aws/aws-sdk-go/aws"
//      "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/session"
        "fmt"
        "io"
        "log"
        "os"
//        "string"
        "ecs/component"
        //  "github.com/aws/aws-sdk-go/service/elb"
        //  "github.com/aws/aws-sdk-go/service/ecs"
        //  "github.com/aws/aws-sdk-go/service/ecr"
)

var (
        Info  *log.Logger
        Error *log.Logger

        elb_name = flag.String("elb.name", "rewards-campaignstg-elb", "The Elb name.")
        elb_port = flag.String("elb.port", "", "The Elb listening port.")
        elb_type = flag.String("elb.type", "", "The Elb type.")
        elb_type = flag.String("elb.scheme", "", "The Elb type.")

        target_group      = flag.Bool("target_group", false, "Target group name for ecs .")
        target_group_name = flag.String("target_group.name", "", "Target group name for ecs .")
        target_group_type = flag.String("target_group.type", "IP", "Target group target type for ecs .")
        target_group_port = flag.Int64("target_group.port", 80, "Target group port for ecs .")
        target_group_vpc  = flag.String("target_group.vpc", "vpc-3cbdfa58", "VPC for Target group used.")

        security_group = flag.String("ecs.security_group", "sg-0b8d77f79f681fc25", "Security groups name like: xxxx xxxx xxxx")
        subnet = flag.String("subnet", "subnet-0f757a90deb6870e0 subnet-068110c72ac9bfc65", "Security groups name like: xxxx xxxx xxxx")

        ecr_name     = flag.String("ecr.name", "", "ECR name.")
        ecr_img_name = flag.String("ecr.img_name", "", "ECR image name.")

        ecs_cluster = flag.String("ecs.ecs_cluster.name", "", "The ecs cluster name.")

        ecs_task_name           = flag.String("ecs.ecs_task.name", "", "The ecs task info.")
        ecs_task_role           = flag.String("ecs.ecs_task.role", "ecsTaskExecutionRole", "The ecs task role name .")
        ecs_task_type           = flag.String("ecs.ecs_task.type", "", "The ecs task info.")
        ecs_task_mem            = flag.Int64("ecs.ecs_task.mem", 1024, "The ecs task info.")
        ecs_task_cpu            = flag.String("ecs.ecs_task.cpu","0.5vCPU","The ecs task info.")
        ecs_task_container_name = flag.String("ecs.ecs_task.container_name", "", "The ecs task info.")
        ecs_task_loger          = flag.String("ecs.ecs_task.loger", "", "The ecs task info.")
        ecs_task_provider       = flag.String("ecs.ecs_task.provider", "Fargate", "The ecs task provider.Fargate|EC2")

        ecs_service_name         = flag.String("ecs.ecs_service.name", "", "The ecs service info.")
        ecs_service_version      = flag.String("ecs.ecs_service.version", "LATEST", "The ecs service info.")
        ecs_service_task_num     = flag.Int64("ecs.ecs_service.task_num", 1, "The ecs service num.")
        ecs_service_pubip        = flag.String("ecs.ecs_service.pubip", "DISABLED", "The ecs service wether public ip DISABLED|ENABLED.")
        ecs_service_check_period = flag.Int64("ecs.ecs_service.check_period", 300, "The ecs service info.")

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
        sess := session.Must(session.NewSessionWithOptions(session.Options{
                SharedConfigState: session.SharedConfigEnable,
        }))
        if *ecr_name != ""{
              ecr_info := component.EcrCreator(sess,*ecr_name)
              fmt.Println(*ecr_info)
        }
}
