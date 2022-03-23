package net

import (
	"fmt"
	"runtime"
	"time"

	json "github.com/json-iterator/go"

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
	Log          *zerolog.Logger
	Addr         string
	ClientID     string
	Username     string
	Password     string
	CleanSession bool
	// Order                   bool
	// WillEnabled             bool
	// WillTopic               string
	// WillPayload             []byte
	// WillQos                 byte
	// WillRetained            bool
	// ProtocolVersion         uint
	// protocolVersionExplicit bool
	// TLSConfig               *tls.Config
	KeepAlive            time.Duration
	PingTimeout          time.Duration
	ConnectTimeout       time.Duration
	MaxReconnectInterval time.Duration
	WriteTimeout         time.Duration
	AutoReconnect        bool
	// Store                   mqtt.Store
	DefaultPublishHandler mqtt.MessageHandler
	OnConnect             mqtt.OnConnectHandler
	OnConnectionLost      mqtt.ConnectionLostHandler
	ResumeSubs            bool
}

// GetManager mqtt obj init
func GetManager(config *Config) (*MqttManager, error) {
	var c mqtt.Client
	manager := &MqttManager{}
	if manager == nil || manager.Client == nil {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(config.Addr)
		opts.SetClientID(config.ClientID)
		opts.SetUsername(config.Username)
		opts.SetPassword(config.Password)
		opts.SetCleanSession(config.CleanSession)
		if config.KeepAlive == 0 {
			config.KeepAlive = 30 * time.Second
		}
		opts.SetKeepAlive(config.KeepAlive)
		if config.PingTimeout == 0 {
			config.PingTimeout = 5 * time.Second
		}
		opts.SetPingTimeout(config.PingTimeout)
		if config.ConnectTimeout == 0 {
			config.ConnectTimeout = 20 * time.Second
		}
		opts.SetConnectTimeout(config.ConnectTimeout)
		if config.MaxReconnectInterval == 0 {
			config.MaxReconnectInterval = 10 * time.Second
		}
		opts.SetMaxReconnectInterval(config.MaxReconnectInterval)
		opts.SetWriteTimeout(config.WriteTimeout)
		opts.SetAutoReconnect(config.AutoReconnect)
		opts.SetResumeSubs(config.ResumeSubs)
		if config.DefaultPublishHandler != nil {
			opts.SetDefaultPublishHandler(config.DefaultPublishHandler)
		}
		if config.OnConnect != nil {
			opts.SetOnConnectHandler(config.OnConnect)
		}
		if config.OnConnectionLost != nil {
			opts.SetConnectionLostHandler(config.OnConnectionLost)
		}
		c = mqtt.NewClient(opts)
		// fmt.Println(opts.Servers)
		token := c.Connect()
		if token == nil {
			err := fmt.Errorf("tokenNil")
			config.Log.Err(err).Interface("config", config).Msg("mqtt connect token error")
			return nil, err
		}
		if token.Wait() && token.Error() != nil {
			config.Log.Err(token.Error()).Interface("config", config).Msg("mqtt connect token error")
			return nil, token.Error()
		}
	}
	if c != nil {
		config.Log.Debug().Interface("config", config).Msg("mqtt connect success")
		manager.Client = c
		manager.Conf = config
		return manager, nil
	}
	return nil, fmt.Errorf("unknownErr")
}

// Publish mqtt publish func wrapper
func (m *MqttManager) Publish(topic string, qos byte, payload interface{}) error {
	defer panicDefer()
	buf, err := json.Marshal(payload)

	if err != nil {
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt publish marshal error, m nil")
		return err
	}
	if m == nil || (m != nil && m.Client == nil) {
		err := fmt.Errorf("disconnectedErr")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt publish error, m nil")
		return err
	}
	token := m.Client.Publish(topic, qos, false, buf)
	if token == nil {
		err := fmt.Errorf("tokenNil")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt publish error")
		return err
	}
	if token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt publish error")
		return fmt.Errorf("publishErr")
	}
	return nil
}

func panicDefer() {
	err := recover()
	if err != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		fmt.Printf("panic: %v\n%s", err, buf)
	}
}

// RetainPublish mqtt publish with retain func wrapper
func (m *MqttManager) RetainPublish(topic string, qos byte, payload interface{}) error {
	buf, err := json.Marshal(payload)

	if err != nil {
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt RetainPublish marshal error, m nil")
		return err
	}
	if m == nil || (m != nil && m.Client == nil) {
		err := fmt.Errorf("disconnectedErr")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt RetainPublish error, m nil")
		return err
	}
	token := m.Client.Publish(topic, qos, true, buf)
	if token == nil {
		err := fmt.Errorf("tokenNil")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt RetainPublish error")
		return err
	}
	if token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt RetainPublish error")
		return fmt.Errorf("publishErr")
	}
	return nil
}

// Subscribe mqtt Subscribe func wrapper
func (m *MqttManager) Subscribe(topic string, qos byte, callback func(mqtt.Client, mqtt.Message)) error {
	if m == nil || m.Client == nil {
		err := fmt.Errorf("disconnectedErr")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt subscribe error, m nil")
		return err
	}
	token := m.Client.Subscribe(topic, qos, callback)
	if token == nil {
		err := fmt.Errorf("tokenNil")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt subscribe error")
		return err
	}
	if token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt subscribe error")
		return token.Error()
	}
	return nil
}

// SubscribeQeue mqtt Subscribe func wrapper
func (m *MqttManager) SubscribeQeue(topic string, qos byte, callback func(msgBody []byte)) error {
	if m == nil || m.Client == nil {
		err := fmt.Errorf("disconnectedErr")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt subscribe error, m nil")
		return err
	}
	token := m.Client.Subscribe(topic, qos, func(c mqtt.Client, msg mqtt.Message) {
		callback(msg.Payload())
	})
	if token == nil {
		err := fmt.Errorf("tokenNil")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt subscribe error")
		return err
	}
	if token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt subscribe error")
		return token.Error()
	}
	return nil
}

// UnSubscribe mqtt Subscribe func wrapper
func (m *MqttManager) UnSubscribe(topic string) error {
	if m == nil || m.Client == nil {
		err := fmt.Errorf("disconnectedErr")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt UnSubscribe error, m nil")
		return err
	}
	token := m.Client.Unsubscribe(topic)
	if token == nil {
		err := fmt.Errorf("tokenNil")
		m.Conf.Log.Err(err).Str("topic", topic).Msg("mqtt Unsubscribe error")
		return err
	}
	if token.Wait() && token.Error() != nil {
		m.Conf.Log.Err(token.Error()).Str("topic", topic).Msg("mqtt Unsubscribe error")
		return token.Error()
	}
	return nil
}

// Close mqtt disconnect when system exit
func (m *MqttManager) Close(t uint) {
	if m != nil && m.Client != nil {
		m.Client.Disconnect(t)
	}
}
