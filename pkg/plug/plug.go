package plug

type Plug struct {
	ID         int
	Name       string
	IsOn       bool
	PowerUsage float32
	isSelected bool
	manager    *PlugManager
}

func (p *Plug) Select() {
	if p.manager != nil {
		p.manager.selectPlug(p)
	}
}

func (p *Plug) Toggle() {
	p.IsOn = !p.IsOn
}
