package config

import "sync"

type mainData struct {
	mu       sync.Mutex
	data     Config
	fileName string
}

var main mainData

func (m *mainData) loadMainConfiguration(fileName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	config, err := LoadConfiguration(fileName)
	if err != nil {
		return err
	}
	m.data = config
	m.fileName = fileName
	return nil
}

func (m *mainData) getMainConfiguration() Config {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data
}

func (m *mainData) getMainFileName() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.fileName
}

func LoadMainConfiguration(fileName string) error {
	return main.loadMainConfiguration(fileName)
}

func GetMainConfiguration() Config {
	return main.getMainConfiguration()
}

func GetMainFileName() string {
	return main.getMainFileName()
}
