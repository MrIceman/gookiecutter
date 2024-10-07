package dynorm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ScanOpt func(*dynamodb.ScanInput)

func ScanWithFilterExpression(filter *string) ScanOpt {
	return func(input *dynamodb.ScanInput) {
		input.FilterExpression = filter
	}
}

func ScanWithTableName(tableName *string) ScanOpt {
	return func(input *dynamodb.ScanInput) {
		input.TableName = tableName
	}
}

func ScanWithProjection(projection *string) ScanOpt {
	return func(input *dynamodb.ScanInput) {
		input.ProjectionExpression = projection
	}
}

func ScanWithExpressionAttributeNames(names map[string]string) ScanOpt {
	return func(input *dynamodb.ScanInput) {
		input.ExpressionAttributeNames = names
	}
}

func ScanWithExpressionAttributeValues(values map[string]types.AttributeValue) ScanOpt {
	return func(input *dynamodb.ScanInput) {
		input.ExpressionAttributeValues = values
	}
}

func ScanExhaustively(db *dynamodb.Client, opt ...ScanOpt) ([]map[string]types.AttributeValue, error) {
	var items []map[string]types.AttributeValue
	var lastEvaluatedKey map[string]types.AttributeValue
	for {
		params := &dynamodb.ScanInput{
			ExclusiveStartKey: lastEvaluatedKey,
		}
		for _, o := range opt {
			o(params)
		}

		res, err := db.Scan(context.TODO(), params)
		if err != nil {
			return nil, err
		}
		items = append(items, res.Items...)
		lastEvaluatedKey = res.LastEvaluatedKey
		if len(lastEvaluatedKey) == 0 {
			break
		}
	}
	return items, nil
}
