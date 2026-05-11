/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/15 09:58
 */
package rlock

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisBloomFilter(t *testing.T) {
	client, _ := redismock.NewClientMock()
	bf := NewRedisBloomFilter(client)
	assert.NotNil(t, bf)
	assert.Equal(t, client, bf.client)
}

func TestAdd(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()

	tests := []struct {
		name        string
		key         string
		members     []any
		expectedErr error
	}{
		{
			name:        "Add Single Member",
			key:         "test_set",
			members:     []any{"member1"},
			expectedErr: nil,
		},
		{
			name:        "Add Multiple Members",
			key:         "test_set",
			members:     []any{"member1", "member2", "member3"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectSAdd(tt.key, tt.members...).SetVal(int64(len(tt.members)))

			bf := NewRedisBloomFilter(client)
			err := bf.Add(ctx, tt.key, tt.members...)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestIsMember(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()

	tests := []struct {
		name        string
		key         string
		member      any
		expected    bool
		expectedErr error
	}{
		{
			name:        "Member Exists",
			key:         "test_set",
			member:      "member1",
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "Member Does Not Exist",
			key:         "test_set",
			member:      "member2",
			expected:    false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectSIsMember(tt.key, tt.member).SetVal(tt.expected)

			bf := NewRedisBloomFilter(client)
			exists, err := bf.IsMember(ctx, tt.key, tt.member)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, exists)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembers(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()

	tests := []struct {
		name        string
		key         string
		expected    []string
		expectedErr error
	}{
		{
			name:        "Get Members",
			key:         "test_set",
			expected:    []string{"member1", "member2", "member3"},
			expectedErr: nil,
		},
		{
			name:        "No Members",
			key:         "test_set",
			expected:    []string{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectSMembers(tt.key).SetVal(tt.expected)

			bf := NewRedisBloomFilter(client)
			members, err := bf.Members(ctx, tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, members)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRemove(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()

	tests := []struct {
		name        string
		key         string
		members     []any
		expectedErr error
	}{
		{
			name:        "Remove Single Member",
			key:         "test_set",
			members:     []any{"member1"},
			expectedErr: nil,
		},
		{
			name:        "Remove Multiple Members",
			key:         "test_set",
			members:     []any{"member1", "member2", "member3"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectSRem(tt.key, tt.members...).SetVal(int64(len(tt.members)))

			bf := NewRedisBloomFilter(client)
			err := bf.Remove(ctx, tt.key, tt.members...)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestClear(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()

	tests := []struct {
		name        string
		key         string
		expectedErr error
	}{
		{
			name:        "Clear Set",
			key:         "test_set",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectDel(tt.key).SetVal(int64(1))

			bf := NewRedisBloomFilter(client)
			err := bf.Clear(ctx, tt.key)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRemoveMember(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()

	tests := []struct {
		name        string
		key         string
		member      any
		expectedErr error
	}{
		{
			name:        "Remove Specific Member",
			key:         "test_set",
			member:      "member1",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectSRem(tt.key, tt.member).SetVal(int64(1))

			bf := NewRedisBloomFilter(client)
			err := bf.RemoveMember(ctx, tt.key, tt.member)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLen(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()

	tests := []struct {
		name        string
		key         string
		expected    int64
		expectedErr error
	}{
		{
			name:        "Get Length",
			key:         "test_set",
			expected:    3,
			expectedErr: nil,
		},
		{
			name:        "Empty Set",
			key:         "test_set",
			expected:    0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectSCard(tt.key).SetVal(tt.expected)

			bf := NewRedisBloomFilter(client)
			length, err := bf.Len(ctx, tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, length)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
