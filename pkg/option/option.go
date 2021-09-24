package option

type Configuration struct {
	Parameters   []Option
	Integrations []Option
}

type Option struct {
	Name        string      `json:"name"`
	Default     interface{} `json:"default"`
	Regex       string      `json:"regex"`
	Description string      `json:"description"`
	DependsOn   []string    `json:"dependsOn"`
	Files       Files       `json:"files"`
}

type Files struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}
