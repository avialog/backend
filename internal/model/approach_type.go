package model

type ApproachType string

const (
	ApproachTypeVisual ApproachType = "VISUAL"
)

var AvailableApproachTypes = []ApproachType{
	ApproachTypeVisual,
}
