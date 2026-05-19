package iterp

import "iter"

func Zip[T any](its ...iter.Seq[T]) iter.Seq[[]T] {
	if len(its) == 0 {
		return Empty[[]T]()
	}

	return func(yield func([]T) bool) {
		done := make(chan struct{})
		defer close(done)

		c := make([]chan T, len(its))
		for i, it := range its {
			c[i] = make(chan T)
			go func() {
				defer close(c[i])

				for v := range it {
					select {
					case c[i] <- v:
					case <-done:
						return
					}
				}
			}()
		}

		for {
			vals := make([]T, len(its))
			for i := range its {
				v, ok := <-c[i]
				if !ok {
					return
				}
				vals[i] = v
			}
			if !yield(vals) {
				return
			}
		}
	}
}

func Zip2[T any, U any](v iter.Seq[T], w iter.Seq[U]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		done := make(chan struct{})
		defer close(done)

		vc := make(chan T)
		go func() {
			defer close(vc)

			for v := range v {
				select {
				case vc <- v:
				case <-done:
					return
				}
			}
		}()

		wc := make(chan U)
		go func() {
			defer close(wc)

			for w := range w {
				select {
				case wc <- w:
				case <-done:
					return
				}
			}
		}()

		for {
			v, ok := <-vc
			if !ok {
				return
			}
			w, ok := <-wc
			if !ok {
				return
			}

			if !yield(v, w) {
				return
			}
		}
	}
}
