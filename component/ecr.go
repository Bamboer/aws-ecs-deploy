package component

import (
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/ecr"
        "io"
        "log"
        "os"
        "strings"
)

var (
        Info  *log.Logger
        Error *log.Logger
)

func init() {
        Info = log.New(os.Stdout,
                "Info: ",
                log.Ldate|log.Ltime|log.Lshortfile)

        Error = log.New(io.MultiWriter(os.Stderr, os.Stdout),
                "Error: ",
                log.Ldate|log.Ltime|log.Lshortfile)
}

func EcrCreator(sess *session.Session, repository_name string) *string {
        if sess == nil || repository_name == "" {
                Error.Println("Not create ECR repository. Pls input repository name.")
                return nil
        }
        svc := ecr.New(sess)
        input := &ecr.DescribeRepositoriesInput{
                RepositoryNames: aws.StringSlice(strings.Split(repository_name, " ")),
        }

        result, err := svc.DescribeRepositories(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecr.ErrCodeServerException:
                                Error.Println(ecr.ErrCodeServerException, aerr.Error())
                        case ecr.ErrCodeInvalidParameterException:
                                Error.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
                        case ecr.ErrCodeRepositoryNotFoundException:
                                Error.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
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
                Info.Println("RepositoryArn: ", *result.Repositories[0].RepositoryArn, "RepositoryUri: ", *result.Repositories[0].RepositoryUri)
                return result.Repositories[0].RepositoryArn
        }
        //
        input := &ecr.CreateRepositoryInput{
                RepositoryName: aws.String(repository_name),
        }
        result, err := svc.CreateRepository(input)
        if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                        switch aerr.Code() {
                        case ecr.ErrCodeServerException:
                                Error.Println(ecr.ErrCodeServerException, aerr.Error())
                        case ecr.ErrCodeInvalidParameterException:
                                Error.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
                        case ecr.ErrCodeInvalidTagParameterException:
                                Error.Println(ecr.ErrCodeInvalidTagParameterException, aerr.Error())
                        case ecr.ErrCodeTooManyTagsException:
                                Error.Println(ecr.ErrCodeTooManyTagsException, aerr.Error())
                        case ecr.ErrCodeRepositoryAlreadyExistsException:
                                Error.Println(ecr.ErrCodeRepositoryAlreadyExistsException, aerr.Error())
                        case ecr.ErrCodeLimitExceededException:
                                Error.Println(ecr.ErrCodeLimitExceededException, aerr.Error())
                        case ecr.ErrCodeKmsException:
                                Error.Println(ecr.ErrCodeKmsException, aerr.Error())
                        default:
                                Error.Println(aerr.Error())
                        }
                } else {
                        Error.Println(err.Error())
                }
                return nil
        }

        Info.Println("RepositoryUri: ", *result.Repository.RepositoryUri, "RepositoryArn: ", *result.Repository.RepositoryArn)
        return result.Repository.RepositoryArn
}
