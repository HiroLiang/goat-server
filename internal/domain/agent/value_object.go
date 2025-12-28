package agent

type ID int64

type Type string

type Status string

type Engine string

const (
	Available    Status = "available"
	Maintaining  Status = "maintaining"
	Discontinued Status = "discontinued"
	Error        Status = "error"

	Local  Type = "local"
	Remote Type = "remote"
	Hybrid Type = "hybrid"

	GGUF   Engine = "gguf"
	ONNX   Engine = "onnx"
	API    Engine = "api"
	Cloud  Engine = "cloud"
	MLC    Engine = "mlc"
	WebGPU Engine = "webGPU"
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
