package component

import (
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/elbv2"
        "strings"
)

func CreateTargetGroup(sess *session.Session, target_group_name, target_group_protocol, target_group_protocolversion, vpcid string, target_group_port int64) *string {
        if sess == nil || target_gropu_name == "" {
                Error.Println("Not create target group. parameter not enough.")
                return nil
        }
        svc := elbv2.New(sess)
        input := &elbv2.DescribeTargetGroupsInput{
                Names: []*string{aws.String(target_group_name)},
        }

        result, err := svc.DescribeTargetGroups(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case elbv2.ErrCodeTargetGroupNotFoundException:
                                Error.Println(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
        } else {
                Info.Println("Exists Target group : ", *result.GoString())
                return result.TargetGroups[0].TargetGroupArn
        }
        // Create target group
        input := &elbv2.CreateTargetGroupInput{
                Name:            aws.String(target_group_name),
                Port:            aws.Int64(target_group_port),
                Protocol:        aws.String(target_group_protocol), //HTTP|HTTPS|TCP|TLS|UDP|GENEVE
                ProtocolVersion: aws.String(target_group_protocolversion),
                VpcId:           aws.String(vpcid),
                TargetType:      aws.String(target_type), //instance|ip|lambda
        }

        result, err := svc.CreateTargetGroup(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case elbv2.ErrCodeDuplicateTargetGroupNameException:
                                Error.Println(elbv2.ErrCodeDuplicateTargetGroupNameException, aerr.Error())
                        case elbv2.ErrCodeTooManyTargetGroupsException:
                                Error.Println(elbv2.ErrCodeTooManyTargetGroupsException, aerr.Error())
                        case elbv2.ErrCodeInvalidConfigurationRequestException:
                                Error.Println(elbv2.ErrCodeInvalidConfigurationRequestException, aerr.Error())
                        case elbv2.ErrCodeTooManyTagsException:
                                Error.Println(elbv2.ErrCodeTooManyTagsException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
        }
        Info.Println("Create target group: ", *result.TargetGroups[0].TargetGroupArn)
        return result.TargetGroups[0].TargetGroupArn
}

func CreateLB(sess *session.Session, lbname, elb_type, elb_scheme, security_group, subnets string) *string {
        if sess == nil || lbname == "" {
                Error.Println("Not create ELB. parameter not enough.")
                return nil
        }
        svc := elbv2.New(sess)
        input := &elbv2.DescribeLoadBalancersInput{
                Names: aws.StringSplit(strings.Split(elb_name, " ")),
        }
        result, err := svc.DescribeLoadBalancers(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case elbv2.ErrCodeLoadBalancerNotFoundException:
                                Error.Println(elbv2.ErrCodeLoadBalancerNotFoundException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
        } else {
                Info.Println("ELB Exists: ", *result.LoadBalancers[0].LoadBalancerArn)
                return result.LoadBalancers[0].LoadBalancerArn
        }

        input := &elbv2.CreateLoadBalancerInput{
                Name:           aws.String(lbname),
                Scheme:         aws.String(elb_scheme), //internet-facing|internal
                SecurityGroups: aws.StringSlice(strings.Split(security_group, " ")),
                Subnets:        aws.StringSlice(strings.Split(subnets, " ")),
                Type:           aws.String(elb_type), //application|network|gateway
        }

        result, err := svc.CreateLoadBalancer(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case elbv2.ErrCodeCertificateNotFoundException:
                                Error.Println(elbv2.ErrCodeCertificateNotFoundException, aerr.Error())
                        case elbv2.ErrCodeInvalidConfigurationRequestException:
                                Error.Println(elbv2.ErrCodeInvalidConfigurationRequestException, aerr.Error())
                        case elbv2.ErrCodeSubnetNotFoundException:
                                Error.Println(elbv2.ErrCodeSubnetNotFoundException, aerr.Error())
                        case elbv2.ErrCodeInvalidSubnetException:
                                Error.Println(elbv2.ErrCodeInvalidSubnetException, aerr.Error())
                        case elbv2.ErrCodeInvalidSecurityGroupException:
                                Error.Println(elbv2.ErrCodeInvalidSecurityGroupException, aerr.Error())
                        case elbv2.ErrCodeInvalidSchemeException:
                                Error.Println(elbv2.ErrCodeInvalidSchemeException, aerr.Error())
                        case elbv2.ErrCodeTooManyTagsException:
                                Error.Println(elbv2.ErrCodeTooManyTagsException, aerr.Error())
                        case elbv2.ErrCodeDuplicateTagKeysException:
                                Error.Println(elbv2.ErrCodeDuplicateTagKeysException, aerr.Error())
                        case elbv2.ErrCodeUnsupportedProtocolException:
                                Error.Println(elbv2.ErrCodeUnsupportedProtocolException, aerr.Error())
                        case elbv2.ErrCodeOperationNotPermittedException:
                                Error.Println(elbv2.ErrCodeOperationNotPermittedException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        // Print the error, cast err to awserr.Error to get the Code and
                        // Message from an error.
                        Error.Println(err.Error())
                }
        }

        Info.Println("Create ELB: ", result.LoadBalancers[0].LoadBalancerArn)
        return result.LoadBalancers[0].LoadBalancerArn
}

func CreateListener(sess *session.Session, lbname, protocol, elb_arn, targetgroup_arn string, port int64) *string{
        if sess == nil || targetgroup_name == "" || lbname == "" {
                Error.Println("Not Create Listener.")
                return nil
        }
        svc = elbv2.New(sess)
        input := &elbv2.CreateListenerInput{
                DefaultActions: []*elbv2.Action{
                        {
                                TargetGroupArn: aws.String(targetgroup_arn),
                                Type:           aws.String("forward"),
                        },
                },
                LoadBalancerArn: aws.String(elb_arn),
                Port:            aws.Int64(port),
                Protocol:        aws.String(protocol),
        }

        result, err := svc.CreateListener(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case elbv2.ErrCodeDuplicateListenerException:
                                Error.Println(elbv2.ErrCodeDuplicateListenerException, aerr.Error())
                        case elbv2.ErrCodeTooManyListenersException:
                                Error.Println(elbv2.ErrCodeTooManyListenersException, aerr.Error())
                        case elbv2.ErrCodeTooManyCertificatesException:
                                Error.Println(elbv2.ErrCodeTooManyCertificatesException, aerr.Error())
                        case elbv2.ErrCodeLoadBalancerNotFoundException:
                                Error.Println(elbv2.ErrCodeLoadBalancerNotFoundException, aerr.Error())
                        case elbv2.ErrCodeTargetGroupNotFoundException:
                                Error.Println(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
                        case elbv2.ErrCodeTargetGroupAssociationLimitException:
                                Error.Println(elbv2.ErrCodeTargetGroupAssociationLimitException, aerr.Error())
                        case elbv2.ErrCodeInvalidConfigurationRequestException:
                                Error.Println(elbv2.ErrCodeInvalidConfigurationRequestException, aerr.Error())
                        case elbv2.ErrCodeIncompatibleProtocolsException:
                                Error.Println(elbv2.ErrCodeIncompatibleProtocolsException, aerr.Error())
                        case elbv2.ErrCodeSSLPolicyNotFoundException:
                                Error.Println(elbv2.ErrCodeSSLPolicyNotFoundException, aerr.Error())
                        case elbv2.ErrCodeCertificateNotFoundException:
                                Error.Println(elbv2.ErrCodeCertificateNotFoundException, aerr.Error())
                        case elbv2.ErrCodeUnsupportedProtocolException:
                                Error.Println(elbv2.ErrCodeUnsupportedProtocolException, aerr.Error())
                        case elbv2.ErrCodeTooManyRegistrationsForTargetIdException:
                                Error.Println(elbv2.ErrCodeTooManyRegistrationsForTargetIdException, aerr.Error())
                        case elbv2.ErrCodeTooManyTargetsException:
                                Error.Println(elbv2.ErrCodeTooManyTargetsException, aerr.Error())
                        case elbv2.ErrCodeTooManyActionsException:
                                Error.Println(elbv2.ErrCodeTooManyActionsException, aerr.Error())
                        case elbv2.ErrCodeInvalidLoadBalancerActionException:
                                Error.Println(elbv2.ErrCodeInvalidLoadBalancerActionException, aerr.Error())
                        case elbv2.ErrCodeTooManyUniqueTargetGroupsPerLoadBalancerException:
                                Error.Println(elbv2.ErrCodeTooManyUniqueTargetGroupsPerLoadBalancerException, aerr.Error())
                        case elbv2.ErrCodeALPNPolicyNotSupportedException:
                                Error.Println(elbv2.ErrCodeALPNPolicyNotSupportedException, aerr.Error())
                        case elbv2.ErrCodeTooManyTagsException:
                                Error.Println(elbv2.ErrCodeTooManyTagsException, aerr.Error())
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
        Info.Println("Create Listener: ", result.Listeners[0].ListenerArn)
        return result.Listeners[0].ListenerArn
}
