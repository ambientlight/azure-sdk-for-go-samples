// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package maps

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/services/maps/mgmt/2021-02-01/maps"
)

var (
	mapsAccount    *maps.Account
	creatorAccount *maps.Creator
)

func addLocalEnvAndParse() error {
	// parse env at top-level (also controls dotenv load)
	err := config.ParseEnvironment()
	if err != nil {
		return fmt.Errorf("failed to add top-level env: %+v", err)
	}

	return nil
}

func addLocalFlagsAndParse() error {
	// add top-level flags
	err := config.AddFlags()
	if err != nil {
		return fmt.Errorf("failed to add top-level flags: %+v", err)
	}

	// parse all flags
	flag.Parse()
	return nil
}

var usesADAuth = flag.Bool("ad-auth", false, "uses Azure Maps AD authentication instead of Shared Key if set")

func setup() error {
	var err error
	err = addLocalEnvAndParse()
	if err != nil {
		return err
	}

	err = addLocalFlagsAndParse()
	if err != nil {
		return err
	}

	newMapsAccount, newCreatorAccount := CreateResourceGroupWithMapAndCreatorAccount()
	mapsAccount = &newMapsAccount
	creatorAccount = &newCreatorAccount
	return nil
}

func teardown() error {
	ctx := context.Background()
	mapsAccount = nil
	creatorAccount = nil
	resources.Cleanup(ctx)
	return nil
}

// TestMain sets up the environment and initiates tests.
func TestMain(m *testing.M) {
	var err error
	var code int

	err = setup()
	if err != nil {
		log.Fatalf("could not set up environment: %+v", err)
	}

	defer func() {
		err = teardown()
		if err != nil {
			log.Fatalf(
				"could not tear down environment: %v\n; original exit code: %v\n",
				err, code)
		}
		os.Exit(code)
	}()

	code = m.Run()
}
