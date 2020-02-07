package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var TemplatesMap = make(map[string]JTemplate)

var yamlTempl = `apiVersion: batch/v1
kind: Job
metadata:
  annotations:
    author:
  labels:
    job-name:
  name:
  namespace:
spec:
  backoffLimit: 0
  completions: 1
  parallelism: 1
  template:
    spec:
      containers:
        - command:
            - "/bin/sh"
            - "-c"
          image:
          env:
            - name: SOURCE_NS
              value: example
          imagePullPolicy: Always
          name:
          terminationMessagePath: "/dev/termination-log"
          terminationMessagePolicy: File
          volumeMounts:
      dnsPolicy: ClusterFirst
      imagePullSecrets:
        - name: imagepullsecret
      restartPolicy: Never
      schedulerName: default-scheduler
      securityContext:
      terminationGracePeriodSeconds: 30
      volumes:
      serviceAccount:
      serviceAccountName:`

type JTemplate struct {
	APIVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Annotations struct {
			Author string `yaml:"author" json:"author"`
		} `yaml:"annotations" json:"annotations"`
		Labels struct {
			JobName string `yaml:"job-name" json:"job-name"`
		} `yaml:"labels" json:"labels"`
		Name      string `yaml:"name" json:"name"`
		Namespace string `yaml:"namespace" json:"namespace"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		BackoffLimit int `yaml:"backoffLimit" json:"backoffLimit"`
		Completions  int `yaml:"completions" json:"completions"`
		Parallelism  int `yaml:"parallelism" json:"parallelism"`
		Template     struct {
			Spec struct {
				Containers []struct {
					Command []string `yaml:"command" json:"command"`
					Image   string   `yaml:"image" json:"image"`
					Env     []struct {
						Name  string `yaml:"name" json:"name"`
						Value string `yaml:"value" json:"value"`
					} `yaml:"env" json:"env"`
					ImagePullPolicy          string      `yaml:"imagePullPolicy" json:"imagePullPolicy"`
					Name                     string      `yaml:"name" json:"name"`
					TerminationMessagePath   string      `yaml:"terminationMessagePath" json:"terminationMessagePath"`
					TerminationMessagePolicy string      `yaml:"terminationMessagePolicy" json:"terminationMessagePolicy"`
					VolumeMounts             []struct {
						MountPath string `yaml:"mountPath" json:"mountPath"`
						Name      string `yaml:"name" json:"name"`
						SubPath   string `yaml:"subPath" json:"subPath"`
					} `yaml:"volumeMounts" json:"volumeMounts"`
				} `yaml:"containers" json:"containers"`
				DNSPolicy        string `yaml:"dnsPolicy" json:"dnsPolicy"`
				ImagePullSecrets []struct {
					Name string `yaml:"name" json:"name"`
				} `yaml:"imagePullSecrets" json:"imagePullSecrets"`
				RestartPolicy                 string      `yaml:"restartPolicy" json:"restartPolicy"`
				SchedulerName                 string      `yaml:"schedulerName" json:"schedulerName"`
				SecurityContext               interface{} `yaml:"securityContext" json:"securityContext"`
				TerminationGracePeriodSeconds int         `yaml:"terminationGracePeriodSeconds" json:"terminationGracePeriodSeconds"`
				Volumes                       []struct {
					Name   string `yaml:"name" json:"name"`
					Secret struct {
						DefaultMode int    `yaml:"defaultMode" json:"defaultMode"`
						SecretName  string `yaml:"secretName" json:"secretName"`
					} `yaml:"secret" json:"secret"`
				} `yaml:"volumes" json:"volumes"`
				ServiceAccount                string      `yaml:"serviceAccount" json:"serviceAccount"`
				ServiceAccountName            string      `yaml:"serviceAccountName" json:"serviceAccountName"`
			} `yaml:"spec" json:"spec"`
		} `yaml:"template" json:"template"`
	} `yaml:"spec" json:"spec"`
}

func MakeFillerJob() JTemplate{
	jobName := "filler-test"

	var fjob = TemplatesMap["job"]

	fjob.Metadata.Labels.JobName = jobName

	fjob.Metadata.Name = jobName
	fjob.Metadata.Namespace = "filler-namespace"

	var SecretsScript = `Filler
						Job
						Activity
						-------`

	fjob.Spec.Template.Spec.Containers[0].Command = append(fjob.Spec.Template.Spec.Containers[0].Command, SecretsScript)

	fjob.Spec.Template.Spec.Containers[0].Env = append(fjob.Spec.Template.Spec.Containers[0].Env,  struct{
		Name  string `yaml:"name" json:"name"`
		Value string `yaml:"value" json:"value"`
	}{
		Name: "NEW_NAMESPACE",
		Value: jobName,
	})
	fjob.Spec.Template.Spec.Containers[0].Env = append(fjob.Spec.Template.Spec.Containers[0].Env,  struct{
		Name  string `yaml:"name" json:"name"`
		Value string `yaml:"value" json:"value"`
	}{
		Name: "K8S_TOKEN",
		Value: "token",
	})
	fjob.Spec.Template.Spec.Containers[0].Name = jobName

	return fjob
}

func MakeDeployJob() JTemplate{
	jobName := "deploy-test"

	var DeployScript = `Deploy
						Job
						Activity
						--------`

	var djob = TemplatesMap["job"]

	djob.Metadata.Labels.JobName = jobName
	djob.Metadata.Name = jobName
	djob.Metadata.Namespace = "deploy-namespace"
	djob.Spec.Template.Spec.Containers[0].Command = append(djob.Spec.Template.Spec.Containers[0].Command, DeployScript)
	djob.Spec.Template.Spec.Containers[0].Env = append(djob.Spec.Template.Spec.Containers[0].Env,  struct{
		Name  string `yaml:"name" json:"name"`
		Value string `yaml:"value" json:"value"`
	}{
		Name: "NEW_NAMESPACE",
		Value: jobName,
	})

	djob.Spec.Template.Spec.Containers[0].VolumeMounts = append(djob.Spec.Template.Spec.Containers[0].VolumeMounts, struct {
		MountPath string `yaml:"mountPath" json:"mountPath"`
		Name      string `yaml:"name" json:"name"`
		SubPath   string `yaml:"subPath" json:"subPath"`
	}{
		MountPath: "/root/.ssh/id_rsa",
		Name: "ssh-key-volume",
		SubPath: "id_rsa",
	})
	djob.Spec.Template.Spec.Volumes = append(djob.Spec.Template.Spec.Volumes, struct {
		Name   string `yaml:"name" json:"name"`
		Secret struct {
			DefaultMode int    `yaml:"defaultMode" json:"defaultMode"`
			SecretName  string `yaml:"secretName" json:"secretName"`
		} `yaml:"secret" json:"secret"`
	}{
		Name: "ssh-key-volume",
		Secret: struct {
			DefaultMode int    `yaml:"defaultMode" json:"defaultMode"`
			SecretName  string `yaml:"secretName" json:"secretName"`
		}{
			DefaultMode: 384,
			SecretName: "gitlab-ssh-key",
		},
	})
	djob.Spec.Template.Spec.ServiceAccount = "jeeves"
	djob.Spec.Template.Spec.ServiceAccountName = "jeeves"
	djob.Spec.Template.Spec.Containers[0].Name = jobName
	return djob
}

func GetTemplates(){
	var jtmpl JTemplate
	err := yaml.Unmarshal([]byte(yamlTempl), &jtmpl)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Preload failed. Unable to parse yaml template")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	TemplatesMap["job"] = jtmpl
}

func init()  {
	GetTemplates()
}

func main() {
	a := MakeFillerJob()
	b := MakeDeployJob()
	fmt.Println(a)
	fmt.Println(b)
}
