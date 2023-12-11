package main

import (
	"context"
	"github.com/jessevdk/go-flags"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	server "micro-manager-tasks/m/v2/app/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	Listen         string        `short:"l" long:"listen" env:"LISTEN" default:":8080" description:"listen address"`
	Secret         string        `short:"s" long:"secret" env:"EVENT_SECRET_KEY" default:"123"`
	PinSize        int           `long:"pinszie" env:"PIN_SIZE" default:"5" description:"pin size"`
	MaxExpire      time.Duration `long:"expire" env:"MAX_EXPIRE" default:"24h" description:"max lifetime"`
	MaxPinAttempts int           `long:"pinattempts" env:"PIN_ATTEMPTS" default:"3" description:"max attempts to enter pin"`
	WebRoot        string        `long:"web" env:"WEB" default:"/" description:"web ui location"`
	MongoUrl       string        `long:"mongo" env:"MONGO_URL" default:"mongodb://localhost:27017" description:"mongo url"`
	Database       string        `long:"db" env:"DATABASE" default:"micro-tasks" description:"database name"`
	MongoUser      string        `long:"dbuser" env:"MONGODB_ROOT_USER" default:"root" description:"database user"`
	MongoPass      string        `long:"dbpass" env:"MONGODB_ROOT_PASSWORD" default:"F4u3b8BYQKad" description:"database password"`
}

var revision string

func main() {
	log.Printf("Micro Manager tasks %s\n", revision)

	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {

		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if x := recover(); x != nil {
			log.Printf("[WARN] run time panic:\n%v", x)
			panic(x)
		}

		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()
	client, err := getMongoConnection(opts, ctx)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.Server{
		Listen:         opts.Listen,
		PinSize:        opts.PinSize,
		MaxExpire:      opts.MaxExpire,
		MaxPinAttempts: opts.MaxPinAttempts,
		WebRoot:        opts.WebRoot,
		Secret:         opts.Secret,
		Version:        revision,
		Client:         client,
		Database:       opts.Database,
	}
	if err := srv.Run(ctx); err != nil {
		log.Printf("[ERROR] failed, %+v", err)
	}
}

func getMongoConnection(opts Options, ctx context.Context) (*mongo.Client, error) {
	cred := options.Credential{
		AuthSource: "admin",
		Username:   opts.MongoUser,
		Password:   opts.MongoPass,
	}
	clientOptions := options.Client().ApplyURI(opts.MongoUrl).SetAuth(cred)
	return mongo.Connect(ctx, clientOptions)
}
