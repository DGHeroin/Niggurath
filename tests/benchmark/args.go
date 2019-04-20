package benchmark

import "uTu/utils"

var (
    counter                 = uint64(0)
    isRunning               = true
    testTimes               = 5
    rpcSenderCoroutineCount = 500
)

func size(s uint64) string {
    return utils.ByteSize(s)
}