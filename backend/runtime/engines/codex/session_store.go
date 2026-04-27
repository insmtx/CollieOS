package codex

import "sync"

// SessionStore maps SingerOS session IDs to Codex thread IDs.
type SessionStore struct {
	mu      sync.RWMutex
	threads map[string]string
}

// NewSessionStore creates an in-memory Codex session store.
func NewSessionStore() *SessionStore {
	return &SessionStore{threads: make(map[string]string)}
}

// Get returns the Codex thread ID for a SingerOS session.
func (s *SessionStore) Get(sessionID string) (string, bool) {
	if s == nil || sessionID == "" {
		return "", false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	threadID, ok := s.threads[sessionID]
	return threadID, ok
}

// Set stores a Codex thread ID for a SingerOS session.
func (s *SessionStore) Set(sessionID string, threadID string) {
	if s == nil || sessionID == "" || threadID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.threads[sessionID] = threadID
}
