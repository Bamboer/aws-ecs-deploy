# aws-ecs-deploy-go


### Parameter

###### 必须参数
-ecr.name  <br>

-ecs.ecs_cluster.name <br>
-ecs.ecs_service.name <br>
-ecs.ecs_task.container_name <br>
-ecs.ecs_task.container_image <br>

-ecs.ecs_task.name <br>

-elb.name <br>

-target_group.name <br>

###### default
-ecs.ecs_service.pubip=ENABLED <br>
-ecs.ecs_service.task_num=1 <br>
-ecs.ecs_service.type=FARGATE <br>
-ecs.ecs_task.container_networkmode=awsvpc <br>
-ecs.ecs_task.container_port=80 <br>
-ecs.ecs_task.cpu=512 <br>
-ecs.ecs_task.executionrole=ecsTaskExecutionRole <br>
-ecs.ecs_task.mem=1024 <br>
-ecs.ecs_task.mode=FARGATE <br>
-ecs.ecs_task.role=ecsTaskExecutionRole <br>

-ecs.security_group= <br>

-elb.listener_port=80 <br>
-elb.scheme=internet-facing <br>
-elb.type=application <br>

-protocol=HTTP <br>
-protocolversion=HTTP1 <br>

-target_group.port=80 <br>
-target_group.type=IP <br>

-vpcid = <br>

-subnets= <br>
