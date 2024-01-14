package chix

import (
	"context"
	"errors"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// HandleCtx обработчик, контроллер
type HandleCtx func(ctx *Ctx) error
type HandlerNext func(next http.Handler) http.Handler

// Router обертка над chi роутером
type Router struct {
	*chi.Mux
	Validate Validate

	statistic map[string]int
}

// NewRouter создание роутера
func NewRouter(validate Validate) *Router {
	return &Router{
		Mux:      chi.NewRouter(),
		Validate: validate,

		statistic: make(map[string]int),
	}
}

// handleCtx запускает обработчик роутера
func (rt *Router) handlerCtx(handler HandleCtx, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ResponseWriter(ctx) == nil {
		ctx = WithResponseWriter(ctx, w)
	}

	if Request(ctx) == nil {
		ctx = WithRequest(ctx, r)
	}

	r = r.WithContext(ctx)
	ctxRoute := NewCtx(r.Context(), rt.Validate)

	if err := handler(ctxRoute); err != nil {
		err = ctxRoute.JSON(Map{
			"success": false,
			"message": err.Error(),
		})
		if err != nil {
			log.Error("Не удалось обработать запрос", slog.String("handlerCtx", err.Error()))
			http.Error(ctxRoute.ResponseWriter, "Не предвиденная ошибка", http.StatusInternalServerError)
		}
	}
}

// Get запрос на получение данных
func (rt *Router) Get(path string, handler HandleCtx) {
	rt.statistic["routes"]++
	rt.Mux.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

// Post запрос на сохранение данных
func (rt *Router) Post(path string, handler HandleCtx) {
	rt.statistic["routes"]++
	rt.Mux.Post(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

// Put запрос на обновление всех данных
func (rt *Router) Put(path string, handler HandleCtx) {
	rt.statistic["routes"]++
	rt.Mux.Put(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

// Patch запрос на обновление конкретных данных
func (rt *Router) Patch(path string, handler HandleCtx) {
	rt.statistic["routes"]++
	rt.Mux.Patch(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

// Delete запрос на удаление данных
func (rt *Router) Delete(path string, handler HandleCtx) {
	rt.statistic["routes"]++
	rt.Mux.Delete(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

// NotFound обрабатывает 404 ошибку
func (rt *Router) NotFound(handler HandleCtx) {
	rt.statistic["routes"]++
	rt.Mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

// MethodNotAllowed обрабатывает 405 ошибку
func (rt *Router) MethodNotAllowed(handler HandleCtx) {
	rt.statistic["routes"]++
	rt.Mux.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

// Use добавляет промежуточное программное обеспечение
func (rt *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	rt.statistic["middlewares"]++
	rt.Mux.Use(middlewares...)
}

// Group группирует роутеры
func (rt *Router) Group(fn func(r *Router)) Router {
	im := rt.With()

	if fn != nil {
		fn(&im)
	}

	return im
}

// With добавляет встроенное промежуточное программное обеспечение для обработчика конечной точки
func (rt *Router) With(middlewares ...func(http.Handler) http.Handler) Router {
	return Router{
		Mux:      rt.Mux.With(middlewares...).(*chi.Mux),
		Validate: rt.Validate,

		statistic: rt.statistic,
	}
}

// Route создает вложенность роутеров
func (rt *Router) Route(pattern string, fn func(r *Router)) *chi.Mux {
	subRouter := Router{
		Mux:      chi.NewRouter(),
		Validate: rt.Validate,

		statistic: rt.statistic,
	}

	fn(&subRouter)
	rt.Mux.Mount(pattern, subRouter.Mux)

	return subRouter.Mux
}

// ServeHTTP возвращает весь пул роутеров
func (rt *Router) ServeHTTP() http.HandlerFunc {
	return rt.Mux.ServeHTTP
}

// Listen запускает сервер
// Реализация: https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
func (rt *Router) Listen(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: rt.ServeHTTP(),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("Не удалось запустить сервер", slog.String("err", err.Error()))
				ch <- ctx.Err()
			}
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		panic(err)
	case <-ctx.Done():
		timeoutCtx, done := context.WithTimeout(context.Background(), time.Second*10)
		defer done()

		go func() {
			<-timeoutCtx.Done()
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				log.Error("Время корректного завершения работы истекло. Принудительный выход", slog.String("err", timeoutCtx.Err().Error()))
			}
		}()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Error("Не удалось остановить сервер", slog.String("err", err.Error()))
		}
	}

	return nil
}

// GetStatistic возвращает статистику использования роутеров
func (rt *Router) GetStatistic() map[string]int {
	return rt.statistic
}