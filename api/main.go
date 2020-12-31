package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/budgets"
)

type APIResponse struct {
	Message string `json:"message"`
}

type BudgetData struct {
	BudgetName  string `json:"BudgetName"`
	LimitAmount string `json:"LimitAmount"`
	LimitUnit   string `json:"LimitUnit"`
	SpendAmount string `json:"SpendAmount"`
	SpendUnit   string `json:"SpendUnit"`
}

type Response events.APIGatewayProxyResponse

const layout  string = "2006-01-02 15:04"
var budgetsClient *budgets.Client

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "getbudget" :
			if name, ok := d["name"]; ok {
				res, e := describeBudget(ctx, name)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: res})
				}
			}
		case "getbudgets" :
			res, e := describeBudgets(ctx)
			if e != nil {
				err = e
			} else {
				jsonBytes, _ = json.Marshal(APIResponse{Message: res})
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: 500,
			Body: string(jsonBytes),
		}, nil
	}
	return Response {
		StatusCode: 200,
		Body: string(jsonBytes),
	}, nil
}

func describeBudget(ctx context.Context, budgetName string)(string, error) {
	if budgetsClient == nil {
		budgetsClient = budgets.NewFromConfig(getConfig(ctx))
	}
	input := &budgets.DescribeBudgetInput{
		AccountId: aws.String(os.Getenv("ACCOUNT_ID")),
		BudgetName: aws.String(budgetName),
	}

	result, err := budgetsClient.DescribeBudget(ctx, input)
	if err != nil {
		log.Print(err)
		return "", err
	}
	resultJson, err := json.Marshal(BudgetData{
		BudgetName: aws.ToString(result.Budget.BudgetName),
		LimitAmount: aws.ToString(result.Budget.BudgetLimit.Amount),
		LimitUnit: aws.ToString(result.Budget.BudgetLimit.Unit),
		SpendAmount: aws.ToString(result.Budget.CalculatedSpend.ActualSpend.Amount),
		SpendUnit: aws.ToString(result.Budget.CalculatedSpend.ActualSpend.Unit),
	})
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(resultJson), nil
}

func describeBudgets(ctx context.Context)(string, error) {
	if budgetsClient == nil {
		budgetsClient = budgets.NewFromConfig(getConfig(ctx))
	}
	input := &budgets.DescribeBudgetsInput{
		AccountId: aws.String(os.Getenv("ACCOUNT_ID")),
	}

	result, err := budgetsClient.DescribeBudgets(ctx, input)
	if err != nil {
		log.Print(err)
		return "", err
	}
	budgetNameList := []string{}
	for _, v := range result.Budgets {
		budgetNameList = append(budgetNameList, aws.ToString(v.BudgetName))
	}
	resultJson, err := json.Marshal(budgetNameList)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(resultJson), nil
}

func getConfig(ctx context.Context) aws.Config {
	var err error
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Print(err)
	}
	return cfg
}

func main() {
	lambda.Start(HandleRequest)
}
