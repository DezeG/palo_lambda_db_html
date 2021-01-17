package main
 
import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"./structs"
	"context"
	"net/url"
	"strconv"


	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


func HandleLambdaEvent(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var patient structs.Patient

	m, _ := url.ParseQuery(req.Body)

	uuid, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
	}

	patient.Name = m["name"][0]
	patient.Illness = m["illness"][0]
	patient.Pain_level, _ = strconv.Atoi(m["painLevel"][0])
	patient.Hospital = m["hospital"][0]
	patient.Uuid = uuid.String()

	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type"
	resp.Headers["content-type"] = "text/html"
	resp.Body = `Success`
	resp.StatusCode = 200

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(patient)
	if err != nil {
		fmt.Println(err)
	}

	tableName := "patients"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err)
	}

	return resp, nil
}
 
func main() {
	lambda.Start(HandleLambdaEvent)
}