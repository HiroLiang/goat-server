package agent

type ID int64

type Type string

type Status string

type Engine string

const (
	Available    Status = "AVAILABLE"
	Maintaining  Status = "MAINTAINING"
	Discontinued Status = "DISCONTINUED"
	Error        Status = "ERROR"

	Local  Type = "LOCAL"
	Remote Type = "REMOTE"
	Hybrid Type = "HYBRID"

	GGUF   Engine = "GGUF"
	ONNX   Engine = "ONNX"
	API    Engine = "API"
	Cloud  Engine = "CLOUD"
	MLC    Engine = "MLC"
	WebGPU Engine = "WebGPU"
)

func (s Status) Desc() string {
	switch s {
	case Available:
		return "available"
	case Maintaining:
		return "maintaining"
	case Discontinued:
		return "not supported"
	case Error:
		return "something went wrong"
	default:
		return "unknown"
	}
}
