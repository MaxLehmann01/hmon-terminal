package plug

type PlugManager struct {
	plugs             []*Plug
	selectedPlugIndex int
}

func NewPlugManager() *PlugManager {
	return &PlugManager{}
}

func (pm *PlugManager) AddPlug(p *Plug) {
	p.manager = pm
	pm.plugs = append(pm.plugs, p)

	if len(pm.plugs) == 1 {
		pm.selectPlug(p)
	}
}

func (pm *PlugManager) selectPlug(p *Plug) {
	for i, plug := range pm.plugs {
		plug.isSelected = (plug == p)
		if plug.isSelected {
			pm.selectedPlugIndex = i
		}
	}
}

func (pm *PlugManager) SelectNext() {
	if len(pm.plugs) == 0 {
		return
	}
	pm.selectedPlugIndex = (pm.selectedPlugIndex + 1) % len(pm.plugs)
	pm.selectPlug(pm.plugs[pm.selectedPlugIndex])
}

func (pm *PlugManager) SelectPlugByID(id int) {
	for _, p := range pm.plugs {
		if p.ID == id {
			pm.selectPlug(p)
			return
		}
	}
}

func (pm *PlugManager) GetSelected() *Plug {
	if len(pm.plugs) == 0 {
		return nil
	}

	return pm.plugs[pm.selectedPlugIndex]
}

func (pm *PlugManager) ToggleSelected() {
	if p := pm.GetSelected(); p != nil {
		p.Toggle()
	}
}

func (pm *PlugManager) Clear() {
	pm.plugs = []*Plug{}
	pm.selectedPlugIndex = 0
}
