package logprovider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventBufferV1(t *testing.T) {
	buf := NewLogBuffer(logger.TestLogger(t), 10, 20, 1)

	buf.Enqueue(big.NewInt(1),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
	)
	buf.Enqueue(big.NewInt(2),
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x2"), LogIndex: 0},
		logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 2},
	)
	results, remaining := buf.Dequeue(int64(1), 10, 1, 2, DefaultUpkeepSelector)
	require.Equal(t, 2, len(results))
	require.Equal(t, 2, remaining)
	require.True(t, results[0].ID.Cmp(results[1].ID) != 0)
	results, remaining = buf.Dequeue(int64(1), 10, 1, 2, DefaultUpkeepSelector)
	require.Equal(t, 2, len(results))
	require.Equal(t, 0, remaining)
}

func TestLogEventBufferV1_Dequeue(t *testing.T) {
	tests := []struct {
		name         string
		logsInBuffer map[*big.Int][]logpoller.Log
		args         dequeueArgs
		lookback     int
		results      []logpoller.Log
		remaining    int
	}{
		{
			name:         "empty",
			logsInBuffer: map[*big.Int][]logpoller.Log{},
			args:         newDequeueArgs(10, 1, 1, 10, nil),
			lookback:     20,
			results:      []logpoller.Log{},
		},
		{
			name: "happy path",
			logsInBuffer: map[*big.Int][]logpoller.Log{
				big.NewInt(1): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 0},
					{BlockNumber: 14, TxHash: common.HexToHash("0x15"), LogIndex: 1},
				},
			},
			args:     newDequeueArgs(10, 5, 3, 10, nil),
			lookback: 20,
			results: []logpoller.Log{
				{}, {},
			},
		},
		{
			name: "with upkeep limits",
			logsInBuffer: map[*big.Int][]logpoller.Log{
				big.NewInt(1): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 1},
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 0},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 0},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 1},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 1},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 2},
				},
				big.NewInt(2): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 11},
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 10},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 10},
					{BlockNumber: 13, TxHash: common.HexToHash("0x13"), LogIndex: 11},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 11},
					{BlockNumber: 14, TxHash: common.HexToHash("0x14"), LogIndex: 12},
				},
			},
			args:     newDequeueArgs(10, 5, 2, 10, nil),
			lookback: 20,
			results: []logpoller.Log{
				{}, {}, {}, {},
			},
			remaining: 8,
		},
		{
			name: "with max results",
			logsInBuffer: map[*big.Int][]logpoller.Log{
				big.NewInt(1): append(createDummyLogSequence(2, 0, 12, common.HexToHash("0x12")), createDummyLogSequence(2, 0, 13, common.HexToHash("0x13"))...),
				big.NewInt(2): append(createDummyLogSequence(2, 10, 12, common.HexToHash("0x12")), createDummyLogSequence(2, 10, 13, common.HexToHash("0x13"))...),
			},
			args:     newDequeueArgs(10, 5, 3, 4, nil),
			lookback: 20,
			results: []logpoller.Log{
				{}, {}, {}, {},
			},
			remaining: 4,
		},
		{
			name: "with upkeep selector",
			logsInBuffer: map[*big.Int][]logpoller.Log{
				big.NewInt(1): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 0},
					{BlockNumber: 14, TxHash: common.HexToHash("0x15"), LogIndex: 1},
				},
			},
			args:     newDequeueArgs(10, 5, 5, 10, func(id *big.Int) bool { return false }),
			lookback: 20,
			results:  []logpoller.Log{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewLogBuffer(logger.TestLogger(t), uint(tc.lookback), uint(tc.args.blockRate), uint(tc.args.upkeepLimit))
			for id, logs := range tc.logsInBuffer {
				added, dropped := buf.Enqueue(id, logs...)
				require.Equal(t, len(logs), added+dropped)
			}
			results, remaining := buf.Dequeue(tc.args.block, tc.args.blockRate, tc.args.upkeepLimit, tc.args.maxResults, tc.args.upkeepSelector)
			require.Equal(t, len(tc.results), len(results))
			require.Equal(t, tc.remaining, remaining)
		})
	}
}

func TestLogEventBufferV1_Enqueue(t *testing.T) {
	tests := []struct {
		name                             string
		logsToAdd                        map[*big.Int][]logpoller.Log
		added, dropped                   map[string]int
		sizeOfRange                      map[*big.Int]int
		rangeStart, rangeEnd             int64
		lookback, blockRate, upkeepLimit uint
	}{
		{
			name:        "empty",
			logsToAdd:   map[*big.Int][]logpoller.Log{},
			added:       map[string]int{},
			dropped:     map[string]int{},
			sizeOfRange: map[*big.Int]int{},
			rangeStart:  0,
			rangeEnd:    10,
			blockRate:   1,
			upkeepLimit: 1,
			lookback:    20,
		},
		{
			name: "happy path",
			logsToAdd: map[*big.Int][]logpoller.Log{
				big.NewInt(1): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 0},
					{BlockNumber: 14, TxHash: common.HexToHash("0x15"), LogIndex: 1},
				},
				big.NewInt(2): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 11},
				},
			},
			added: map[string]int{
				big.NewInt(1).String(): 2,
				big.NewInt(2).String(): 1,
			},
			dropped: map[string]int{
				big.NewInt(1).String(): 0,
				big.NewInt(2).String(): 0,
			},
			sizeOfRange: map[*big.Int]int{
				big.NewInt(1): 2,
				big.NewInt(2): 1,
			},
			rangeStart:  10,
			rangeEnd:    20,
			blockRate:   5,
			upkeepLimit: 1,
			lookback:    20,
		},
		{
			name: "above limits",
			logsToAdd: map[*big.Int][]logpoller.Log{
				big.NewInt(1): createDummyLogSequence(11, 0, 12, common.HexToHash("0x12")),
				big.NewInt(2): {
					{BlockNumber: 12, TxHash: common.HexToHash("0x12"), LogIndex: 11},
				},
			},
			added: map[string]int{
				big.NewInt(1).String(): 11,
				big.NewInt(2).String(): 1,
			},
			dropped: map[string]int{
				big.NewInt(1).String(): 1,
				big.NewInt(2).String(): 0,
			},
			sizeOfRange: map[*big.Int]int{
				big.NewInt(1): 10,
				big.NewInt(2): 1,
			},
			rangeStart:  10,
			rangeEnd:    20,
			blockRate:   10,
			upkeepLimit: 1,
			lookback:    20,
		},
		{
			name: "out of block range",
			logsToAdd: map[*big.Int][]logpoller.Log{
				big.NewInt(1): append(createDummyLogSequence(2, 0, 1, common.HexToHash("0x1")), createDummyLogSequence(2, 0, 100, common.HexToHash("0x1"))...),
			},
			added: map[string]int{
				big.NewInt(1).String(): 2,
			},
			dropped: map[string]int{
				big.NewInt(1).String(): 0,
			},
			sizeOfRange: map[*big.Int]int{
				big.NewInt(1): 2,
			},
			rangeStart:  1,
			rangeEnd:    101,
			blockRate:   10,
			upkeepLimit: 10,
			lookback:    20,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewLogBuffer(logger.TestLogger(t), tc.lookback, tc.blockRate, tc.upkeepLimit)
			for id, logs := range tc.logsToAdd {
				added, dropped := buf.Enqueue(id, logs...)
				sid := id.String()
				if _, ok := tc.added[sid]; !ok {
					tc.added[sid] = 0
				}
				if _, ok := tc.dropped[sid]; !ok {
					tc.dropped[sid] = 0
				}
				require.Equal(t, tc.added[sid], added)
				require.Equal(t, tc.dropped[sid], dropped)
			}
			for id, size := range tc.sizeOfRange {
				q, ok := buf.(*logBuffer).getUpkeepQueue(id)
				require.True(t, ok)
				require.Equal(t, size, q.sizeOfRange(tc.rangeStart, tc.rangeEnd))
			}
		})
	}
}

func TestLogEventBufferV1_UpkeepQueue_clean(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q := newUpkeepLogBuffer(logger.TestLogger(t), big.NewInt(1), newLogBufferOptions(10, 1, 1))

		q.clean(10)
	})

	t.Run("happy path", func(t *testing.T) {
		buf := NewLogBuffer(logger.TestLogger(t), 10, 5, 1)

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 0},
			logpoller.Log{BlockNumber: 2, TxHash: common.HexToHash("0x1"), LogIndex: 1},
		)
		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x111"), LogIndex: 0},
			logpoller.Log{BlockNumber: 11, TxHash: common.HexToHash("0x111"), LogIndex: 1},
		)

		q, ok := buf.(*logBuffer).getUpkeepQueue(big.NewInt(1))
		require.True(t, ok)
		require.Equal(t, 4, q.sizeOfRange(1, 11))

		buf.Enqueue(big.NewInt(1),
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x171"), LogIndex: 0},
			logpoller.Log{BlockNumber: 17, TxHash: common.HexToHash("0x171"), LogIndex: 1},
		)

		require.Equal(t, 4, q.sizeOfRange(1, 18))
		require.Equal(t, 0, q.clean(12))
		require.Equal(t, 2, q.sizeOfRange(1, 18))
	})
}

type dequeueArgs struct {
	block          int64
	blockRate      int
	upkeepLimit    int
	maxResults     int
	upkeepSelector func(id *big.Int) bool
}

func newDequeueArgs(block int64, blockRate int, upkeepLimit int, maxResults int, upkeepSelector func(id *big.Int) bool) dequeueArgs {
	args := dequeueArgs{
		block:          block,
		blockRate:      blockRate,
		upkeepLimit:    upkeepLimit,
		maxResults:     maxResults,
		upkeepSelector: upkeepSelector,
	}

	if upkeepSelector == nil {
		args.upkeepSelector = DefaultUpkeepSelector
	}
	if blockRate == 0 {
		args.blockRate = 1
	}
	if maxResults == 0 {
		args.maxResults = 10
	}
	if upkeepLimit == 0 {
		args.upkeepLimit = 1
	}

	return args
}

func createDummyLogSequence(n, startIndex int, block int64, tx common.Hash) []logpoller.Log {
	logs := make([]logpoller.Log, n)
	for i := 0; i < n; i++ {
		logs[i] = logpoller.Log{
			BlockNumber: block,
			TxHash:      tx,
			LogIndex:    int64(i + startIndex),
		}
	}
	return logs
}