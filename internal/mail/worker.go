package mail

import (
	"context"
	"log"
)

type OTPJob struct {
	To       string
	Username string
	Otp      string
}

type Worker struct {
	cfg   SMTPConfig
	queue chan OTPJob
}

func NewWorker(cfg SMTPConfig, buffer int) *Worker {
	return &Worker{
		cfg:   cfg,
		queue: make(chan OTPJob, buffer),
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
				case job := <-w.queue:
				if err := SendOtp(
					w.cfg,
					job.To,
					job.Username,
					job.Otp,
				); err != nil {
					log.Println("send otp failed", err)
				}
			case <-ctx.Done():
				log.Println("emailworker stopped")
				return
			}
		}
	}()
}

func (w *Worker) Enqueue(job OTPJob) {
	w.queue <- job
}
