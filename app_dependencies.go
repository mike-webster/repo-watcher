package main

import (
	dispatchers "github.com/mike-webster/repo-watcher/dispatchers"
	"github.com/sirupsen/logrus"
)

type AppDependencies struct {
	logger      *logrus.Logger
	dispatchers dispatchers.Dispatchers
}
