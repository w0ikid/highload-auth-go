package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const txKey contextKey = "tx"

type ContextTransaction struct {
	pool *pgxpool.Pool
}

func NewContextTransaction(pool *pgxpool.Pool) *ContextTransaction {
	return &ContextTransaction{pool: pool}
}

func (ct *ContextTransaction) StartTransaction(ctx context.Context) (context.Context, error) {
	// Если транзакция уже есть в контексте, не создаем новую (поддержка вложенных вызовов)
	if _, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return ctx, nil
	}
	
	// Начинаем новую транзакцию
	tx, err := ct.pool.Begin(ctx)
	if err != nil {
		return ctx, err
	}
	
	// Сохраняем транзакцию в контекст
	return context.WithValue(ctx, txKey, tx), nil
}

func (ct *ContextTransaction) FinalizeTransaction(ctx context.Context, err *error) error {
	tx := RetrieveTx(ctx)
	if tx == nil {
		return nil // Транзакции нет в контексте, ничего не делаем
	}
	
	// Если во время выполнения функции произошла ошибка - делаем Rollback
	if err != nil && *err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return nil
	}
	
	// Если всё ок - делаем Commit
	return tx.Commit(ctx)
}

// RetrieveTx извлекает транзакцию из контекста
func RetrieveTx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txKey).(pgx.Tx)
	return tx
}
