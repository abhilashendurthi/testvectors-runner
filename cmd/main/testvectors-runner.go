/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stratum/testvectors-runner/pkg/logger"
	"github.com/stratum/testvectors-runner/pkg/test"
)

var log = logger.NewLogger()

// main reads test data and utilize testing package to drive the tests. Currently two types of test data are supported.
// One is Test Vectors (see README for more details) and the other is Go function based tests (see examples under tests folder)
// To run with Test Vectors, specify Test Vector files using tvDir and tvName (optional) flag, otherwise specify test function
// names using testNames flag. A target file (tgfile) and a portmap file (pmFile) are mandatory in both cases.
func main() {
	testNames := flag.String("test-names", "", "Names of the tests to run, separated by comma")
	tvName := flag.String("tv-name", ".*", "Test Vector name specified by regular expression")
	tvDir := flag.String("tv-dir", "", "Directory of Test Vector files")
	tgFile := flag.String("target", "", "Path to the Target file")
	pmFile := flag.String("portmap", "", "Path to the portmap file")
	dpMode := flag.String("dp-mode", "direct", "Data plane mode: 'direct' or 'loopback'")
	matchType := flag.String("match-type", "exact", "Data plane match type: 'exact' or 'in'")
	logDir := flag.String("log-dir", "/tmp", "Location to store logs")
	logLevel := flag.String("log-level", "warn", "Log Level")
	templateConfig := flag.String("template-config", "", "Path to template config file")

	help := flag.Bool("help", false, "Help")
	h := flag.Bool("h", false, "Help")
	//Add -test.v to list of arguments for verbose go test output
	os.Args = append(os.Args, "-test.v")

	flag.Parse()
	flag.Usage = usage

	if *tgFile == "" || *pmFile == "" || *tvDir == "" {
		flag.Usage()
		os.Exit(3)
	}
	if *help || *h {
		flag.Usage()
		os.Exit(0)
	}

	setupLog(*logDir, *logLevel)
	testSuiteSlice := test.CreateSuite(*testNames, *tvDir, *tvName, *templateConfig)
	test.Run(*tgFile, *dpMode, *matchType, *pmFile, testSuiteSlice)
}

func setupLog(logDir string, logLevel string) {
	log.SetLogLevel(logLevel)
	log.SetLogFolder(logDir)
}

func usage() {
	usage := `Usage:
***mandatory arguments***
	[--target <filename>]               	run testvectors against the provided target proto file
	[--portmap <filename>]             		use the provided port mapping file
	[--tv-dir <directory>]              	run all the testvectors from provided directory

***optional arguments***
	[--template-config <filename>]			use the provided config file to convert templates to test vectors
	[--tv-name <regex>]                 	run all the testvectors matching provided regular expression
	[--dp-mode <mode>]                  	run the testvectors in provided mode
											default is direct; acceptable modes are <direct, loopbak>
	[--match-type <type>]               	match packets based on the provided match-type
											default is exact; acceptable modes <exact, in>
	[--log-level <level>]               	run tvrunner binary with provided log level
											default is warn; acceptable levels are <panic, fatal, error, warn, info, debug>
	[--log-dir <directory>]             	save logs to provided directory
											default is /tmp
`
	fmt.Println(usage)
}
