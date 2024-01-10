package sparql

type Head struct {
	Vars []string `json:"vars"`
}

type Value struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Results struct {
	Bindings []map[string]Value `json:"bindings"`
}

type Result struct {
	Head    Head    `json:"head"`
	Results Results `json:"results"`
}
