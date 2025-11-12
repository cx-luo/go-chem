package core

type SessionPool struct {
	pool chan *Indigo
	size int
}

func NewSessionPool(size int) *SessionPool {
	pool := make(chan *Indigo, size)
	for i := 0; i < size; i++ {
		indigo, _ := IndigoInit()
		pool <- indigo
	}
	return &SessionPool{pool: pool, size: size}
}

func (p *SessionPool) Get() *Indigo {
	return <-p.pool
}

func (p *SessionPool) Put(indigo *Indigo) {
	select {
	case p.pool <- indigo:
	default:
		indigo.Close() // 池满时直接释放
	}
}
