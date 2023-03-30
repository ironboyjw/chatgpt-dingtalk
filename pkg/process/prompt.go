package process

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

// GeneratePrompt 生成当次请求的 Prompt
func GeneratePrompt(msg string) (rst string, err error) {
func GeneratePrompt(msg string) (rst string) {
	for _, prompt := range *public.Prompt {
		if strings.HasPrefix(msg, prompt.Title) {
			if strings.TrimSpace(msg) == prompt.Title {
				rst = fmt.Sprintf("%s：\n%s___输入内容___%s", prompt.Title, prompt.Prefix, prompt.Suffix)
				err = errors.New("消息内容为空") // 当提示词之后没有文本，抛出异常，以便直接返回Prompt所代表的内容
			} else {
				rst = prompt.Prefix + strings.TrimSpace(strings.Replace(msg, prompt.Title, "", -1)) + prompt.Suffix
			}
			rst = prompt.Content + strings.Replace(msg, prompt.Title, "", -1)
>>>>>>> parent of 71a464b (perf: 当使用prompt但内容为空时，直接返回prompt的内容 (#138))
			return
		} else {
			rst = msg
		}
	}
	return
}
