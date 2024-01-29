package healthchecker

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

const defaultInterval = time.Minute

type (
	CheckFn func(ctx context.Context) bool

	Checker struct {
		log      *zerolog.Logger
		cancel   context.CancelFunc
		isHealth CheckFn
		cfg      Config
	}
)

func New(l zerolog.Logger, checkFn CheckFn, cfg Config) *Checker {
	l = l.With().Str("cmp", "healthchecker").Logger()

	if cfg.Interval == 0 {
		cfg.Interval = defaultInterval
	}

	return &Checker{
		log:      &l,
		isHealth: checkFn,
		cfg:      cfg,
	}
}

func (c *Checker) Start(_ context.Context) error {
	if !c.cfg.Enabled {
		return nil
	}

	go c.run()
	return nil
}

func (c *Checker) Stop(_ context.Context) error {
	if !c.cfg.Enabled {
		return nil
	}

	c.cancel()
	return nil
}

func (c *Checker) run() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	if c.cfg.StartDelay > 0 {
		<-time.Tick(c.cfg.StartDelay) //nolint:staticcheck
	}

	ticker := time.NewTicker(c.cfg.Interval)
	defer ticker.Stop()

	c.log.Debug().Msg("checker started")

	for {
		select {
		case <-ctx.Done():
			c.log.Info().Msg("checker stopped")
			return
		case <-ticker.C:
			c.log.Debug().Msg("checking health")

			func() {
				ctx2, cancel2 := context.WithTimeout(ctx, c.cfg.Interval/2)
				defer cancel2()

				if !c.isHealth(ctx2) {
					if c.cfg.FatalOnCheck {
						c.log.Fatal().Msg("service is not healthy")
						return
					}

					c.log.Warn().Msg("service is not healthy")
				}
			}()
		}
	}
}
