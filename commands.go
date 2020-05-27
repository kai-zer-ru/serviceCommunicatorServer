package serviceCommunicatorServer

type CommandStruct struct {
	Name           string            `json:"name,omitempty"`
	Description    string            `json:"description,omitempty"`
	Params         map[string]string `json:"params,omitempty"`
	Method         string            `json:"method,omitempty"`
	RequiredParams []string          `json:"required_params,omitempty"`
}
