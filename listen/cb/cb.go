// @Title
// @Description 处理回调
// @Author  55
// @Date  2022/5/16
package cb

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/zngw/zchan"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"
)

var cbTotalNum = 5 // 总回调次数

var topics = sync.Map{}

// 消息结构
type Msg struct {
	Id    string // id
	Url   string // 回调地址
	Count int    // 回调次数
	Time  int64  // 回调时间
	Data  []byte // 回调数据
}

// 初始化回调，不同次回调在不同管道
func init() {
	for i := 1; i <= cbTotalNum; i++ {
		zchan, err := zchan.New(100)
		if err == nil {
			topics.Store(i, zchan)
		}
	}

	topics.Range(func(key, value interface{}) bool {
		go func() {
			for v := range value.(*zchan.ZChan).Out {
				msg := v.(Msg)
				success := doCb(msg.Url, msg.Time, msg.Data)
				if !success && msg.Count < cbTotalNum {
					// 加入下次回调
					Push(msg.Url, msg.Count+1, msg.Data)
				}
			}
		}()
		return true
	})
}

// 添加到回调队列
// url 回调地址
// count 当前回调次数
// data 回调数据
func Push(url string, count int, data []byte) {
	msg := Msg{
		Id:    uuid.New().String(),
		Url:   url,
		Count: count,
		Time:  (int64)(math.Pow(3, float64(count)))*1000 + time.Now().UnixNano()/1e6,
		Data:  data,
	}

	queue, ok := topics.Load(count)
	if !ok {
		return
	}

	queue.(*zchan.ZChan).In <- msg
}

// 执行回调
func doCb(url string, tm int64, data []byte) (success bool) {
	wait := tm - time.Now().UnixNano()/1e6

	// 距离回调大于100ms的，直接用sleep等待，因为同次回调，后插入的肯定调用时间越后
	if wait > 100 {
		time.Sleep(time.Millisecond * time.Duration(wait))
	}

	// 回调
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return
	}

	str, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	return string(str) == "success"
}
