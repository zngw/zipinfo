// @Title
// @Description $无限缓存Channel队列
// @Author  55
// @Date  2022/5/16
package cb

import (
	"container/list"
)

type T interface{}

type Queue struct {
	In     chan<- T   // 写入channel
	Out    <-chan T   // 读取channel
	buffer *list.List // 使用list实现自适应扩缩容Buf
}

// Len uc中总共的元素数量
func (uc Queue) Len() int {
	return len(uc.In) + uc.BufLen() + len(uc.Out)
}

// BufLen uc的buf中的元素数量
func (uc Queue) BufLen() int {
	return uc.buffer.Len()
}

// New 新建一个无限缓存的Channel，并指定In和Out大小(In和Out设置得一样大)
func New(initCapacity int) *Queue {
	in := make(chan T, initCapacity)
	out := make(chan T, initCapacity)
	ch := Queue{In: in, Out: out, buffer: list.New()}

	go process(in, out, &ch)

	return &ch
}

// 内部Worker Goroutine实现
func process(in, out chan T, ch *Queue) {
	defer close(out) // in 关闭，数据读取后也把out关闭

	// 不断从in中读取数据放入到out或者buf中
loop:
	for {
		// 第一步：从in中读取数据
		value, ok := <-in
		if !ok {
			// in 关闭了，退出loop
			break loop
		}

		// 第二步：将数据存储到out或者buf中
		if ch.buffer.Len() > 0 {
			// 当buf中有数据时，新数据优先存放到buf中，确保数据FIFO原则
			ch.buffer.PushBack(value)
		} else {
			// out 没有满,数据放入out中
			select {
			case out <- value:
				continue
			default:
			}

			// out 满了，数据放入buf中
			ch.buffer.PushBack(value)
		}

		// 第三步：处理buf，一直尝试把buf中的数据放入到out中，直到buf中没有数据
		for ch.buffer.Len() > 0 {
			select {
			// 为了避免阻塞in，还要尝试从in中读取数据
			case val, ok := <-in:
				if !ok {
					// in 关闭了，退出loop
					break loop
				}
				// 因为这个时候out是满的，新数据直接放入buf中
				ch.buffer.PushBack(val)

			// 将buf中数据放入out
			case out <- ch.buffer.Front():
			}
		}
	}

	// in被关闭退出loop后，buf中还有可能有未处理的数据，将他们塞入out中，并重置buf
	for ch.buffer.Len() > 0 {
		out <- ch.buffer.Front()
	}
}
