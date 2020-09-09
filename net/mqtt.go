package net

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

// MqttManager mqtt obj
type MqttManager struct {
	Client mqtt.Client
	Conf   *Config
}

// Config mqtt config struct
type Config struct {
	Addr     string //tcp://iot.eclipse.org:1883
	ClientID string
	Username string
	Password string
	Log      *zerolog.Logger
}

// GetManager mqtt obj init
func GetManager(config *Config) (*MqttManager, error) {
	var c mqtt.Client
	manager := &MqttManager{}
	if manager == nil || manager.Client == nil {
		opts := mqtt.NewClientOptions().AddBroker(config.Addr)
		opts.SetClientID(config.ClientID)
		opts.SetAutoReconnect(true)
		if config.Username != "" {
			opts.SetUsername(config.Username)
		}
		if config.Password != "" {
			opts.SetPassword(config.Password)
		}
		opts.SetKeepAlive(20 * time.Second)
		opts.SetCleanSession(false)
		// opts.CleanSession = false
		//opts.SetDefaultPublishHandler(f)
		opts.SetPingTimeout(1 * time.Second)
		c = mqtt.NewClient(opts)
		// fmt.Println(opts.Servers)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			config.Log.Err(token.Error()).Interface("config", config).Msg("mqtt connect token error")
		}
	}
	if c != nil {
		config.Log.Info().Interface("config", config).Msg("mqtt connect success")
		manager.Client = c
		manager.Conf = config
		return manager, nil
	}
	return nil, fmt.Errorf("unknownErr")
}

// Publish mqtt publish func wrapper
func (m *MqttManager) Publish(topic string, qos byte, payload interface{}) error {
	if m == nil {
		m.Conf.Log.Err(fmt.Errorf("disconnectedErr")).Str("topic", topic).Msg("mqtt publish error, m nil")
		return fmt.Errorf("disconnectedErr")
	}
	buf, err := json.Marshal(payload)

	// // fmt.Println("pub msg:", string(buf))
	if err != nil {
		// 	// panic(err.Error())
		return err
	}
	if token := m.Client.Publish(topic, qos, false, buf); token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt publish error")
		return fmt.Errorf("publishErr")
	}
	// token.WaitTimeout(time.Second * 60)
	// err = token.Error()
	// if err != nil {
	// 	// panic(err)
	// 	return err
	// }
	return nil
}

// RetainPublish mqtt publish with retain func wrapper
func (m *MqttManager) RetainPublish(topic string, qos byte, payload interface{}) error {
	if m == nil {
		m.Conf.Log.Err(fmt.Errorf("disconnectedErr")).Str("topic", topic).Msg("mqtt publish error, m nil")
		return fmt.Errorf("disconnectedErr")
	}
	buf, err := json.Marshal(payload)

	// // fmt.Println("pub msg:", string(buf))
	if err != nil {
		// 	// panic(err.Error())
		return err
	}
	if token := m.Client.Publish(topic, qos, true, buf); token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt publish error")
		return fmt.Errorf("publishErr")
	}
	// token.WaitTimeout(time.Second * 60)
	// err = token.Error()
	// if err != nil {
	// 	// panic(err)
	// 	return err
	// }
	return nil
}

// Subscribe mqtt Subscribe func wrapper
func (m *MqttManager) Subscribe(topic string, qos byte, callback func(mqtt.Client, mqtt.Message)) error {
	if m == nil {
		m.Conf.Log.Err(fmt.Errorf("disconnectedErr")).Str("topic", topic).Msg("mqtt subscribe error, m nil")
		return fmt.Errorf("disconnectedErr")
	}
	if token := m.Client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt subscribe error")
		return token.Error()
	}
	return nil
}

// UnSubscribe mqtt Subscribe func wrapper
func (m *MqttManager) UnSubscribe(topic string) error {
	if m == nil {
		m.Conf.Log.Err(fmt.Errorf("disconnectedErr")).Str("topic", topic).Msg("mqtt UnSubscribe error, m nil")
		return fmt.Errorf("disconnectedErr")
	}
	if token := m.Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt Unsubscribe error")
		return token.Error()
	}
	return nil
}

// Close mqtt disconnect when system exit
func (m *MqttManager) Close(t uint) {
	m.Client.Disconnect(t)
}
