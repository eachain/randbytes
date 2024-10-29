# randbytes

randbytes基于sha256算法提供了一种安全的随机bytes获取，随机bytes中隐含时间及机器信息（不可逆，不能被解读出来）。本库可用于生成trace_id等。

## 示例

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/eachain/randbytes"
)

func main() {
	fmt.Println(hex.EncodeToString(randbytes.New(32)))
	fmt.Println(randbytes.UUID())
}
```
