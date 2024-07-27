package homecaresdk

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
	SendMessage(messageType string, payload any)
	ReceiveMessage() (messageType string, payload any, err error)
}

func NewAppCommunicator(sendMsg func(string, any), recvMsg func() (string, any, error)) AppCommunicator {
	return &appCommunicator{sendMsg: sendMsg, recvMsg: recvMsg}
}

type appCommunicator struct {
	sendMsg func(string, any)
	recvMsg func() (string, any, error)
}

func (ac *appCommunicator) SendMessage(messageType string, payload any) {
	if ac.sendMsg != nil {
		ac.sendMsg(messageType, payload)
	}
}

func (ac *appCommunicator) ReceiveMessage() (string, any, error) {
	if ac.recvMsg != nil {
		return ac.recvMsg()
	}
	return "", nil, nil
}

func MarshalJSON(messageType string, payload any) (string, error) {
	msg := map[string]any{
		"type":    messageType,
		"payload": payload,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func UnmarshalJSON(data string) (messageType string, payload any, err error) {
	var msg map[string]any
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return "", nil, err
	}
	messageType, _ = msg["type"].(string)
	payload, _ = msg["payload"]
	return messageType, payload, nil
}
