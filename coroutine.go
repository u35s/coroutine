package coroutine

import "time"

var waitDuration time.Duration = 4 * 1e9

func SetWaitSecond(d time.Duration) { waitDuration = d }

func NewCoroutine() *Coroutine {
	return &Coroutine{
		mYield:  make(chan struct{}),
		mResume: make(chan struct{}),
		mDone:   make(chan struct{}),

		mRunDone:      make(chan error),
		mWaitDuration: waitDuration,
	}
}

type Coroutine struct {
	mYield        chan struct{}
	mResume       chan struct{}
	mDone         chan struct{}
	mRunDone      chan error
	mWaitDuration time.Duration
}

func (co *Coroutine) SetWaitDuration(d time.Duration) {
	co.mWaitDuration = d
}

func (co *Coroutine) Run(exe func() error) error {
	go func() {
		err := exe()

		ticker := time.NewTicker(co.mWaitDuration)
		defer ticker.Stop()
		select {
		case co.mRunDone <- err:
		case <-ticker.C: // Yield 之后mRunDone <- err将会阻塞 如果用default,下面的阻塞代码可能还没执行
		}
	}()
	var err error
	select {
	case err = <-co.mRunDone:
	case <-co.mYield:
		co.mYield = nil
	}
	return err
}

func (co *Coroutine) Yield() {
	if co.mYield != nil {
		co.mYield <- struct{}{}
	}
	select {
	case <-co.mResume:
	}
}

func (co *Coroutine) Resume() {
	co.mResume <- struct{}{}
	ticker := time.NewTicker(co.mWaitDuration)
	defer ticker.Stop()
	select {
	case <-co.mDone:
	case <-ticker.C: // 当coroutine长时间不调用Done,tick.C返回导致程序并发,把tick设置的大一点
		// 当Resume后并发不会有问题可以去掉此处的select阻塞,一般不建议这么做
		// 此处tick主要防止忘记调用Done,阻塞主线程
	}
}

func (co *Coroutine) Done() {
	co.mDone <- struct{}{}
}
