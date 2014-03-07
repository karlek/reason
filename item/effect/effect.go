package effect

// Item effects
const (
	Strength Type = iota + 1
	Defense  Type = iota + 1
)

type Type int
type Magnitude int

func (t Type) String() string {
	switch t {
	case Strength:
		return "Strength"
	case Defense:
		return "Defense"
	}
	return "Invalid type"
}
