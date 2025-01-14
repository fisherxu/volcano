/*
Copyright 2019 The Volcano Authors.

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

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vkv1 "volcano.sh/volcano/pkg/apis/batch/v1alpha1"
)

var _ = Describe("TensorFlow E2E Test", func() {
	It("Will Start in pending state and goes through other phases to get complete phase", func() {
		context := initTestContext()
		defer cleanupTestContext(context)

		jobName := "tensorflow-dist-mnist"

		job := &vkv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name: jobName,
			},
			Spec: vkv1.JobSpec{
				MinAvailable:  int32(3),
				SchedulerName: schedulerName,
				Plugins: map[string][]string{
					"svc": {},
					"env": {},
				},
				Policies: []vkv1.LifecyclePolicy{
					{
						Event:  vkv1.PodEvictedEvent,
						Action: vkv1.RestartJobAction,
					},
				},
				Tasks: []vkv1.TaskSpec{
					{
						Replicas: int32(1),
						Name:     "ps",
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								RestartPolicy: v1.RestartPolicyNever,
								Containers: []v1.Container{
									{
										Command: []string{
											"sh",
											"-c",
											"PS_HOST=`cat /etc/volcano/ps.host | sed 's/$/&:2222/g' | sed 's/^/\"/;s/$/\"/' | tr \"\n\" \",\"`; WORKER_HOST=`cat /etc/volcano/worker.host | sed 's/$/&:2222/g' | sed 's/^/\"/;s/$/\"/' | tr \"\n\" \",\"`; export TF_CONFIG={\\\"cluster\\\":{\\\"ps\\\":[${PS_HOST}],\\\"worker\\\":[${WORKER_HOST}]},\\\"task\\\":{\\\"type\\\":\\\"ps\\\",\\\"index\\\":${VK_TASK_INDEX}},\\\"environment\\\":\\\"cloud\\\"}; echo ${TF_CONFIG}; python /var/tf_dist_mnist/dist_mnist.py --train_steps 1000",
										},
										Image: defaultTFImage,
										Name:  "tensorflow",
										Ports: []v1.ContainerPort{
											{
												Name:          "tfjob-port",
												ContainerPort: int32(2222),
											},
										},
									},
								},
							},
						},
					},
					{
						Replicas: int32(2),
						Name:     "worker",
						Policies: []vkv1.LifecyclePolicy{
							{
								Event:  vkv1.TaskCompletedEvent,
								Action: vkv1.CompleteJobAction,
							},
						},
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								RestartPolicy: v1.RestartPolicyNever,
								Containers: []v1.Container{
									{
										Command: []string{
											"sh",
											"-c",
											"PS_HOST=`cat /etc/volcano/ps.host | sed 's/$/&:2222/g' | sed 's/^/\"/;s/$/\"/' | tr \"\n\" \",\"`; WORKER_HOST=`cat /etc/volcano/worker.host | sed 's/$/&:2222/g' | sed 's/^/\"/;s/$/\"/' | tr \"\n\" \",\"`; export TF_CONFIG={\\\"cluster\\\":{\\\"ps\\\":[${PS_HOST}],\\\"worker\\\":[${WORKER_HOST}]},\\\"task\\\":{\\\"type\\\":\\\"worker\\\",\\\"index\\\":${VK_TASK_INDEX}},\\\"environment\\\":\\\"cloud\\\"}; echo ${TF_CONFIG}; python /var/tf_dist_mnist/dist_mnist.py --train_steps 1000",
										},
										Image: defaultTFImage,
										Name:  "tensorflow",
										Ports: []v1.ContainerPort{
											{
												Name:          "tfjob-port",
												ContainerPort: int32(2222),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		created, err := context.vcclient.BatchV1alpha1().Jobs("test").Create(job)
		Expect(err).NotTo(HaveOccurred())

		err = waitJobStates(context, created, []vkv1.JobPhase{vkv1.Pending, vkv1.Inqueue, vkv1.Running, vkv1.Completed}, twoMinute)
		Expect(err).NotTo(HaveOccurred())
	})

})
