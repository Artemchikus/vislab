package types

type ServiceToPqTableConn struct {
	From *Service
	To   *PostgresqlTable
}

type ServiceToKubernetesJobConn struct {
	From *Service
	To   *KubernetesJob
}

type ServiceToKubernetesServiceConn struct {
	From *Service
	To   *KubernetesService
}

type ServiceToCronJobConn struct {
	From *Service
	To   *CronJob
}

type ServiceToDeploymentConn struct {
	From *Service
	To   *Deployment
}

type ServiceToIngressConn struct {
	From *Service
	To   *Ingress
}

type ServiceToPipelineJobConn struct {
	From *Service
	To   *PipelineJob
}

type ServiceToTeamConn struct {
	From *Service
	To   *Team
}

type ServiceToKafkaQueueConn struct {
	From *Service
	To   *KafkaQueue
}

type ServiceToRabbitQueueConn struct {
	From *Service
	To   *RabbitQueue
}
