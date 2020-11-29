# aws-ecs-deploy-go
------------------
######Parameter
*必须参数
-ecr.name  </br>
-ecs.ecs_cluster.name <br>
-ecs.ecs_service.name
-ecs.ecs_task.container_name
-ecs.ecs_task.container_image

-ecs.ecs_task.name
-elb.name
-target_group.name

*default
-ecs.ecs_service.pubip=ENABLED
-ecs.ecs_service.task_num=1
-ecs.ecs_service.type=FARGATE
-ecs.ecs_task.container_networkmode=awsvpc
-ecs.ecs_task.container_port=80
-ecs.ecs_task.cpu=512
-ecs.ecs_task.executionrole=ecsTaskExecutionRole
-ecs.ecs_task.mem=1024
-ecs.ecs_task.mode=FARGATE
-ecs.ecs_task.role=ecsTaskExecutionRole
 -ecs.security_group=
 -elb.listener_port=80
 -elb.scheme=internet-facing
 -elb.type=application
 -protocol=HTTP
 -protocolversion=HTTP1
 -target_group.port=80
 -target_group.type=IP
 -vpcid =
 
