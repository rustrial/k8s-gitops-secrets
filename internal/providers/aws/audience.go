package aws

import (
	"maps"
	"strings"
)

type K8sAudience struct {
	Generic    map[string]string
	Namespaces []string
	Names      []string
}

type AwsAudience struct {
	K8sAudience
	Partitions []string
	OrgUnits   []string
	Regions    []string
}

const (
	NamespacesKey = "namespace"
	NamesKey      = "name"
	PartitionsKey = "partition"
	OrgUnitsKey   = "orgUnits"
	RegionsKey    = "region"
)

func AwsAudienceFromMap(input map[string]string) *AwsAudience {
	audience := AwsAudience{}
	for key, value := range input {
		switch key {
		case NamespacesKey:
			audience.Namespaces = strings.Split(value, ",")
		case NamesKey:
			audience.Names = strings.Split(value, ",")
		case PartitionsKey:
			audience.Partitions = strings.Split(value, ",")
		case OrgUnitsKey:
			audience.OrgUnits = strings.Split(value, ",")
		case RegionsKey:
			audience.Regions = strings.Split(value, ",")
		default:
			if audience.Generic == nil {
				audience.Generic = make(map[string]string)
			}
			audience.Generic[key] = value
		}
	}
	return &audience
}

func (self *AwsAudience) Audience() map[string]string {
	audience := make(map[string]string)
	if self.Generic != nil {
		maps.Copy(audience, self.Generic)
	}
	if len(self.Namespaces) > 0 {
		audience[NamespacesKey] = strings.Join(self.Namespaces, ",")
	}
	if len(self.Names) > 0 {
		audience[NamesKey] = strings.Join(self.Names, ",")
	}
	if len(self.Partitions) > 0 {
		audience[PartitionsKey] = strings.Join(self.Partitions, ",")
	}
	if len(self.OrgUnits) > 0 {
		audience[OrgUnitsKey] = strings.Join(self.OrgUnits, ",")
	}
	if len(self.Regions) > 0 {
		audience[RegionsKey] = strings.Join(self.Regions, ",")
	}
	if len(audience) == 0 {
		return nil
	}
	return audience
}
