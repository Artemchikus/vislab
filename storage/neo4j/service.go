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
	neo4jServiceRepo struct {
		db neo4j.DriverWithContext
	}
)

func (n *Neo4jStorage) Service() storage.ServiceRepository {
	if n.serviceRepo != nil {
		return n.serviceRepo
	}

	n.serviceRepo = &neo4jServiceRepo{db: n.db}
	return n.serviceRepo
}

func (n *neo4jServiceRepo) Create(ctx context.Context, service *types.Service) (string, error) {
	// query := `CREATE
	// (s:Service {
	// 	name: $name,
	// 	link: $link,
	// 	group: $group,
	// 	mainBranch: $mainBranch,
	// 	lastTag: $lastTag,
	// 	language: $language,
	// 	description: $description,
	// 	status: $status
	// })
	// RETURN s
	// `

	query := `CREATE
	(s:Service {
		name: $name,
		group: $group,
		fullName: $fullName
	})
	RETURN s
	`

	args := map[string]any{
		"name":     service.Name,
		"group":    service.Group,
		"fullName": service.FullName,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("service node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "s")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jServiceRepo) CreatePort(ctx context.Context, port *types.ServicePort) (string, error) {
	query := `CREATE
	(s:ServicePort {
		number: $number
	})
	RETURN s
	`

	args := map[string]any{
		"number": port.Number,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("service port node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "s")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jServiceRepo) Get(ctx context.Context, name string) (*types.Service, error) {
	query := `MATCH
	(s:Service)
	WHERE s.name = $name
	RETURN s
	`

	args := map[string]any{
		"name": name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("service not found: %s", name)
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "s")
	if err != nil {
		return nil, err
	}

	// service := &types.Service{
	// 	UID:         itemNode.ElementId,
	// 	Name:        itemNode.Props["name"].(string),
	// 	FullName:    itemNode.Props["fullName"].(string),
	// 	Group:       itemNode.Props["group"].(string),
	// 	Link:        itemNode.Props["link"].(string),
	// 	MainBranch:  itemNode.Props["mainBranch"].(string),
	// 	LastTag:     itemNode.Props["lastTag"].(string),
	// 	Language:    itemNode.Props["language"].(string),
	// 	Description: itemNode.Props["description"].(string),
	// 	Status:      itemNode.Props["status"].(string),
	// }

	service := &types.Service{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		service.Name = &name
	}
	if fullNameAny, ok := itemNode.Props["fullName"]; ok {
		fullName := fullNameAny.(string)
		service.FullName = &fullName
	}
	if groupAny, ok := itemNode.Props["group"]; ok {
		group := groupAny.(string)
		service.Group = &group
	}

	return service, nil
}

func (n *neo4jServiceRepo) GetPorts(ctx context.Context, uid string) ([]*types.ServicePort, error) {
	query := `MATCH
	(sp:ServicePort)-[:IN]->(s:Service)
	WHERE elementId(s) = $uid
	RETURN sp
	`

	args := map[string]any{
		"uid": uid,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	var ports []*types.ServicePort

	if len(res.Records) == 0 {
		return ports, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "sp")
		if err != nil {
			return nil, err
		}

		port := &types.ServicePort{
			UID: &itemNode.ElementId,
		}

		if numberAny, ok := itemNode.Props["number"]; ok {
			number := numberAny.(int64)
			port.Number = &number
		}

		ports = append(ports, port)
	}

	return ports, nil
}

func (n *neo4jServiceRepo) Delete(ctx context.Context, uid string) error {
	query := `MATCH
	(s:Service)
	WHERE elementId(s) = $uid
	DETACH DELETE s
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

func (n *neo4jServiceRepo) DeletePort(ctx context.Context, uid string) error {
	query := `MATCH
	(sp:ServicePort)
	WHERE elementId(sp) = $uid
	DETACH DELETE sp
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

func (n *neo4jServiceRepo) Update(ctx context.Context, service *types.Service) (*types.Service, error) {
	query := `MATCH
	(s:Service)
	WHERE s.name = $name
	SET
	`

	params := []string{}

	if service.Name == nil {
		return nil, fmt.Errorf("service cannot be updated, name is required")
	}
	if service.FullName != nil {
		params = append(params, "s.fullName = $fullName")
	}
	if service.Group != nil {
		params = append(params, "s.group = $group")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	args := map[string]any{
		"name":     service.Name,
		"fullName": service.FullName,
		"group":    service.Group,
	}

	query += strings.Join(params, ", ")
	query += " RETURN s"

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("service not updated")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "s")
	if err != nil {
		return nil, err
	}

	newService := &types.Service{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newService.Name = &name
	}
	if fullNameAny, ok := itemNode.Props["fullName"]; ok {
		fullName := fullNameAny.(string)
		newService.FullName = &fullName
	}
	if groupAny, ok := itemNode.Props["group"]; ok {
		group := groupAny.(string)
		newService.Group = &group
	}

	return newService, nil
}

func (n *neo4jServiceRepo) UpdatePort(ctx context.Context, port *types.ServicePort) (*types.ServicePort, error) {
	query := `MATCH
	(sp:ServicePort)
	WHERE elementId(sp) = $uid
	SET
	`

	params := []string{}

	if port.Number != nil {
		params = append(params, "sp.number = $number")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	args := map[string]any{
		"uid":    port.UID,
		"number": port.Number,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("service port not updated")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "sp")
	if err != nil {
		return nil, err
	}

	newPort := &types.ServicePort{
		UID: &itemNode.ElementId,
	}

	if numberAny, ok := itemNode.Props["number"]; ok {
		number := numberAny.(int64)
		newPort.Number = &number
	}

	return newPort, nil
}
