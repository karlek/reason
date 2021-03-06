package item

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/karlek/reason/item/effect"
	"github.com/karlek/reason/util"

	"github.com/mewkiz/pkg/errutil"
	"github.com/nsf/termbox-go"
)

// Items is a map where names of the item is the key mapping to that
// item object.
var Items = map[string]DrawItemer{}

func New(di DrawItemer) DrawItemer {
	switch i := di.(type) {
	case *Weapon:
		j := *i
		return &j
	case *Potion:
		j := *i
		return &j
	case *Tool:
		j := *i
		return &j
	case *Ring:
		j := *i
		return &j
	case *Scroll:
		j := *i
		return &j
	case *Amulet:
		j := *i
		return &j
	case *Gold:
		j := *i
		return &j
	case *Corpse:
	default:
		log.Printf("new: unknown type %T", i)
	}
	return nil
}

// Load initializes the Items map with creatures.
func Load() (err error) {
	filenames, err := util.DirFiles("github.com/karlek/reason/item/data/")
	if err != nil {
		return errutil.Err(err)
	}
	for _, filename := range filenames {
		i, err := load(filename)
		if err != nil {
			return errutil.Err(err)
		}
		if i == nil {
			log.Println("asdf")
			continue
		}
		Items[i.Name()] = i
	}
	return nil
}

// load parses a JSON data file into a go item object.
func load(filename string) (i DrawItemer, err error) {
	// jsonItem is a temporary struct for easier conversion between JSON and
	// go structs.
	type jsonItem struct {
		Name     string
		Category string
		Flavor   string
		Use      string
		Rarity   string
		Num      int
		Graphics struct {
			Ch string
			Fg map[string]string
			Bg map[string]string
		}
		Effects []struct {
			Type      string
			Magnitude int
		}
		Pathable bool
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errutil.Err(err)
	}

	var ji jsonItem
	err = json.Unmarshal(buf, &ji)
	if err != nil {
		return nil, errutil.Err(err)
	}

	fg, err := util.ParseColor(ji.Graphics.Fg)
	if err != nil {
		return nil, errutil.Err(err)
	}
	bg, err := util.ParseColor(ji.Graphics.Bg)
	if err != nil {
		return nil, errutil.Err(err)
	}

	rarity, err := parseRarity(ji.Rarity)
	if err != nil {
		return nil, errutil.Err(err)
	}

	effs, err := parseEffects(ji.Effects)
	if err != nil {
		return nil, errutil.Err(err)
	}
	log.Printf("%#v\n", effs)
	j := Item{
		name:     ji.Name,
		flavor:   ji.Flavor,
		use:      ji.Use,
		rarity:   rarity,
		count:    ji.Num,
		category: ji.Category,
		effects:  effs,
	}
	j.SetPathable(ji.Pathable)
	j.SetGraphics(termbox.Cell{
		Ch: rune(ji.Graphics.Ch[0]),
		Fg: fg,
		Bg: bg,
	})
	switch ji.Category {
	case "weapon":
		i = &Weapon{Item: j}
	case "potion":
		i = &Potion{Item: j}
	case "tool":
		i = &Tool{Item: j}
	case "ring":
		i = &Ring{Item: j}
	case "corpse":
		i = &Corpse{Item: j}
	case "amulet":
		i = &Amulet{Item: j}
	case "gold":
		i = &Gold{Item: j}
	case "scroll":
		i = &Scroll{Item: j}
	default:
		log.Fatalln("implement %s", ji.Category)
	}

	return i, nil
}

func parseRarity(rarity string) (int, error) {
	switch rarity {
	case "common":
		return Common, nil
	case "magical":
		return Magical, nil
	case "artifact":
		return Artifact, nil
	default:
		return 0, errutil.NewNoPosf("invalid rarity: %s", rarity)
	}
}

func parseEffects(jsEffects []struct {
	Type      string
	Magnitude int
}) (map[effect.Type]effect.Magnitude, error) {
	itemEffects := make(map[effect.Type]effect.Magnitude)
	for _, eff := range jsEffects {
		var t effect.Type
		switch eff.Type {
		case "Strength":
			t = effect.Strength
		case "Defense":
			t = effect.Defense
		default:
			return nil, errutil.Newf("invalid type: %s", eff.Type)
		}
		itemEffects[t] = effect.Magnitude(eff.Magnitude)
	}
	return itemEffects, nil
}
