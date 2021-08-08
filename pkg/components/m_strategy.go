package components

import (
	appv1 "k8s.io/api/apps/v1"
	"strings"
)

type Strategy struct {
	Type  string
	RollingUpdateConfig map[string]string
}

func (mp *Strategy) Convert() *appv1.DeploymentStrategy {
	return new(appv1.DeploymentStrategy)
}

// RollingUpdate:maxUnavailable=25%,maxSurge=25%
func ParseStrategy(strategyString string) (*Strategy, error) {
	if strategyString == "" {
		return nil, nil
	}

	strategy := new(Strategy)
	split := strings.Split(strategyString, ":")

	strategy.Type = split[0]
	strategy.RollingUpdateConfig = make(map[string]string)

	if len(split) > 1 {
		kvs := strings.Split(split[1], ",")

		for _, kv := range kvs {
			config := strings.Split(strings.Trim(kv," "), "=")
			if len(config) > 1 {
				strategy.RollingUpdateConfig[config[0]] = config[1]
			}
		}
	}

	return strategy, nil
}
