// Copyright © 2019 Intel Corporation
//
// SPDX-License-Identifier: GPL-3.0-only

package syscheck

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/clearlinux/clr-installer/log"
)

func getCPUFeature(feature string) error {
	cpuInfo, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Error("Unable to read /proc/cpuinfo")
		return errors.New("Unable to read /proc/cpuinfo")
	}
	if strings.Contains(string(cpuInfo), feature) {
		return nil
	}

	return errors.New("Missing CPU feature: " + feature)
}

func getEFIExist() error {
	if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
		return errors.New("Failed to find EFI firmware")
	}

	return nil
}

// RunSystemCheck checks compatibility for clear linux. (e.g. EFI firmware, CPU featureset)
func RunSystemCheck() error {
	log.Info("Running system compatibility checks.")

	//Check the following CPU features from /proc/cpuinfo
	cpuFeatures := []string{
		"lm",
		"sse4_2",
		"sse4_1",
		"pclmulqdq",
		"aes",
		"ssse3",
	}
	for _, feature := range cpuFeatures {
		fmt.Printf("Checking for required CPU feaure: %s", feature)
		err := getCPUFeature(feature)
		if err != nil {
			fmt.Printf(" [*failed*]\n")
			fmt.Println(err)
			log.ErrorError(err)

			return err
		}

		fmt.Println(" [success]")
	}

	//Check if we have EFI firmware
	fmt.Printf("Checking for required EFI firmware")
	err := getEFIExist()
	if err != nil {
		fmt.Printf(" [*failed*]\n")
		fmt.Println(err)
		log.ErrorError(err)

		return err
	}

	fmt.Println(" [success]")
	fmt.Println("Success: System is compatible")

	log.Info("Success: System is compatible")
	return nil
}