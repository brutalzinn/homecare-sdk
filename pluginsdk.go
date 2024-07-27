package pluginsdk

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
)

///PLUGIN SDK
/// LoadTemplates loads HTML templates from a directory.
/// HOMECARE 0.0.1 ALEXA - HO - PC INTEGRATION - 2021-09-01

func LoadTemplates(dir string) (*template.Template, error) {
	return template.ParseGlob(filepath.Join(dir, "*.html"))
}

func ServeAssets(mux *http.ServeMux, assetDir string, urlPrefix string) {
	mux.Handle(urlPrefix+"/", http.StripPrefix(urlPrefix, http.FileServer(http.Dir(assetDir))))
}

type AppCommunicator interface {
	SendMessage(messageType string, payload interface{})
	ReceiveMessage() (messageType string, payload interface{}, err error)
}

func NewAppCommunicator(sendMsg func(string, interface{}), recvMsg func() (string, interface{}, error)) AppCommunicator {
	return &appCommunicator{sendMsg: sendMsg, recvMsg: recvMsg}
}

type appCommunicator struct {
	sendMsg func(string, interface{})
	recvMsg func() (string, interface{}, error)
}

func (ac *appCommunicator) SendMessage(messageType string, payload interface{}) {
	if ac.sendMsg != nil {
		ac.sendMsg(messageType, payload)
	}
}

func (ac *appCommunicator) ReceiveMessage() (string, interface{}, error) {
	if ac.recvMsg != nil {
		return ac.recvMsg()
	}
	return "", nil, nil
}

func MarshalJSON(messageType string, payload interface{}) (string, error) {
	msg := map[string]interface{}{
		"type":    messageType,
		"payload": payload,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func UnmarshalJSON(data string) (messageType string, payload interface{}, err error) {
	var msg map[string]interface{}
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return "", nil, err
	}
	messageType, _ = msg["type"].(string)
	payload, _ = msg["payload"]
	return messageType, payload, nil
}
