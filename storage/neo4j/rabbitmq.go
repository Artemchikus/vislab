package neo4j

import (
	"context"
	"fmt"
	"strings"
	"vislab/storage"
	"vislab/storage/neo4j/types"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type (
	neo4jRabbitRepo struct {
		db neo4j.DriverWithContext
	}
)

func (n *Neo4jStorage) RabbitMQ() storage.RabbitMQRepository {
	if n.rabbitRepo != nil {
		return n.rabbitRepo
	}

	n.rabbitRepo = &neo4jRabbitRepo{db: n.db}
	return n.rabbitRepo
}

func (n *neo4jRabbitRepo) Create(ctx context.Context, rabbit *types.RabbitMQ) (string, error) {
	query := `CREATE
	(r:RabbitMQ {
		host: $host,
		port: $port
	})
	RETURN r
	`

	args := map[string]any{
		"host": rabbit.Host,
		"port": rabbit.Port,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("rabbit node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "r")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jRabbitRepo) CreateQueue(ctx context.Context, queue *types.RabbitQueue) (string, error) {
	query := `CREATE
	(rq:RabbitQueue {
		name: $name
	})
	RETURN rq
	`

	args := map[string]any{
		"name": queue.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("rabbit queue node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "rq")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jRabbitRepo) Get(ctx context.Context, host string) (*types.RabbitMQ, error) {
	query := `MATCH
	(r:RabbitMQ)
	WHERE r.host = $host
	RETURN r
	`

	args := map[string]any{
		"host": host,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("rabbit node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "r")
	if err != nil {
		return nil, err
	}

	rabbitMQ := &types.RabbitMQ{
		UID: &itemNode.ElementId,
	}

	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		rabbitMQ.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		rabbitMQ.Port = &port
	}
	if userAny, ok := itemNode.Props["user"]; ok {
		user := userAny.(string)
		rabbitMQ.User = &user
	}

	return rabbitMQ, nil
}

func (n *neo4jRabbitRepo) GetQueues(ctx context.Context, rabbitUid string) ([]*types.RabbitQueue, error) {
	query := `MATCH
	(rq:RabbitQueue)-[:IN]->(r:RabbitMQ)
	WHERE elementId(r) = $uid
	RETURN rq
	`

	args := map[string]any{
		"uid": rabbitUid,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	queues := []*types.RabbitQueue{}

	if len(res.Records) == 0 {
		return queues, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "rq")
		if err != nil {
			return nil, err
		}

		queue := &types.RabbitQueue{
			UID: &itemNode.ElementId,
		}

		if nameAny, ok := itemNode.Props["name"]; ok {
			name := nameAny.(string)
			queue.Name = &name
		}

		queues = append(queues, queue)
	}

	return queues, nil
}

func (n *neo4jRabbitRepo) Delete(ctx context.Context, uid string) error {
	query := `MATCH
	(r:RabbitMQ)
	WHERE elementId(r) = $uid
	DETACH DELETE r
	`

	args := map[string]any{
		"uid": uid,
	}

	_, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	return nil
}

func (n *neo4jRabbitRepo) DeleteQueue(ctx context.Context, uid string) error {
	query := `MATCH
	(rq:RabbitQueue)
	WHERE elementId(rq) = $uid
	DETACH DELETE rq
	`

	args := map[string]any{
		"uid": uid,
	}

	_, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	return nil
}

func (n *neo4jRabbitRepo) Update(ctx context.Context, rabbit *types.RabbitMQ) (*types.RabbitMQ, error) {
	query := `MATCH
	(r:RabbitMQ)
	WHERE r.host = $host
	SET
	`
	params := []string{}

	if rabbit.Host == nil {
		return nil, fmt.Errorf("rabbitmq cannot be updated, host field is required")
	}
	if rabbit.Port != nil {
		params = append(params, "r.port = $port")
	}
	if rabbit.User != nil {
		params = append(params, "r.user = $user")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN r"

	args := map[string]any{
		"host": rabbit.Host,
		"port": rabbit.Port,
		"user": rabbit.User,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("rabbit node not updated")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "r")
	if err != nil {
		return nil, err
	}

	newRabbitMQ := &types.RabbitMQ{
		UID: &itemNode.ElementId,
	}

	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		newRabbitMQ.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		newRabbitMQ.Port = &port
	}
	if userAny, ok := itemNode.Props["user"]; ok {
		user := userAny.(string)
		newRabbitMQ.User = &user
	}

	return newRabbitMQ, nil
}

func (n *neo4jRabbitRepo) UpdateQueue(ctx context.Context, queue *types.RabbitQueue) (*types.RabbitQueue, error) {
	query := `MATCH
	(rq:RabbitQueue)
	WHERE elementId(rq) = $uid
	SET
	`
	params := []string{}

	if queue.UID == nil {
		return nil, fmt.Errorf("rabbit queue cannot be updated, uid field is required")
	}
	if queue.Name != nil {
		params = append(params, "rq.name = $name")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN rq"

	args := map[string]any{
		"uid":  queue.UID,
		"name": queue.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("rabbit queue node not updated")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "rq")
	if err != nil {
		return nil, err
	}

	newRabbitQueue := &types.RabbitQueue{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newRabbitQueue.Name = &name
	}

	return newRabbitQueue, nil
}
