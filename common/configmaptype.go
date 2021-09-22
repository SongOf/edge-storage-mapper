package common

import "encoding/json"

// DeviceProfile is structure to store in configMap.
type DeviceProfile struct {
	DeviceInstances []DeviceInstance `json:"deviceInstances,omitempty"`
	DeviceModels    []DeviceModel    `json:"deviceModels,omitempty"`
	Protocols       []Protocol       `json:"protocols,omitempty"`
}

// DeviceInstance is structure to store device in deviceProfile.json in configmap.
type DeviceInstance struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Protocol  string `json:"protocol,omitempty"`
	PProtocol Protocol
	Model     string `json:"model,omitempty"`
	Twins     []Twin `json:"twins,omitempty"`
}

// DeviceModel is structure to store deviceModel in deviceProfile.json in configmap.
type DeviceModel struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Properties  []Property `json:"properties,omitempty"`
}

// Property is structure to store deviceModel property.
type Property struct {
	Name         string      `json:"name,omitempty"`
	DataType     string      `json:"dataType,omitempty"`
	Description  string      `json:"description,omitempty"`
	AccessMode   string      `json:"accessMode,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Minimum      int64       `json:"minimum,omitempty"`
	Maximum      int64       `json:"maximum,omitempty"`
	Unit         string      `json:"unit,omitempty"`
}

// Protocol is structure to store protocol in deviceProfile.json in configmap.
type Protocol struct {
	Name                 string          `json:"name,omitempty"`
	Protocol             string          `json:"protocol,omitempty"`
	ProtocolConfigs      json.RawMessage `json:"protocolConfig,omitempty"`
	ProtocolCommonConfig json.RawMessage `json:"protocolCommonConfig,omitempty"`
}

// Metadata is the metadata for data.
type Metadata struct {
	Timestamp string `json:"timestamp,omitempty"`
	Type      string `json:"type,omitempty"`
}

// Twin is the set/get pair to one register.
type Twin struct {
	PropertyName string       `json:"propertyName,omitempty"`
	Desired      DesiredData  `json:"desired,omitempty"`
	Reported     ReportedData `json:"reported,omitempty"`
}

// DesiredData is the desired data.
type DesiredData struct {
	Value     string   `json:"value,omitempty"`
	Metadatas Metadata `json:"metadata,omitempty"`
}

// ReportedData is the reported data.
type ReportedData struct {
	Value     string   `json:"value,omitempty"`
	Metadatas Metadata `json:"metadata,omitempty"`
}
