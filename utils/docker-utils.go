package utils

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// func DockerRun(tempDatabaseDir string, tempDBName string, uuid string, dockerName string) {
// 	dockerContainerID := `docker run --name ` + dockerName + ` \
// 		-v ` + tempDatabaseDir + `:/var/lib/postgresql/data:Z \
// 		-e POSTGRES_PASSWORD=gitops \
// 		-e POSTGRES_DB=` + tempDBName + ` \
// 		-p 6432:6432 \
// 		--network gitops-net-` + uuid + ` \
// 		-d \
// 		postgres:13 \
// 		-c log_statement='all' \
// 		-c log_min_duration_statement=0`
// 	fmt.Println(dockerContainerID)

// 	fmt.Println(dockerContainerID)

// }

func Docker_apply() {
	dockerName := "managed-gitops-postgres-test"
	dockerNetworkcmd := "docker network create %s"
	uuid := "sam"
	tempDBName := "db-" + uuid
	s := fmt.Sprintf(dockerNetworkcmd, "gitops-net-"+uuid)

	// To print which command is running
	fmt.Println("\nRunning: ", s)

	// To get the output of the command
	out, err := exec.Command("docker", "network", "create", "gitops-net-"+uuid).Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out), "`docker create network` command Run Successful!")

	tempDatabaseDircmd := "mktemp -d -t postgres-XXXXXXXXXX"
	s = fmt.Sprintf(tempDatabaseDircmd)

	// To print which command is running
	fmt.Println("\nRunning: ", s)

	// To actually run the command (runs in background)
	tempDatabaseDir, err_run := exec.Command("mktemp", "-d", "-t", "postgres-XXXXXXXXXX").Output()
	if err_run != nil {
		log.Fatal(err_run)
	}
	fmt.Println(string(tempDatabaseDir), "`mktemp dir` command Run Successful!")

	// running a docker container
	dockerContainerIDcmd := `docker run --name ` + dockerName + ` \
	-v ` + string(tempDatabaseDir) + `:/var/lib/postgresql/data:Z \
	-e POSTGRES_PASSWORD=gitops \
	-e POSTGRES_DB=` + tempDBName + ` \
	-p 6432:6432 \
	--network gitops-net-` + uuid + ` \
	-d \
	postgres:13 \
	-c log_statement='all' \
	-c log_min_duration_statement=0`

	fmt.Println("\nRunning:", dockerContainerIDcmd)

	err = wait.Poll(5*time.Second, 2*time.Minute, func() (bool, error) {
		dockerContainerID, errDockerRun := exec.Command("docker", "run", "--name", dockerName,
			"-v", string(tempDatabaseDir)+":/var/lib/postgresql/data:Z",
			"-e", "POSTGRES_PASSWORD=gitops",
			"-e", "POSTGRES_DB="+tempDBName,
			"-p", "6432:6432",
			"--network", "gitops-net-"+uuid,
			"-d",
			"postgres:13",
			"-c", "log_statement=all",
			"-c", "log_min_duration_statement=0").Output()

		if errDockerRun != nil {
			log.Fatal(errDockerRun)
		}
		fmt.Println(string(dockerContainerID), "`docker run` command Run Successful!")
		return true, nil
	})

	if err != nil {
		log.Fatal("error in executing docker run command")
	}

}
