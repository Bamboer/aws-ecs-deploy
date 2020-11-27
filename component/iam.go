package component

import (
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/iam"
)

func GetRole(sess *session.Session, role_name string) *iam.GetRoleOutput {
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
        }
        Info.Println(result)
        return result
}
