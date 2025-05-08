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
	neo4jKafkaRepo struct {
		db neo4j.DriverWithContext
	}
)

func (n *Neo4jStorage) Kafka() storage.KafkaRepository {
	if n.kafkaRepo != nil {
		return n.kafkaRepo
	}

	n.kafkaRepo = &neo4jKafkaRepo{db: n.db}
	return n.kafkaRepo
}

func (n *neo4jKafkaRepo) Create(ctx context.Context, kafka *types.Kafka) (string, error) {
	query := `CREATE
	(k:Kafka {
		host: $host,
		port: $port
	})
	RETURN k
	`

	args := map[string]any{
		"host": kafka.Host,
		"port": kafka.Port,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("kafka node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "k")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jKafkaRepo) CreateQueue(ctx context.Context, queue *types.KafkaQueue) (string, error) {
	query := `CREATE
	(kq:KafkaQueue {
		name: $name,
		queueType: $queueType,
		topic: $topic,
		typeName: $typeName
	})
	RETURN kq
	`

	args := map[string]any{
		"name":      queue.Name,
		"queueType": queue.QueueType,
		"topic":     queue.Topic,
		"typeName":  queue.TypeName,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("kafka queue node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "kq")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jKafkaRepo) Delete(ctx context.Context, uid string) error {
	query := `MATCH
	(k:Kafka)
	WHERE elementId(k) = $uid
	DETACH DELETE k
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

func (n *neo4jKafkaRepo) DeleteQueue(ctx context.Context, uid string) error {
	query := `MATCH
	(kq:KafkaQueue)
	WHERE elementId(kq) = $uid
	DETACH DELETE kq
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

func (n *neo4jKafkaRepo) Get(ctx context.Context, host string) (*types.Kafka, error) {
	query := `MATCH
	(k:Kafka)
	WHERE k.host = $host
	RETURN k
	`

	args := map[string]any{
		"host": host,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("kafka node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "k")
	if err != nil {
		return nil, err
	}

	kafka := &types.Kafka{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		kafka.Name = &name
	}
	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		kafka.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		kafka.Port = &port
	}

	return kafka, nil
}

func (n *neo4jKafkaRepo) GetQueues(ctx context.Context, kafkaUid string) ([]*types.KafkaQueue, error) {
	query := `MATCH
	(kq:KafkaQueue)-[:IN]->(k:Kafka)
	WHERE elementId(k) = $uid
	RETURN kq
	`

	args := map[string]any{
		"uid": kafkaUid,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	var queues []*types.KafkaQueue

	if len(res.Records) == 0 {
		return queues, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "kq")
		if err != nil {
			return nil, err
		}

		queue := &types.KafkaQueue{
			UID: &itemNode.ElementId,
		}

		if nameAny, ok := itemNode.Props["name"]; ok {
			name := nameAny.(string)
			queue.Name = &name
		}
		if queueTypeAny, ok := itemNode.Props["queueType"]; ok {
			queueType := queueTypeAny.(string)
			queue.QueueType = &queueType
		}
		if topicAny, ok := itemNode.Props["topic"]; ok {
			topic := topicAny.(string)
			queue.Topic = &topic
		}
		if typeNameAny, ok := itemNode.Props["typeName"]; ok {
			typeName := typeNameAny.(string)
			queue.TypeName = &typeName
		}

		queues = append(queues, queue)
	}

	return queues, nil
}

func (n *neo4jKafkaRepo) Update(ctx context.Context, kafka *types.Kafka) (*types.Kafka, error) {
	query := `MATCH
	(k:Kafka)
	WHERE k.host = $host
	SET
	`
	params := []string{}

	if kafka.Host == nil {
		return nil, fmt.Errorf("kafka cannot be updated, host field is required")
	}
	if kafka.Name != nil {
		params = append(params, "k.name = $name")
	}
	if kafka.Port != nil {
		params = append(params, "k.port = $port")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN k"

	args := map[string]any{
		"host": kafka.Host,
		"name": kafka.Name,
		"port": kafka.Port,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("kafka node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "k")
	if err != nil {
		return nil, err
	}

	newKafka := &types.Kafka{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newKafka.Name = &name
	}
	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		newKafka.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		newKafka.Port = &port
	}

	return newKafka, nil
}

func (n *neo4jKafkaRepo) UpdateQueue(ctx context.Context, queue *types.KafkaQueue) (*types.KafkaQueue, error) {
	query := `MATCH
	(kq:KafkaQueue)
	WHERE elementId(kq) = $uid
	SET
	`
	params := []string{}

	if queue.UID == nil {
		return nil, fmt.Errorf("kafka queue cannot be updated, uid field is required")
	}
	if queue.Name != nil {
		params = append(params, "kq.name = $name")
	}
	if queue.QueueType != nil {
		params = append(params, "kq.queueType = $queueType")
	}
	if queue.Topic != nil {
		params = append(params, "kq.topic = $topic")
	}
	if queue.TypeName != nil {
		params = append(params, "kq.typeName = $typeName")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN kq"

	args := map[string]any{
		"uid":       queue.UID,
		"name":      queue.Name,
		"queueType": queue.QueueType,
		"topic":     queue.Topic,
		"typeName":  queue.TypeName,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("kafka queue node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "kq")
	if err != nil {
		return nil, err
	}

	newQueue := &types.KafkaQueue{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newQueue.Name = &name
	}
	if queueTypeAny, ok := itemNode.Props["queueType"]; ok {
		queueType := queueTypeAny.(string)
		newQueue.QueueType = &queueType
	}
	if topicAny, ok := itemNode.Props["topic"]; ok {
		topic := topicAny.(string)
		newQueue.Topic = &topic
	}
	if typeNameAny, ok := itemNode.Props["typeName"]; ok {
		typeName := typeNameAny.(string)
		newQueue.TypeName = &typeName
	}

	return newQueue, nil
}
