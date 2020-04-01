package app

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	h2 "net/http"

	"github.com/d7561985/questions/delivery/http"
	"github.com/d7561985/questions/internal/tr"
	"github.com/d7561985/questions/repository/filerepo"
	"github.com/d7561985/questions/repository/memcache"
	"github.com/d7561985/questions/repository/translate"
	"github.com/d7561985/questions/usecase/simple"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	"github.com/urfave/cli/v2"
)

const (
	ucTimeOut = time.Second * 5 // ToDo: move ot Env
	paramFile = "file"
)

var runFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    paramFile,
		Value:   "questions.json",
		Usage:   "file inside data folder",
		Aliases: []string{"f"},
		//Required: true,
	},
}

func run(ctx *cli.Context) error {
	logger := tr.NewZero(&zl.Logger)
	logger.LogFields(zerolog.InfoLevel,
		log.String("build", version.BuildContext()), log.String("version", version.Info()),
		log.String("event", "start"))

	if err := godotenv.Load(".env"); err != nil {
		logger.LogFields(zerolog.ErrorLevel, log.Error(err), log.String("event", "load .env file"))
	}

	go func() {
		if err := h2.ListenAndServe(os.Getenv("PROMETHEUS_HTTP_ADDR"), promhttp.Handler()); err != nil {
			logger.LogFields(zerolog.FatalLevel, log.String("event", "prometheus"), log.Error(err))
		}
	}()

	// WIP: for delivery and repository should be own trace instances.
	trace, closer := jaeger.NewTracer("WIP", jaeger.NewConstSampler(false), jaeger.NewNullReporter())
	defer closer.Close()

	listen, err := net.Listen("tcp", os.Getenv("DELIVERY_HTTP_ADDR"))
	if err != nil {
		logger.LogFields(zerolog.FatalLevel, log.String("addr", os.Getenv("DELIVERY_HTTP_ADDR")),
			log.Error(err))
		return err
	}

	repo := filerepo.New(logger, trace)
	if err = repo.Load(ctx.String(paramFile)); err != nil {
		logger.LogFields(zerolog.FatalLevel, log.Error(err))
		return err
	}

	trnsl, err := translate.New(os.Getenv("GOOGLE_APPLICATION_PROJECT_ID"), trace, ucTimeOut)
	if err != nil {
		logger.LogFields(zerolog.InfoLevel, log.Error(err))
	}

	uc := simple.NewService(repo, trnsl, memcache.New(), logger, ucTimeOut)
	deliver := http.New(logger, uc, trace)
	go deliver.Serve(listen)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)

	<-s

	deliver.Stop()
	if err := repo.Close(); err != nil {
		logger.LogFields(zerolog.ErrorLevel, log.Error(err))
		return err
	}

	// ToDo: cache or repo instances should also stop option
	return nil
}
