package model

type ApproachType string

const (
	ApproachTypeVisual ApproachType = "VISUAL"
)

var availableApproachTypes = []ApproachType{
	ApproachTypeVisual,
}
