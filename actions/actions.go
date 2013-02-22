package actions

import "github.com/nsf/termbox-go"
import "github.com/karlek/worc/terrain"
import "github.com/karlek/worc/creature"
import "github.com/karlek/worc/menu"
import "github.com/karlek/reason/environment"

func Look(a *terrain.Area, hero creature.Creature, as menu.AreaScreen) (err error) {
	x, y := hero.X, hero.Y

lookLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			// Movement
			case termbox.KeyEsc:
				termbox.HideCursor()
				break lookLoop
			case termbox.KeyArrowUp:
				x, y = x, y-1
			case termbox.KeyArrowDown:
				x, y = x, y+1
			case termbox.KeyArrowLeft:
				x, y = x-1, y
			case termbox.KeyArrowRight:
				x, y = x+1, y
			}
			// Prevent from moving outside the top
			if y < 0 {
				y = 0
			}

			// Prevent from moving outside the bottom
			if y == a.Height {
				y = a.Height - 1
			}

			// Prevent from moving outside the left
			if x < 0 {
				x = 0
			}

			// Prevent from moving outside the right
			if x == a.Width {
				x = a.Width - 1
			}

			termbox.SetCursor(x, y)
			if object, found := a.Objects[terrain.Coord{X: x + as.XOffset, Y: y + as.YOffset}]; found {
				switch obj := object.(type) {
				case creature.Creature:
					menu.PrintStatus("You see: " + obj.Name)
				default:
					menu.PrintStatus("You wat")
				}
			} else {
				menu.PrintStatus("You see: " + a.Terrain[y+as.YOffset][x+as.XOffset].Name)
			}
		}
	}

	termbox.Flush()
	return nil
}

func Open(a *terrain.Area, hero creature.Creature, as menu.AreaScreen) (err error) {
	menu.PrintStatus("Open/Close door - In which direction lies the door?")
	var x, y int
openLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			// Movement
			case termbox.KeyArrowUp:
				x, y = hero.X, hero.Y-1
			case termbox.KeyArrowDown:
				x, y = hero.X, hero.Y+1
			case termbox.KeyArrowLeft:
				x, y = hero.X-1, hero.Y
			case termbox.KeyArrowRight:
				x, y = hero.X+1, hero.Y
			}
			// Prevent from moving outside the top
			if y < 0 {
				y = 0
			}

			// Prevent from moving outside the bottom
			if y == a.Height {
				y = a.Height - 1
			}

			// Prevent from moving outside the left
			if x < 0 {
				x = 0
			}

			// Prevent from moving outside the right
			if x == a.Width {
				x = a.Width - 1
			}

			if a.Terrain[y][x] == environment.ClosedDoor {
				a.Terrain[y][x] = environment.OpenDoor
				a.UpdateCoord(x, y, as)
			} else if a.Terrain[y][x] == environment.OpenDoor {
				a.Terrain[y][x] = environment.ClosedDoor
				a.UpdateCoord(x, y, as)
			} else {
				menu.PrintStatus("You can't open/close that.")
			}
			break openLoop
		}
	}

	termbox.Flush()
	return nil
}

func WalkedIntoDoor(a *terrain.Area, as menu.AreaScreen, x, y int) {
doorLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'n':
				break doorLoop
			case 'y':
				a.Terrain[y][x] = environment.OpenDoor
				a.Draw(as)
				break doorLoop
			}

			switch ev.Key {
			case termbox.KeyEsc:
				break doorLoop
			case termbox.KeyEnter:
				a.Terrain[y][x] = environment.OpenDoor
				a.Draw(as)
				break doorLoop
			}
		}
	}
}
