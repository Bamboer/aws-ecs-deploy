package component

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"strings"
)

func GetVpc(sess *session.Session, vpcid string) bool {
	svc := ec2.New(sess)
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(vpcid),
		},
	}

	result, err := svc.DescribeVpcs(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				Error.Println(aerr.Error())
			}
		} else {
			Error.Println(err.Error())
		}
		return false
	}
	Info.Println(*result)
	return true
}

func GetSubnets(sess *session.Session, subnets string) bool {
	svc := ec2.New(sess)
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(strings.Split(subnets, " ")),
	}
	result, err := svc.DescribeSubnets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				Error.Println(aerr.Error())
			}
		} else {
			Error.Println(err.Error())
		}
		return false
	}
	Info.Println(*result)
	return true
}

func GetRole(sess *session.Session, role_name string) *string {
	svc := iam.New(sess)
	input := &iam.GetRoleInput{
		RoleName: aws.String(role_name),
	}

	result, err := svc.GetRole(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				Error.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				Error.Println(iam.ErrCodeServiceFailureException, aerr.Error())
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
	Info.Println(*result)
	return result.Role.Arn
}
