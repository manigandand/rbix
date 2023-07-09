package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/manigandand/adk/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type k8s struct {
	client *kubernetes.Clientset
}

// NewK8sOrchestrator - create new k8s orchestrator
func NewK8sOrchestrator() (Orchestrator, *errors.AppError) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.InternalServer("could not create k8s config: " + err.Error())
	}
	// creates the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.InternalServer("could not create k8s client: " + err.Error())
	}

	if _, err := client.RbacV1().ClusterRoleBindings().Get(context.Background(), "ns:platform-u:default-r:cluster-admin", metav1.GetOptions{}); err != nil {
		return nil, errors.InternalServer("could not get cluster role binding(ns:platform-u:default-r:cluster-admin): " + err.Error())
	}

	return &k8s{
		client: client,
	}, nil
}

// StartRBIInstance - start new sqrx-rbi instance
func (k *k8s) StartRBIInstance(ctx context.Context, containerUniqeID string,
) (*ContainerInfo, *errors.AppError) {
	containerInfo := &ContainerInfo{
		ContainerID: containerUniqeID,
		CreatedAt:   time.Now(),
	}

	// create namespace
	newns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: containerUniqeID,
		},
	}
	nsres, err := k.client.CoreV1().Namespaces().Create(ctx, newns, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.InternalServer("could not create namespace: " + err.Error())
	}
	log.Printf("namespace created: %+v\n", nsres)

	// create deployment
	depRes, err := k.client.AppsV1().Deployments(containerUniqeID).Create(
		ctx, newSqrxRBIDeployment(containerUniqeID), metav1.CreateOptions{},
	)
	if err != nil {
		return nil, errors.InternalServer("could not create deployment: " + err.Error())
	}

	// create service
	svcRes, err := k.client.CoreV1().Services(containerUniqeID).Create(
		ctx, newSqrxRBIService(containerUniqeID), metav1.CreateOptions{},
	)
	if err != nil {
		return nil, errors.InternalServer("could not expose container: " + err.Error())
	}
	log.Println(containerUniqeID, " -> deployment exposed: status: ", svcRes.Status)
	containerInfo.StartedAt = time.Now()

	containerInfo.K8sContainer = depRes
	containerInfo.ValidTill = time.Now().Add(10 * time.Minute)
	containerInfo.Session = fmt.Sprintf("%s/%s/ws", SqrxWSLoadbalncerHost, containerUniqeID)
	containerInfo.TerminationToken = uuid.New().String()

	// save to db
	db.SaveContainerInfo(containerUniqeID, containerInfo)
	db.SaveTerminationTokenInfo(containerInfo.TerminationToken, containerUniqeID)

	// spinup background process to delete the container after 10 minutes
	go func() {
		<-time.After(10 * time.Minute)
		if err := k.DestroyRBIInstance(context.Background(), containerInfo.TerminationToken); err.NotNil() {
			log.Println("couldn't able to delete container: ", err.Error())
		}
	}()

	return containerInfo, nil
}

// DestroyRBIInstance - destroy sqrx-rbi instance
func (k *k8s) DestroyRBIInstance(ctx context.Context, terminationToken string) *errors.AppError {
	cInfo, err := db.GetContainerInfoByTermToken(terminationToken)
	if err.NotNil() {
		return err
	}

	if err := k.client.CoreV1().Namespaces().Delete(ctx, cInfo.ContainerID, metav1.DeleteOptions{}); err != nil {
		return errors.InternalServer("could not delete container: " + err.Error())
	}

	db.DeleteContainer(cInfo.ContainerID)

	return nil
}

// Stop - stop orchestrator
func (k *k8s) Stop() {
}

func newSqrxRBIDeployment(containerUniqeID string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sqrx-rbi",
			Namespace: containerUniqeID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"run": "sqrx-rbi"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"run": "sqrx-rbi"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "sqrx-rbi",
							Image: SqrxRbiImage,
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 8888,
								},
							},
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("2500m"),
									v1.ResourceMemory: resource.MustParse("712Mi"),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("2000m"),
									v1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "ENV",
									Value: EnvLocalK8s,
								},
								{
									Name:  "PORT",
									Value: "8888",
								},
								{
									Name:  "CONTAINER_ID",
									Value: containerUniqeID,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func newSqrxRBIService(containerUniqeID string) *v1.Service {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sqrx-rbi",
			Namespace: containerUniqeID,
			Labels:    map[string]string{"run": "sqrx-rbi"},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"run": "sqrx-rbi"},
			Type:     v1.ServiceTypeClusterIP,
			Ports: []v1.ServicePort{
				{
					Name:     "http",
					Protocol: v1.ProtocolTCP,
					Port:     8888,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8888,
					},
				},
			},
		},
	}
	return svc
}

func int32Ptr(i int32) *int32 { return &i }
