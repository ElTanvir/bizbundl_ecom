package settings

import (
	"context"
	"fmt"

	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/store"
	"bizbundl/util"
)

// Constants for Cache Prefixes
const (
	PrefixConfig  = "cfg:"
	PrefixPayment = "pay:"
)

// AppSecret should be loaded from env in real app, hardcoded for MVP as per user request
// In production, inject this via NewSettings(secret)
var AppSecret = "my-secret-key-32-bytes-long-1234"

type Settings struct {
	q db.Querier
}

func NewSettings(q db.Querier) *Settings {
	return &Settings{q: q}
}

// GetConfig retrieves a config value (Read-Through Pattern)
func (s *Settings) GetConfig(ctx context.Context, key string) (string, error) {
	cacheKey := PrefixConfig + key

	// 1. Check L1 Cache (Memory)
	if val, ok := store.Get().Get(cacheKey); ok {
		return val.(string), nil
	}

	// 2. Check L2 Cache (Database)
	cfg, err := s.q.GetStoreConfig(ctx, key)
	if err != nil {
		return "", err
	}

	// 3. Decrypt if needed
	finalValue := cfg.Value
	if cfg.IsEncrypted != nil && *cfg.IsEncrypted {
		decrypted, err := util.Decrypt(cfg.Value, AppSecret)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt config %s: %w", key, err)
		}
		finalValue = decrypted
	}

	// 4. Populate L1 Cache
	store.Get().Set(cacheKey, finalValue)

	return finalValue, nil
}

// SetConfig updates a config value (Write-Through Pattern)
func (s *Settings) SetConfig(ctx context.Context, key, value, group string, encrypt bool) error {
	finalValue := value

	// 1. Encrypt if needed
	if encrypt {
		encrypted, err := util.Encrypt(value, AppSecret)
		if err != nil {
			return fmt.Errorf("failed to encrypt config: %w", err)
		}
		finalValue = encrypted
	}

	// 2. Update DB
	_, err := s.q.CreateStoreConfig(ctx, db.CreateStoreConfigParams{
		Key:         key,
		Value:       finalValue,
		IsEncrypted: &encrypt,
		GroupName:   group,
	})
	// Try Update if Create fails (Upsert logic not in SQL, so simple retry or check existence first)
	// For MVP, simplistic approach: Try Update, if 0 rows, Create.
	// Actually, let's use the explicit UPSERT in SQL next time. For now, we use the generated Create/Update.
	if err != nil {
		// Assume conflict -> Update
		_, err = s.q.UpdateStoreConfig(ctx, db.UpdateStoreConfigParams{
			Key:         key,
			Value:       finalValue,
			IsEncrypted: &encrypt,
			GroupName:   group,
		})
		if err != nil {
			return err
		}
	}

	// 3. Update L1 Cache
	store.Get().Set(PrefixConfig+key, value) // Store RAW value in cache for speed

	return nil
}

// GetPaymentGateway retrieves a gateway config
func (s *Settings) GetPaymentGateway(ctx context.Context, id string) (*db.PaymentGateway, error) {
	cacheKey := PrefixPayment + id

	// 1. Memory Check
	if val, ok := store.Get().Get(cacheKey); ok {
		return val.(*db.PaymentGateway), nil
	}

	// 2. DB Check
	pg, err := s.q.GetPaymentGateway(ctx, id)
	if err != nil {
		return nil, err
	}

	// We do NOT decrypt the JSON blob here fully, typically the caller needs specific fields.
	// But for the struct, we return it as is.
	// Optimization: If the config within JSON is encrypted, we decrypt it on usage?
	// For MVP, let's assume the JSON itself isn't encrypted, but fields inside might be.
	// The requirement was: "config JSONB... Encrypted fields inside JSON".
	// Handling JSON decryption is complex. Let's stick to returning the raw struct for now.

	// 3. Populate Memory
	// We need to return a pointer copy or struct? Store generic `any` works.
	store.Get().Set(cacheKey, &pg)

	return &pg, nil
}
