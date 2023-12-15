package whatsapp

func (w *Whatsapp) NewKillChannel(userID int) {
	w.killchannel[userID] = make(chan bool)
}

func (w *Whatsapp) SendKillChannel(userID int) {
	w.killchannel[userID] <- true
}
