package dynorm

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"log"
)

var (
	ErrItemNotFound = errors.New("not found")
)

type queryOption func(*dynamodb.QueryInput)

func WithIndex[T Entity](t T) queryOption {
	return func(input *dynamodb.QueryInput) {
		mi := t.MetaInfo()

		input.IndexName = &mi.IndexName
		expr, err := createExpression(
			expression.Key(mi.IndexKey).Equal(expression.Value(mi.IndexKeyValue)), []expression.ConditionBuilder{})
		if err != nil {
			log.Fatalf("error creating expression: %v", err)
		}
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
		input.KeyConditionExpression = expr.KeyCondition()
	}
}

func Get[T Entity](c *dynamodb.Client, result *T) (err error) {
	i := *result
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			i.MetaInfo().PartitionKey: &types.AttributeValueMemberS{Value: i.MetaInfo().PartitionkeyValue},
		},

		TableName: aws.String(i.MetaInfo().Table),
	}
	if i.MetaInfo().SortKey != "" && i.MetaInfo().SortKeyValue != "" {
		switch i.MetaInfo().SortKeyValue.(type) {
		case string:
			input.Key[i.MetaInfo().SortKey] = &types.AttributeValueMemberS{Value: i.MetaInfo().SortKeyValue.(string)}
		case int64:
			input.Key[i.MetaInfo().SortKey] = &types.AttributeValueMemberN{Value: i.MetaInfo().SortKeyValue.(string)}
		}
	}

	response, err := c.GetItem(
		context.TODO(),
		input)
	if err != nil {
		return err
	}
	if response.Item == nil {
		result = nil
		return ErrItemNotFound
	}

	err = attributevalue.UnmarshalMap(response.Item, &result)
	return err
}

func GetAllByPartitionKey[T Entity](c *dynamodb.Client, ki MetaInfo, option ...queryOption) (result []T, err error) {
	expr, err := createExpression(
		expression.Key(ki.PartitionKey).Equal(expression.Value(ki.PartitionkeyValue)),
		[]expression.ConditionBuilder{})
	if err != nil {
		return nil, err
	}
	q := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 &ki.Table,
	}
	for _, opt := range option {
		opt(q)
	}

	res, err := c.Query(context.TODO(), q)
	if err != nil {
		log.Printf("error querying table: %v", err)
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(res.Items, &result)
	return result, err
}

func Scan[T Entity](c *dynamodb.Client, ki MetaInfo) (result []T, err error) {
	res, err := c.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: &ki.Table,
	})
	if err != nil {
		log.Printf("error querying table: %v", err)
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(res.Items, &result)
	return result, err

}

func GetAllByPartitionKeyWithFilter[T Entity](c *dynamodb.Client, ki MetaInfo, expr expression.Expression) (result []T, err error) {
	res, err := c.Query(context.TODO(), &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 &ki.Table,
	})
	if err != nil {
		log.Printf("error querying table: %v", err)
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(res.Items, &result)
	return result, err
}

func Put[T Entity](c *dynamodb.Client, i T) (err error) {
	av, err := attributevalue.MarshalMap(i)
	tableName := i.MetaInfo().Table
	log.Printf("table: %s", tableName)
	_, err = c.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		})
	return err
}

func Delete[T Entity](c *dynamodb.Client, i T) (err error) {
	deleteInput := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			i.MetaInfo().PartitionKey: &types.AttributeValueMemberS{Value: i.MetaInfo().PartitionkeyValue},
		},
		TableName: aws.String(i.MetaInfo().Table),
	}
	if i.MetaInfo().SortKey != "" {
		switch i.MetaInfo().SortKeyValue.(type) {
		case string:
			deleteInput.Key[i.MetaInfo().SortKey] = &types.AttributeValueMemberS{Value: i.MetaInfo().SortKeyValue.(string)}
		case int64:
			deleteInput.Key[i.MetaInfo().SortKey] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i.MetaInfo().SortKeyValue.(int64))}
		}
	}

	_, err = c.DeleteItem(
		context.TODO(),
		deleteInput)
	return err
}

func createExpression(keyCondition expression.KeyConditionBuilder, filter []expression.ConditionBuilder) (expression.Expression, error) {
	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
	if len(filter) > 0 {
		builder = builder.WithFilter(buildAndConditions(filter))
	}
	expr, err := builder.Build()
	return expr, err
}

// buildAndConditions builds a ConditionBuilder that joins all given with AND
// if conditions is empty the builder will fail later, so make sure not to pass in an empty list
func buildAndConditions(conditions []expression.ConditionBuilder) expression.ConditionBuilder {
	switch len(conditions) {
	case 0:
		return expression.ConditionBuilder{}
	case 1:
		return conditions[0]
	default:
		return expression.And(conditions[0], conditions[1], conditions[2:]...)
	}
}
