package main

import (
	"fmt"
	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

var SEPARATOR = "\n---\n"  // This is what separates manifests in a Kustomize output


func main() {

	formatter := aurora.NewAurora(true)

	if len(os.Args) != 3 {
		fmt.Println(fmt.Sprintf("usage: %s [filename] [filename]", os.Args[0]))
		os.Exit(1)
	}

	file1 := os.Args[1]
	file2 := os.Args[2]

	errors := checkFiles(file1, file2)
	failOnErr(formatter, errors...)

	contents1, err := ioutil.ReadFile(file1)
	if err != nil {
		failOnErr(formatter, err)
	}
	contents2, err := ioutil.ReadFile(file2)
	if err != nil {
		failOnErr(formatter, err)
	}

	// Compare the right file with the left file, treating the right file as the new one
	for _, string2 := range strings.Split(string(contents2), SEPARATOR) {
		m2 := Manifest{}
		err := yaml.Unmarshal([]byte(string2), &m2)
		if err != nil {
			continue
		}
		matched := false
		fmt.Println("\n--- " + m2.Kind + ": " + m2.Metadata.Name + " (" + m2.Metadata.Namespace + ")")
		for _, string1 := range strings.Split(string(contents1), SEPARATOR) {
			m1 := Manifest{}
			err := yaml.Unmarshal([]byte(string1), &m1)
			if err != nil {
				continue
			}
			
			if manifestsMatch(&m1, &m2) {
				
				y1, err := unmarshal(string1)
				if err != nil {
					break
				}
				y2, err := unmarshal(string2)
				if err != nil {
					break
				}
				matched = true
				fmt.Println(computeDiff(formatter, y1, y2))
				break

			}
		}
		if matched != true {
			fmt.Println(formatter.Bold(formatter.Cyan("> Present only in file: " + file2)).String())
		}
	}

	// Check that there were no manifests in the left file that are not in the right file
	for _, string1 := range strings.Split(string(contents1), SEPARATOR) {
		m1 := Manifest{}
		err := yaml.Unmarshal([]byte(string1), &m1)
		if err != nil {
			continue
		}
		matched := false
		for _, string2 := range strings.Split(string(contents2), SEPARATOR) {
			m2 := Manifest{}
			err := yaml.Unmarshal([]byte(string2), &m2)
			if err != nil {
				continue
			}

			if m1.Kind == m2.Kind &&
				m1.Metadata.Name == m2.Metadata.Name &&
				m1.Metadata.Namespace == m2.Metadata.Namespace {

				matched = true
				break
			}
		}
		if matched != true {
			fmt.Println("\n--- " + m1.Kind + ": " + m1.Metadata.Name + " (" + m1.Metadata.Namespace + ")")
			fmt.Println(formatter.Bold(formatter.Cyan("< Present only in file: " + file1)).String())
		}
	}
}

type Manifest struct {
	ApiVersion string `json:"apiVersion",required:"true"`
	Kind string `json:"kind",required:"true"`
	Metadata Metadata `json:"metadata",required:"true"`
	Spec map[string]interface{} `json:"spec",required:"true"`
}

type Metadata struct {
	Name string `json:"name",required:"true"`
	Namespace string `json:"namespace",required:"false"`
	Labels map[string]interface{} `json:"labels",required:"false"`
	Annotations map[string]interface{} `json:"name",required:"false"`
}

func checkFiles(filenames ...string) []error {
	var errs []error
	for _, filename := range filenames {
		_, err := os.Stat(filename)
		if err != nil {
			errs = append(errs, fmt.Errorf("cannot find file: %v", filename))
		}
	}
	return errs
}

func unmarshal(y string) (interface{}, error) {
	contents := []byte(y)
	var ret interface{}
	err := yaml.Unmarshal(contents, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func failOnErr(formatter aurora.Aurora, errs ...error) {
	if len(errs) == 0 {
		return
	}
	var errMessages []string
	for _, err := range errs {
		errMessages = append(errMessages, err.Error())
	}
	fmt.Sprintf("%v\n", formatter.Red(strings.Join(errMessages, "\n")))
	os.Exit(1)
}

func computeDiff(formatter aurora.Aurora, a interface{}, b interface{}) string {
	diffs := make([]string, 0)
	for _, s := range strings.Split(pretty.Compare(a, b), "\n") {
		switch {
		case strings.HasPrefix(s, "+"):
			diffs = append(diffs, formatter.Bold(formatter.Green(s)).String())
		case strings.HasPrefix(s, "-"):
			diffs = append(diffs, formatter.Bold(formatter.Red(s)).String())
		}
	}
	if len(diffs) == 0 {
		return formatter.Bold(formatter.Green("Manifests match")).String()
	} else {
		return strings.Join(diffs, "\n")
	}
}

func manifestsMatch(l *Manifest, r *Manifest) bool {
	if l.Kind == r.Kind &&
		l.Metadata.Name == r.Metadata.Name &&
		l.Metadata.Namespace == r.Metadata.Namespace {
		return true
	} else {
		return false
	}
}
