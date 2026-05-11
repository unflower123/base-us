package snowflake

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantErr   bool
		errString string
	}{
		{
			name: "valid config",
			cfg: Config{
				Epoch:    DefaultEpoch,
				NodeBits: DefaultNodeBits,
				StepBits: DefaultStepBits,
				WorkerID: 1,
			},
			wantErr: false,
		},
		{
			name: "worker ID too low",
			cfg: Config{
				Epoch:    DefaultEpoch,
				NodeBits: DefaultNodeBits,
				StepBits: DefaultStepBits,
				WorkerID: -1,
			},
			wantErr:   true,
			errString: "worker ID must be between 0 and 7",
		},
		{
			name: "worker ID too high",
			cfg: Config{
				Epoch:    DefaultEpoch,
				NodeBits: DefaultNodeBits,
				StepBits: DefaultStepBits,
				WorkerID: 8,
			},
			wantErr:   true,
			errString: "worker ID must be between 0 and 7",
		},
		{
			name: "system time before epoch",
			cfg: Config{
				Epoch:    time.Now().UnixNano()/1e6 + 10000, // Future epoch
				NodeBits: DefaultNodeBits,
				StepBits: DefaultStepBits,
				WorkerID: 1,
			},
			wantErr:   true,
			errString: "system time is before the epoch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGenerator(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewDefaultGenerator(t *testing.T) {
	tests := []struct {
		name      string
		workerID  int64
		wantErr   bool
		errString string
	}{
		{
			name:     "valid worker ID",
			workerID: 1,
			wantErr:  false,
		},
		{
			name:      "invalid worker ID",
			workerID:  8,
			wantErr:   true,
			errString: "worker ID must be between 0 and 7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDefaultGenerator(tt.workerID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
