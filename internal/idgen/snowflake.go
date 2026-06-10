package idgen

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/spaolacci/murmur3"
)

var node *snowflake.Node

// init 初始化雪花 ID 节点
func init() {
	host, err := os.Hostname()
	if err != nil {
		panic(fmt.Errorf("get hostname: %w", err))
	}
	machineID := murmur3.Sum64([]byte(host))
	idNode, err := snowflake.NewNode(int64(machineID % 1023))
	if err != nil {
		panic(fmt.Errorf("new snowflake node: %w", err))
	}
	node = idNode
}

// GenUUID 获取字符串雪花 ID
func GenUUID() string {
	return node.Generate().String()
}

// GenIntUUID 获取 int64 雪花 ID
func GenIntUUID() int64 {
	return node.Generate().Int64()
}
