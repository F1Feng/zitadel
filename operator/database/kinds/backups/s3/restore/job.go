package restore

import (
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/pkg/databases/db"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getJob(
	namespace string,
	nameLabels *labels.Name,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	accessKeyIDName string,
	accessKeyIDKey string,
	secretAccessKeyName string,
	secretAccessKeyKey string,
	sessionTokenName string,
	sessionTokenKey string,
	image string,
	command string,
	env *corev1.EnvVar,

) *batchv1.Job {

	var envs []corev1.EnvVar
	if env != nil {
		envs = []corev1.EnvVar{*env}
	}

	return &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: namespace,
			Labels:    labels.MustK8sMap(nameLabels),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeSelector:  nodeselector,
					Tolerations:   tolerations,
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{{
						Name:  nameLabels.Name(),
						Image: image,
						Command: []string{
							"/bin/bash",
							"-c",
							command,
						},
						Env: envs,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      internalSecretName,
							MountPath: certPath,
						}, {
							Name:      accessKeyIDKey,
							SubPath:   accessKeyIDKey,
							MountPath: accessKeyIDPath,
						}, {
							Name:      secretAccessKeyKey,
							SubPath:   secretAccessKeyKey,
							MountPath: secretAccessKeyPath,
						}, {
							Name:      sessionTokenKey,
							SubPath:   sessionTokenKey,
							MountPath: sessionTokenPath,
						}},
						ImagePullPolicy: corev1.PullIfNotPresent,
					}},
					Volumes: []corev1.Volume{{
						Name: internalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  db.CertsSecret,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: accessKeyIDKey,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: accessKeyIDName,
							},
						},
					}, {
						Name: secretAccessKeyKey,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: secretAccessKeyName,
							},
						},
					}, {
						Name: sessionTokenKey,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: sessionTokenName,
							},
						},
					}},
				},
			},
		},
	}
}