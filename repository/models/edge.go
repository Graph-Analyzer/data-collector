package models

type EdgeLabel struct {
	label string
}

var (
	EdgeLabelConnectsTo = EdgeLabel{label: "CONNECTS_TO"}
)

func (e EdgeLabel) String() string {
	return e.label
}

type Edge struct {
	EdgeID       string
	EdgeKey      string
	FromRouterID string
	ToRouterID   string
	Weight       float64
	Label        EdgeLabel
}
