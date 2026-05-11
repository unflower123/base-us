package snowflake

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

const (
	DefaultEpoch    = 1288834974657
	DefaultNodeBits = 3
	DefaultStepBits = 12
)

type Generator interface {
	Generate() uint64
	BatchGenerate(count int) []uint64
}

type snowflakeGenerator struct {
	node     *snowflake.Node
	nodeOnce sync.Once
	nodeErr  error
	workerID int64
}

type Config struct {
	Epoch    int64
	NodeBits uint8
	StepBits uint8
	WorkerID int64
}

func NewGenerator(cfg Config) (Generator, error) {
	maxWorkerID := int64(1<<cfg.NodeBits) - 1
	if cfg.WorkerID < 0 || cfg.WorkerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}

	snowflake.Epoch = cfg.Epoch
	snowflake.NodeBits = cfg.NodeBits
	snowflake.StepBits = cfg.StepBits

	current := time.Now().UnixNano() / 1e6
	if current < cfg.Epoch {
		return nil, fmt.Errorf("system time is before the epoch")
	}

	g := &snowflakeGenerator{
		workerID: cfg.WorkerID,
	}

	g.nodeOnce.Do(func() {
		g.node, g.nodeErr = snowflake.NewNode(g.workerID)
		fmt.Errorf("snowflake init failed: %w", g.nodeErr)
	})

	return g, nil
}

func NewDefaultGenerator(workerID int64) (Generator, error) {
	return NewGenerator(Config{
		Epoch:    DefaultEpoch,
		NodeBits: DefaultNodeBits,
		StepBits: DefaultStepBits,
		WorkerID: workerID,
	})
}

func (g *snowflakeGenerator) init() error {

	return g.nodeErr
}

func (g *snowflakeGenerator) Generate() uint64 {
	return uint64(g.node.Generate())
}

func (g *snowflakeGenerator) BatchGenerate(count int) []uint64 {
	ids := make([]uint64, 0, count)
	for i := 0; i < count; i++ {
		ids = append(ids, uint64(g.node.Generate()))
	}
	return ids
}
