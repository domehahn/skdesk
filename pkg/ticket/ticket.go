package ticket

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/domehahn/skdesk/pkg/signer"
)

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusRejected Status = "REJECTED"
)

type ReviewTicket struct {
	ID          string    `json:"id"`
	SkillName   string    `json:"skill_name"`
	SkillDigest string    `json:"skill_digest"`
	Environment string    `json:"environment"` // e.g. "production", "staging"
	Status      Status    `json:"status"`
	Reason      string    `json:"reason,omitempty"`
	Reviewer    string    `json:"reviewer,omitempty"`
	Signature   string    `json:"signature,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Manager struct {
	mu      sync.RWMutex
	tickets map[string]*ReviewTicket
}

func NewManager() *Manager {
	return &Manager{
		tickets: make(map[string]*ReviewTicket),
	}
}

func (m *Manager) Create(skillName, skillDigest, env, reason string) *ReviewTicket {
	m.mu.Lock()
	defer m.mu.Unlock()

	b := make([]byte, 8)
	_, _ = rand.Read(b)
	now := time.Now().UTC()

	t := &ReviewTicket{
		ID:          "rev_" + hex.EncodeToString(b),
		SkillName:   skillName,
		SkillDigest: skillDigest,
		Environment: env,
		Status:      StatusPending,
		Reason:      reason,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	m.tickets[t.ID] = t
	return t
}

func (m *Manager) List() []*ReviewTicket {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []*ReviewTicket
	for _, t := range m.tickets {
		list = append(list, t)
	}
	return list
}

func (m *Manager) Get(id string) (*ReviewTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.tickets[id]
	if !ok {
		return nil, errors.New("ticket not found")
	}
	return t, nil
}

func (m *Manager) Approve(id, reviewer, privKeyHex string) (*ReviewTicket, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tickets[id]
	if !ok {
		return nil, errors.New("ticket not found")
	}
	if t.Status != StatusPending {
		return nil, fmt.Errorf("cannot approve ticket with status %s", t.Status)
	}

	payload := []byte(fmt.Sprintf("%s|%s|%s|APPROVED|%s", t.ID, t.SkillDigest, t.Environment, reviewer))
	sig, err := signer.SignPayload(payload, privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}

	t.Status = StatusApproved
	t.Reviewer = reviewer
	t.Signature = sig
	t.UpdatedAt = time.Now().UTC()
	return t, nil
}

func (m *Manager) Reject(id, reviewer, reason string) (*ReviewTicket, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tickets[id]
	if !ok {
		return nil, errors.New("ticket not found")
	}
	if t.Status != StatusPending {
		return nil, fmt.Errorf("cannot reject ticket with status %s", t.Status)
	}

	t.Status = StatusRejected
	t.Reviewer = reviewer
	t.Reason = reason
	t.UpdatedAt = time.Now().UTC()
	return t, nil
}
