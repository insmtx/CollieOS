package websocket

import (
	"sync"
)

// ConnectorInterface defines methods needed from the connector for messaging
type ConnectorInterface interface {
	GetMessenger() *Messenger
	GetAllClientIDs() []string
	SendMessageToClient(clientID string, message ServerMessage) bool
	BroadcastSend(message ServerMessage)
}

// Manager holds a reference to the connector for messaging
type Manager struct {
	mutex       sync.RWMutex
	connector   ConnectorInterface
	initialized bool
}

var defaultManager = &Manager{}

// GetManager returns the default singleton instance of manager
func GetManager() *Manager {
	return defaultManager
}

// SetConnector sets the connector instance that will be used for messaging
func (m *Manager) SetConnector(connector ConnectorInterface) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.connector = connector
	m.initialized = true
}

// IsInitialized returns true if the manager has been properly initialized
func (m *Manager) IsInitialized() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.initialized
}

// GetMessenger returns an interface for sending messages to clients, or nil if not initialized
func (m *Manager) GetMessenger() *Messenger {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if !m.initialized {
		return nil
	}

	return m.connector.GetMessenger()
}

// SendMessage broadcasts a message to a specific client or all clients
func (m *Manager) SendMessage(clientID, msgType string, payload map[string]interface{}) error {
	messager := m.GetMessenger()
	if messager == nil {
		return nil
	}

	dest := MessageDestination{ClientID: clientID}
	return messager.SendMessage(dest, msgType, payload)
}

// SendAgentStatusUpdate sends an agent status update to clients
func (m *Manager) SendAgentStatusUpdate(clientID, taskID, status, message string) error {
	messager := m.GetMessenger()
	if messager == nil {
		return nil
	}

	return messager.SendAgentStatusUpdate(clientID, taskID, status, message)
}

// SendAgentStepUpdate sends a step-by-step update from an agent during execution
func (m *Manager) SendAgentStepUpdate(clientID, taskID, step, details string) error {
	messager := m.GetMessenger()
	if messager == nil {
		return nil
	}

	return messager.SendAgentStepUpdate(clientID, taskID, step, details)
}

// SendAgentResult sends the final result of an agent's work to clients
func (m *Manager) SendAgentResult(clientID, taskID, resultType, result string) error {
	messager := m.GetMessenger()
	if messager == nil {
		return nil
	}

	return messager.SendAgentResult(clientID, taskID, resultType, result)
}

// SendLogMessage sends a detailed log message during agent execution
func (m *Manager) SendLogMessage(clientID, taskID, logLevel, message string) error {
	messager := m.GetMessenger()
	if messager == nil {
		return nil
	}

	return messager.SendLogMessage(clientID, taskID, logLevel, message)
}

// GetConnectedClients returns a list of all connected client IDs
func (m *Manager) GetConnectedClients() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if !m.initialized {
		return []string{}
	}

	return m.connector.GetAllClientIDs()
}
