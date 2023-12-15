package whatsapp

// webhook for regular messages
func (w *Whatsapp) CallHook(myurl string, payload map[string]string, id int) {
	w.log.Info().Str("url", myurl).Msg("Sending POST")
	_, err := w.clientHttp[id].R().SetFormData(payload).Post(myurl)

	if err != nil {
		w.log.Debug().Str("error", err.Error())
	}
}

// webhook for messages with file attachments
func (w *Whatsapp) CallHookFile(myurl string, payload map[string]string, id int, file string) {
	w.log.Info().Str("file", file).Str("url", myurl).Msg("Sending POST")
	w.clientHttp[id].R().SetFiles(map[string]string{
		"file": file,
	}).SetFormData(payload).Post(myurl)
}
