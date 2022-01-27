package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func Docker_apply() {

	// cleanup := func(dockerContainerID string, networkBridge string) {
	// 	// kill the docker container
	// 	// rm the docker container
	// 	// => docker rm -f container id
	// 	//  remove the network bridge
	// }

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

	fmt.Println(string(out), "`docker network create` command Run Successful!")

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
	-p 6432:5432 \
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
			"-p", "6432:5432",
			"--network", "gitops-net-"+uuid,
			"-d",
			"postgres:13",
			"-c", "log_statement=all",
			"-c", "log_min_duration_statement=0").Output()

		if errDockerRun != nil {
			log.Fatal(errDockerRun)
		}
		if dockerContainerID == nil {
			return false, errDockerRun
		}
		// check for container status
		status, _ := exec.Command("docker", "container", "inspect", "-f", "{{.State.Status}}", string(dockerContainerID)).Output()
		if string(status) == "running" {
			return true, nil
		}

		fmt.Println("Docker Container ID: " + string(dockerContainerID))
		fmt.Println("`docker run` command Run Successful!")

		return true, nil
	})
	if err != nil {
		log.Fatal("error in executing docker run command")
	}

	dbcmd := "PGPASSWORD=gitops psql -h localhost -d %s -U postgres -p 6432 -c 'select 1'"
	s = fmt.Sprintf(dbcmd, tempDBName)

	fmt.Println("\nRunning: ", s)
	// To get the output of the command
	err = wait.Poll(5*time.Second, 2*time.Minute, func() (bool, error) {
		psqlcmd := exec.Command("psql", "-h", "localhost", "-d", tempDBName, "-U", "postgres", "-p", "6432", "-c", "select 1")
		psqlcmd.Env = os.Environ()
		psqlcmd.Env = append(psqlcmd.Env, "PGPASSWORD=gitops")
		var outb, errb bytes.Buffer
		psqlcmd.Stdout = &outb
		psqlcmd.Stderr = &errb

		psqlErr := psqlcmd.Run()

		fmt.Println("out:", outb.String(), "err:", errb.String())
		if psqlErr != nil {
			return false, psqlErr
		}
		fmt.Println("database is ready to use")
		return true, nil
	})

	// creating a new database
	newDBName := "postgres"
	dbcmd = "PGPASSWORD=gitops psql -h localhost -d %s -U postgres -p 6432"
	s = fmt.Sprintf(dbcmd, newDBName)
	fmt.Println("\nRunning: ", s)

	psqlcmd := exec.Command("psql", "-h", "localhost", "-d", newDBName, "-U", "postgres", "-p", "6432")
	psqlcmd.Env = os.Environ()
	psqlcmd.Env = append(psqlcmd.Env, "PGPASSWORD=gitops")
	var errb bytes.Buffer
	psqlcmd.Stderr = &errb

	psqlErr := psqlcmd.Run()

	fmt.Println("err:", errb.String())
	if psqlErr != nil {
		log.Fatal(psqlErr)
	}

	fmt.Println(newDBName, "database is created and ready to use")

	// Following command is used to populate the database tables from the db-schema.sql (defined in the monorepo)
	dbcmd = "PGPASSWORD=gitops psql -h localhost -d %s -U postgres -p 6432 -q -f db-schema.sql"
	s = fmt.Sprintf(dbcmd, newDBName)
	fmt.Println("\nRunning: ", s)
	psqlcmd = exec.Command("psql", "-h", "localhost", "-d", newDBName, "-U", "postgres", "-p", "6432", "-q", "-f", "db-schema.sql")
	psqlcmd.Env = os.Environ()
	psqlcmd.Env = append(psqlcmd.Env, "PGPASSWORD=gitops")
	psqlErr = psqlcmd.Run()

	fmt.Println("err:", errb.String())
	if psqlErr != nil {
		log.Fatal(psqlErr)
	}
	fmt.Println("schema executed in the postgres")

}
