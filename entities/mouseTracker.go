package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}
