package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"image"
	"image/color"
	"trafficManager/entities"
	"trafficManager/systems"
)

type mainScene struct {

}

func (s mainScene) Type() string {
	return "Traffic Manager"
}
// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (s mainScene) Preload() {
	err := engo.Files.Load("textures/citySheet.png", "tilemap/TrafficMap.tmx")
	if err != nil {
		panic(err)
	}
}
// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (s mainScene) Setup(u engo.Updater)  {
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})
	kbs := common.NewKeyboardScroller(
		300, engo.DefaultHorizontalAxis,engo.DefaultVerticalAxis,
		)
	world.AddSystem(kbs)
	world.AddSystem(&common.EdgeScroller{ScrollSpeed: 300, EdgeMargin: 20})
	world.AddSystem(&common.MouseZoomer{ZoomSpeed: -0.125})

	world.AddSystem(&systems.CityBuildingSystem{})

	hud := entities.HUD{BasicEntity: ecs.NewBasic()}
	hud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0,engo.WindowHeight() - 200},
		Width: 200,
		Height: 200,
	}
	hudImage := image.NewUniform(color.RGBA{R: 205, G: 205, B: 205, A: 255})
	hudNRGBA := common.ImageToNRGBA(hudImage, 200, 200)
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)

	hud.RenderComponent = common.RenderComponent{
		Drawable: hudTexture,
		Scale:    engo.Point{X: 1, Y: 1},
		Repeat:   common.Repeat,
	}

	hud.RenderComponent.SetShader(common.HUDShader)
	hud.RenderComponent.SetZIndex(1)
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}

	resource, err := engo.Files.Resource("tilemap/TrafficMap.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	tiles := make([]*entities.Tile, 0)
	for _,tileLayer := range levelData.TileLayers{
		for _, tileElement := range tileLayer.Tiles{
			if tileElement.Image != nil {
				tile := &entities.Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: tileElement.Image,
					Scale:    engo.Point{1, 1},
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    0,
					Height:   0,
				}
				tiles = append(tiles, tile)
			}
		}

	}
	// add the tiles to the RenderSystem
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range tiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
	common.CameraBounds = levelData.Bounds()
}

func main()  {
	opts := engo.RunOptions{
		Title: "Traffic Manager",
		Width: 800,
		Height: 450,
		StandardInputs: true,
	}
	engo.Run(opts, &mainScene{})
}
