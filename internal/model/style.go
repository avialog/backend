package model

type Style string

const (
	StyleVFR Style = "VFR"
	StyleIFR Style = "IFR"
	StyleY   Style = "Y"
	StyleZ   Style = "Z"
	StyleZ2  Style = "Z2"
)

var availableStyles = []Style{
	StyleVFR,
	StyleIFR,
	StyleY,
	StyleZ,
	StyleZ2,
}
