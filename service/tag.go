package service

type Tag string

func (t Tag) String() string {
	return string(t)
}

const (
	DevelopmentTag Tag = "development"
	ProductionTag  Tag = "production"
)
