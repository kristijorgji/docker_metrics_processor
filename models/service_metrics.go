package models

// ServiceMetrics Contains metrics about a docker service and the time
type ServiceMetrics struct {
	Datetime         string
	ContainerID      string
	ContainerName    string
	CPUPercentage    float32
	MemoryUsageInMib float32
	MemoryLimitInMib float32
	MemoryPercentage float32
}
