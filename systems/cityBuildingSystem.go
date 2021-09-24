package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"math/rand"
	"time"
	"trafficManager/entities"
)


var Spritesheet *common.Spritesheet


var cities = [...][12]int{
	{99, 100, 101,
		454, 269, 455,
		415, 195, 416,
		452, 306, 453,
	},
	{99, 100, 101,
		268, 269, 270,
		268, 269, 270,
		305, 306, 307,
	},
	{75, 76, 77,
		446, 261, 447,
		446, 261, 447,
		444, 298, 445,
	},
	{75, 76, 77,
		407, 187, 408,
		407, 187, 408,
		444, 298, 445,
	},
	{75, 76, 77,
		186, 150, 188,
		186, 150, 188,
		297, 191, 299,
	},
	{83, 84, 85,
		413, 228, 414,
		411, 191, 412,
		448, 302, 449,
	},
	{83, 84, 85,
		227, 228, 229,
		190, 191, 192,
		301, 302, 303,
	},
	{91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 946, 947,
	},
	{91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 803, 947,
	},
	{91, 92, 93,
		238, 239, 240,
		238, 239, 240,
		312, 313, 314,
	},
}

type CityBuildingSystem struct {
	world *ecs.World
	mouseTracker entities.MouseTracker
	usedTiles []int

	elapsed, buildTime float32
	built int
}
func (c *CityBuildingSystem) New(world *ecs.World) {
	c.world = world
	c.mouseTracker.BasicEntity = ecs.NewBasic()
	c.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&c.mouseTracker.BasicEntity, &c.mouseTracker.MouseComponent, nil, nil)
		}
	}

	Spritesheet = common.NewSpritesheetWithBorderFromFile("textures/citySheet.png", 16, 16, 1, 1)


	rand.Seed(time.Now().UnixNano())

	c.updateBuildTime()

}

func (*CityBuildingSystem) Remove(entity ecs.BasicEntity){

}

func (c *CityBuildingSystem) Update(dt float32)  {
	c.elapsed += dt
	if c.elapsed >= c.buildTime {
		c.generateCity()
		c.elapsed = 0
		c.updateBuildTime()
		c.built++
	}
}



// generateCity randomly generates a city in a random location on the map
func (c *CityBuildingSystem) generateCity() {
	x := rand.Intn(18)
	y := rand.Intn(18)
	t := x + y*18

	for c.isTileUsed(t) {
		if len(c.usedTiles) > 300 {
			break //to avoid infinite loop
		}
		x = rand.Intn(18)
		y = rand.Intn(18)
		t = x + y*18
	}
	c.usedTiles = append(c.usedTiles, t)

	city := rand.Intn(len(cities))
	cityTiles := make([]*entities.City, 0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			tile := &entities.City{BasicEntity: ecs.NewBasic()}
			tile.SpaceComponent.Position = engo.Point{
				X: float32(((x+1)*64)+8) + float32(i*16),
				Y: float32(((y + 1) * 64)) + float32(j*16),
			}
			tile.RenderComponent.Drawable = Spritesheet.Cell(cities[city][i+3*j])
			tile.RenderComponent.SetZIndex(1)
			cityTiles = append(cityTiles, tile)
		}
	}

	for _, system := range c.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range cityTiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}

func (c *CityBuildingSystem) isTileUsed(tile int) bool {
	for _, t := range c.usedTiles {
		if tile == t {
			return true
		}
	}
	return false
}
func (c *CityBuildingSystem) updateBuildTime() {
	switch {
	case c.built < 2:
		// 10 to 15 seconds
		c.buildTime = 5*rand.Float32() + 10
	case c.built < 5:
		// 60 to 90 seconds
		c.buildTime = 30*rand.Float32() + 60
	case c.built < 10:
		// 30 to 90 seconds
		c.buildTime = 60*rand.Float32() + 30
	case c.built < 20:
		// 30 to 65 seconds
		c.buildTime = 35*rand.Float32() + 30
	case c.built < 25:
		// 30 to 60 seconds
		c.buildTime = 30*rand.Float32() + 30
	default:
		// 20 to 40 seconds
		c.buildTime = 20*rand.Float32() + 20
	}
}