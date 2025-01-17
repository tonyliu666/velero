/*
Copyright 2021 the Velero contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filtering

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	. "github.com/vmware-tanzu/velero/test/e2e"
	. "github.com/vmware-tanzu/velero/test/e2e/test"
	. "github.com/vmware-tanzu/velero/test/e2e/util/k8s"
)

/*
exclude-namespaces
Exclude namespace1 from the cluster backup.
velero backup create <backup-name> --exclude-namespaces <namespace1>

Exclude two namespaces during a restore.
velero restore create <backup-name> --exclude-namespaces <namespace1>,<namespace2>
*/

type ExcludeNamespaces struct {
	FilteringCase
	nsExcluded         *[]string
	namespacesExcluded int
}

var BackupWithExcludeNamespaces func() = TestFunc(&ExcludeNamespaces{FilteringCase: testInBackup})
var RestoreWithExcludeNamespaces func() = TestFunc(&ExcludeNamespaces{FilteringCase: testInRestore})

func (e *ExcludeNamespaces) Init() error {
	e.FilteringCase.Init()
	e.namespacesExcluded = e.NamespacesTotal / 2
	e.NSBaseName = "exclude-namespaces-" + UUIDgen.String()
	if e.IsTestInBackup {
		e.BackupName = "backup-exclude-namespaces-" + UUIDgen.String()
		e.RestoreName = "restore-" + UUIDgen.String()
		e.TestMsg = &TestMSG{
			Desc:      "Backup resources with exclude namespace test",
			FailedMSG: "Failed to backup and restore with namespace include",
			Text:      fmt.Sprintf("should not backup %d namespaces of %d", e.namespacesExcluded, e.NamespacesTotal),
		}
	} else {
		e.BackupName = "backup-" + UUIDgen.String()
		e.RestoreName = "restore-exclude-namespaces-" + UUIDgen.String()
		e.TestMsg = &TestMSG{
			Desc:      "Restore resources with exclude namespace test",
			FailedMSG: "Failed to restore with namespace exclude",
			Text:      fmt.Sprintf("should not backup %d namespaces of %d", e.namespacesExcluded, e.NamespacesTotal),
		}
	}
	e.nsExcluded = &[]string{}
	for nsNum := 0; nsNum < e.NamespacesTotal; nsNum++ {
		createNSName := fmt.Sprintf("%s-%00000d", e.NSBaseName, nsNum)
		if nsNum < e.namespacesExcluded {
			*e.nsExcluded = append(*e.nsExcluded, createNSName)
		} else {
			*e.NSIncluded = append(*e.NSIncluded, createNSName)
		}
	}
	if e.IsTestInBackup {
		e.BackupArgs = []string{
			"create", "--namespace", VeleroCfg.VeleroNamespace, "backup", e.BackupName,
			"--exclude-namespaces", strings.Join(*e.nsExcluded, ","),
			"--include-namespaces", strings.Join(*e.NSIncluded, ","),
			"--default-volumes-to-fs-backup", "--wait",
		}

		e.RestoreArgs = []string{
			"create", "--namespace", VeleroCfg.VeleroNamespace, "restore", e.RestoreName,
			"--from-backup", e.BackupName, "--wait",
		}

	} else {
		*e.NSIncluded = append(*e.NSIncluded, *e.nsExcluded...)
		e.BackupArgs = []string{
			"create", "--namespace", VeleroCfg.VeleroNamespace, "backup", e.BackupName,
			"--include-namespaces", strings.Join(*e.NSIncluded, ","),
			"--default-volumes-to-fs-backup", "--wait",
		}

		e.RestoreArgs = []string{
			"create", "--namespace", VeleroCfg.VeleroNamespace, "restore", e.RestoreName,
			"--exclude-namespaces", strings.Join(*e.nsExcluded, ","),
			"--from-backup", e.BackupName, "--wait",
		}
	}
	return nil
}

func (e *ExcludeNamespaces) CreateResources() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer ctxCancel()
	for nsNum := 0; nsNum < e.NamespacesTotal; nsNum++ {
		createNSName := fmt.Sprintf("%s-%00000d", e.NSBaseName, nsNum)
		fmt.Printf("Creating namespaces ...%s\n", createNSName)
		if err := CreateNamespace(ctx, e.Client, createNSName); err != nil {
			return errors.Wrapf(err, "Failed to create namespace %s", createNSName)
		}
	}
	return nil
}

func (e *ExcludeNamespaces) Verify() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer ctxCancel()
	// Verify that we got back all of the namespaces we created
	for nsNum := 0; nsNum < e.namespacesExcluded; nsNum++ {
		excludeNSName := fmt.Sprintf("%s-%00000d", e.NSBaseName, nsNum)
		_, err := GetNamespace(ctx, e.Client, excludeNSName)
		if err == nil {
			return errors.Wrapf(err, "Resource filtering with exclude namespace but exclude namespace %s exist", excludeNSName)
		}

		if !apierrors.IsNotFound(err) {
			return errors.Wrapf(err, "Resource filtering with exclude namespace failed with checking namespace %s", excludeNSName)
		}
	}

	for nsNum := e.namespacesExcluded; nsNum < e.NamespacesTotal; nsNum++ {
		checkNSName := fmt.Sprintf("%s-%00000d", e.NSBaseName, nsNum)
		checkNS, err := GetNamespace(ctx, e.Client, checkNSName)
		if err != nil {
			return errors.Wrapf(err, "Could not retrieve test namespace %s", checkNSName)
		}
		if checkNS.Name != checkNSName {
			return errors.Errorf("Retrieved namespace for %s has name %s instead", checkNSName, checkNS.Name)
		}
	}
	return nil
}
