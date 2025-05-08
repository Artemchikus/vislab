package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	gitlabcollector "vislab/collector/gitlab"
	"vislab/config"
	"vislab/sources/gitlab"
	"vislab/storage/neo4j"
	defaultupdater "vislab/updater/default"
)

var (
	confFile string
	debug    bool
)

func init() {
	flag.StringVar(&confFile, "conf", "./config.yaml", "Path to config file")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
}

func main() {
	flag.Parse()
	// file, err := os.OpenFile("logs", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level})))

	ctx := context.Background()

	config, err := config.Get(confFile)
	if err != nil {
		slog.Error("failed to get config", "err", err)
		panic(err)
	}

	store, err := neo4j.MustStorage(ctx, config.Storage)
	if err != nil {
		slog.Error("failed to create storage", "err", err)
		panic(err)
	}

	gitlabOptions := gitlab.GetOptions(config.Collector.GitLab.Client)

	gitlabClient, err := gitlab.NewClient(config.Collector.GitLab.Client.Token, config.Collector.GitLab.Client.BaseURL, gitlabOptions...)
	if err != nil {
		slog.Error("failed to create gitlab client", "err", err)
		panic(err)
	}

	collectorOptions, err := gitlabcollector.GetOptions(config.Collector, config.Sources)
	if err != nil {
		slog.Error("failed to get collector options", "err", err)
		panic(err)
	}

	collector, err := gitlabcollector.New(gitlabClient, store, collectorOptions...)
	if err != nil {
		slog.Error("failed to create collector", "err", err)
		panic(err)
	}

	updater, err := defaultupdater.New(collector, config.Updater.Port)
	if err != nil {
		slog.Error("failed to create updater", "err", err)
		panic(err)
	}

	if err := collector.Collect(ctx); err != nil {
		slog.Error("failed to collect data", "err", err)
		panic(err)
	}

	slog.Info("starting updater")
	if err := updater.Start(ctx); err != nil {
		slog.Error("failed to run updater", "err", err)
		panic(err)
	}

	if err := store.Disconnect(ctx); err != nil {
		slog.Error("failed to disconnect from storage", "err", err)
		panic(err)
	}
}
