package rcon

func (r *Rcon) OnClose(fn func(error)) {
	r.onClose = fn
}

func (r *Rcon) OnWarn(fn func(Warn)) {
	r.onWarn = fn
}

func (r *Rcon) OnKick(fn func(Kick)) {
	r.onKick = fn
}

func (r *Rcon) OnMessage(fn func(Message)) {
	r.onMessage = fn
}

func (r *Rcon) OnPosAdminCam(fn func(PosAdminCam)) {
	r.onPosAdminCam = fn
}

func (r *Rcon) OnUnposAdminCam(fn func(UnposAdminCam)) {
	r.onUnposAdminCam = fn
}

func (r *Rcon) OnSquadCreated(fn func(SquadCreated)) {
	r.onSquadCreated = fn
}

func (r *Rcon) OnListPlayers(fn func(Players)) {
	r.onListPlayers = fn
}

func (r *Rcon) OnListSquads(fn func(Squads)) {
	r.onListSquads = fn
}

func (r *Rcon) OnData(fn func(string)) {
	r.onData = fn
}
