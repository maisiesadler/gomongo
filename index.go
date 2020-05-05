package gomongo

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Connect connects to the database using parameter from Systems Manager parameter store
func Connect(ctx context.Context) bool {
	if connected() {
		return true
	}

	connectionString, err := getConnectionString()

	if err != nil {
		fmt.Printf("Error getting connection string, %v.\n", err)
		return false
	}
	fmt.Printf("Got connection string: len=%v\n", len(connectionString))

	fmt.Printf("Connecting...\n")
	connected := connect(ctx, connectionString)
	fmt.Printf("Connected - %v.\n", connected)

	return true
}

func getConnectionString() (string, error) {
	override := os.Getenv("connectionstring")
	if len(override) > 0 {
		return override, nil
	}

	region := "eu-west-2"
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return "", err
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(region))
	keyname := "mongodb"
	withDecryption := true

	fmt.Println("getting parameter")

	paramOutput, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})

	fmt.Println("request complete")

	if err != nil {
		return "", err
	}

	return *paramOutput.Parameter.Value, nil
}
