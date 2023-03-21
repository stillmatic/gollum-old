package store

import (
	"errors"

	"github.com/stillmatic/gollum/conversation"
)

type Store interface {
	// Get returns the conversation for the given key.
	// If the key is not found, it returns ErrNotFound.
	Get(key string) (*conversation.Conversation, error)
	Set(key string, conversation *conversation.Conversation) error
}

var ErrNotFound = errors.New("conversation not found")

type MemoryStore struct {
	conversations map[string]*conversation.Conversation
}

func NewMemoryStore() MemoryStore {
	return MemoryStore{
		conversations: make(map[string]*conversation.Conversation),
	}
}

func (s MemoryStore) Get(key string) (*conversation.Conversation, error) {
	conversation, ok := s.conversations[key]
	if !ok {
		return conversation, ErrNotFound
	}
	return conversation, nil
}

func (s MemoryStore) Set(key string, conversation *conversation.Conversation) error {
	s.conversations[key] = conversation
	return nil
}
