package notification

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const RedisChannel = "bjdms_notifications"

type NotificationType string

const (
	TypeTaskAssigned    NotificationType = "task_assigned"
	TypeJoinRequest     NotificationType = "join_request"
	TypeComplaintAlert  NotificationType = "complaint_alert"
	TypePerformanceMile NotificationType = "performance_milestone"
)

type Notification struct {
	ID             uuid.UUID        `json:"id"`
	UserID         uuid.UUID        `json:"user_id"`
	Type           NotificationType `json:"type"`
	Title          string           `json:"title"`
	Message        string           `json:"message"`
	Data           json.RawMessage  `json:"data,omitempty"`
	IsRead         bool             `json:"is_read"`
	CreatedAt      time.Time        `json:"created_at"`
	JurisdictionID uuid.UUID        `json:"jurisdiction_id"`
}

type Service struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	clients map[uuid.UUID]chan Notification
	mu      sync.RWMutex
}

func NewService(db *pgxpool.Pool, rdb *redis.Client) *Service {
	s := &Service{
		db:      db,
		redis:   rdb,
		clients: make(map[uuid.UUID]chan Notification),
	}
	
	// Start background listener for Redis Pub/Sub
	go s.listen(context.Background())
	
	return s
}

func (s *Service) Create(ctx context.Context, n *Notification) error {
	n.ID = uuid.New()
	n.CreatedAt = time.Now()
	n.IsRead = false

	query := `
		INSERT INTO notifications (id, user_id, type, title, message, data, is_read, created_at, jurisdiction_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := s.db.Exec(ctx, query,
		n.ID, n.UserID, n.Type, n.Title, n.Message, n.Data, n.IsRead, n.CreatedAt, n.JurisdictionID,
	)
	if err != nil {
		return err
	}

	// Dispatch to real-time channel if user is active
	s.dispatch(n)
	return nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID, limit int) ([]Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, data, is_read, created_at, jurisdiction_id
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := s.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Data, &n.IsRead, &n.CreatedAt, &n.JurisdictionID,
		); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (s *Service) MarkAsRead(ctx context.Context, userID, noteID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	_, err := s.db.Exec(ctx, query, noteID, userID)
	return err
}

func (s *Service) RegisterClient(userID uuid.UUID) chan Notification {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	ch := make(chan Notification, 10)
	s.clients[userID] = ch
	return ch
}

func (s *Service) UnregisterClient(userID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if ch, ok := s.clients[userID]; ok {
		close(ch)
		delete(s.clients, userID)
	}
}

func (s *Service) dispatch(n *Notification) {
	if s.redis == nil {
		// Fallback for single instance without Redis
		s.dispatchLocal(n)
		return
	}

	payload, _ := json.Marshal(n)
	s.redis.Publish(context.Background(), RedisChannel, payload)
}

func (s *Service) dispatchLocal(n *Notification) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if ch, ok := s.clients[n.UserID]; ok {
		select {
		case ch <- *n:
		default:
			// Buffer full
		}
	}
}

func (s *Service) listen(ctx context.Context) {
	if s.redis == nil {
		return
	}

	pubsub := s.redis.Subscribe(ctx, RedisChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var n Notification
		if err := json.Unmarshal([]byte(msg.Payload), &n); err != nil {
			continue
		}
		s.dispatchLocal(&n)
	}
}
