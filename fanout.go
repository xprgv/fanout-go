package fanout

type Fanout[T any] struct {
	dataCh  chan T
	subs    map[chan T]struct{}
	subCh   chan chan T
	unsubCh chan chan T
	closeCh chan struct{}
}

// New creates new fanout with special data type
func New[T any]() Fanout[T] {
	f := Fanout[T]{
		dataCh:  make(chan T, 1),
		subs:    make(map[chan T]struct{}, 1),
		subCh:   make(chan chan T, 1),
		unsubCh: make(chan chan T, 1),
		closeCh: make(chan struct{}, 1),
	}

	go func() {
		for {
			select {
			case data, ok := <-f.dataCh:
				if ok {
					for sub := range f.subs {
						select {
						case sub <- data:
						default: // skip if full
						}
					}
				}
			case sub := <-f.subCh:
				f.subs[sub] = struct{}{}
			case sub := <-f.unsubCh:
				delete(f.subs, sub)
			case <-f.closeCh:
				for sub := range f.subs {
					close(sub)
					delete(f.subs, sub)
				}
				f.subs = nil
				close(f.subCh)
				close(f.unsubCh)
				close(f.dataCh)
				close(f.closeCh)
				f.closeCh = nil
				return
			}
		}
	}()

	return f
}

// Add subscriber to fanout
func (f *Fanout[T]) AddSub(sub chan T) {
	f.subCh <- sub
}

// Delete subscriber from fanout
func (f *Fanout[T]) DelSub(sub chan T) {
	f.unsubCh <- sub
}

// Publish data to all subscribers
func (f *Fanout[T]) Publish(data T) {
	f.dataCh <- data
}

// Close all. Do not use after close
func (f *Fanout[T]) Close() {
	if f.closeCh != nil {
		f.closeCh <- struct{}{}
	}
}
