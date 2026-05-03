package worker

import (
	"context"
	"errors"
	"runtime"

	"github.com/w0ikid/highload-auth-go/pkg/crypto/hash"
)

var (
	ErrPoolFull = errors.New("crypto worker pool is full")
)

type jobType int

const (
	hashJob jobType = iota
	compareJob
)

type job struct {
	kind     jobType
	password string
	hash     string // для сравнения
	resultCh chan result
}

type result struct {
	hash  string
	match bool
	err   error
}

// Pool управляет фиксированным набором горутин для выполнения тяжелых криптографических задач.
type Pool struct {
	jobs chan job
}

// NewPool создает новый пул воркеров.
// Если numWorkers <= 0, используется количество ядер CPU.
func NewPool(numWorkers int, queueSize int) *Pool {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	p := &Pool{
		jobs: make(chan job, queueSize),
	}

	for i := 0; i < numWorkers; i++ {
		go p.worker()
	}

	return p
}

func (p *Pool) worker() {
	for j := range p.jobs {
		var res result
		switch j.kind {
		case hashJob:
			h, err := hash.HashPassword(j.password)
			res = result{hash: h, err: err}
		case compareJob:
			match, err := hash.ComparePassword(j.password, j.hash)
			res = result{match: match, err: err}
		}
		j.resultCh <- res
	}
}

// HashPassword ставит задачу на хэширование в очередь.
func (p *Pool) HashPassword(ctx context.Context, password string) (string, error) {
	resCh := make(chan result, 1)
	j := job{
		kind:     hashJob,
		password: password,
		resultCh: resCh,
	}

	select {
	case p.jobs <- j:
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "", ErrPoolFull
	}

	select {
	case res := <-resCh:
		return res.hash, res.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// ComparePassword ставит задачу на сравнение в очередь.
func (p *Pool) ComparePassword(ctx context.Context, password, hash string) (bool, error) {
	resCh := make(chan result, 1)
	j := job{
		kind:     compareJob,
		password: password,
		hash:     hash,
		resultCh: resCh,
	}

	select {
	case p.jobs <- j:
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		return false, ErrPoolFull
	}

	select {
	case res := <-resCh:
		return res.match, res.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}
