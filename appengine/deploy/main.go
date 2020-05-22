package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Copy(dst, src string) error {
	fsrc, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file - %w", err)
	}
	defer fsrc.Close()

	fdst, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("could not open destination file - %w", err)
	}
	defer fdst.Close()

	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}

	if err := fdst.Sync(); err != nil {
		return err
	}

	return fdst.Close()
}

func main() {
	path := flag.String("path", "", "Top level directory where all the repositories are stored (eg, appengine/test/test-deploy-dir)")
	entry := flag.String("entry", "", "Directory project containing your app (eg, github.com/ccontavalli/myapp)")
	config := flag.String("config", "", "Path to the app.yaml file to use for the deploy")
	gomod := flag.String("gomod", "", "Path to a go.mod file to prepare for deploy")
	gosum := flag.String("gosum", "", "Path to a go.sum file to prepare for deploy")
	gcloud := flag.String("gcloud", "/usr/bin/gcloud", "Path to the gcloud binary to use")
	quiet := flag.Bool("quiet", false, "Be more quiet")
	flag.Parse()

	if *entry == "" || *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	project := filepath.Join(*path, "src", *entry)
	if s, err := os.Stat(project); err != nil || !s.Mode().IsDir() {
		log.Fatalf("Couldn't enter directory supplied with --entry=%s: %s directory caused %s", *entry, project, err)
	}

	if *gomod != "" {
		dest := filepath.Join(project, "go.mod")
		if err := Copy(dest, *gomod); err != nil {
			log.Fatalf("Couldn't copy %s - supplied with -gomod - to %s: %s", *gomod, dest, err)
		}
	}
	if *gosum != "" {
		dest := filepath.Join(project, "go.sum")
		if err := Copy(dest, *gosum); err != nil {
			log.Fatalf("Couldn't copy %s - supplied with -gosum - to %s: %s", *gosum, dest, err)
		}
	}

	if *config != "" {
		dest := filepath.Join(project, "app.yaml")
		if err := Copy(dest, *config); err != nil {
			log.Fatalf("Couldn't copy %s - supplied with -config - to %s: %s", *config, dest, err)
		}
	}

	fmt.Fprintf(os.Stderr, "Deploying '%s' to cloud\n", project)

	if !*quiet {
		cmd := exec.Command("/usr/bin/find", ".")
		cmd.Dir = *path
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("find failed: %s", err)
		}
	}

	if *gcloud != "" {
		cmd := exec.Command(*gcloud, "app", "deploy")
		cmd.Dir = filepath.Join(project)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("%s failed: %s", *gcloud, err)
		}
	}
}
