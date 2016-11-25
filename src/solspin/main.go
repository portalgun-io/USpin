//
// Copyright © 2016 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package main provides the solspin command line utility for image generation
package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"libspin"
	"libspin/image"
	"libspin/pkg"
	"os"
)

var log *logrus.Logger

// Set up the main logger formatting used in SolSpin
func init() {
	form := &logrus.TextFormatter{}
	form.FullTimestamp = true
	form.TimestampFormat = "15:04:05"
	log = logrus.New()
	log.Out = os.Stderr
	log.Formatter = form
}

// SolSpin is the main solspin binary lifetime tracking object
type SolSpin struct {
	logImage   *logrus.Entry
	logPackage *logrus.Entry

	builder  image.Builder
	packager pkg.Manager
	spec     *libspin.ImageSpec
}

// NewSolSpin will return a new SolSpin instance which stores global
// state for the duration of an image spin process.
func NewSolSpin(path string) (*SolSpin, error) {
	ret := &SolSpin{}
	var err error

	// Attempt to get the image spec first
	if ret.spec, err = libspin.NewImageSpec(path); err != nil {
		return nil, err
	}

	// Get a builder
	buildType := ret.spec.Config.Image.Type
	if ret.builder, err = image.NewBuilder(buildType); err != nil {
		return nil, err
	}

	// Get our image log
	ret.logImage = log.WithFields(logrus.Fields{"imageType": buildType})

	// TODO: Stop hardcoding this!
	pkgType := "eopkg"

	// Get our package manager
	if ret.packager, err = pkg.NewManager(pkgType); err != nil {
		return nil, err
	}

	// Get packager log
	ret.logPackage = log.WithFields(logrus.Fields{"packageManager": pkgType})

	return ret, nil
}

func printUsage(exitCode int) {
	var fd *os.File
	if exitCode == 0 {
		fd = os.Stdout
	} else {
		fd = os.Stderr
	}

	fmt.Fprintf(fd, "%s [image.spin]\n", os.Args[0])
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) < 2 {
		printUsage(1)
	}

	spin, err := NewSolSpin(os.Args[1])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if err := spin.Build(); err != nil {
		os.Exit(1)
	}
}
