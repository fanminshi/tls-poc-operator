package stub

import (
	"context"

	"github.com/fanminshi/tls-poc-operator/pkg/apis/security/v1alpha1"

	"github.com/fanminshi/operator-sdk/pkg/util/tlsutil"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHandler(ca tlsutil.CA) sdk.Handler {
	return &Handler{ca: ca}
}

type Handler struct {
	ca tlsutil.CA
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch cr := event.Object.(type) {
	case *v1alpha1.Security:
		// Create a simple-server-service that is going to be sit in front the simple server pod.
		svc := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "simple-server-service",
				Namespace: "default",
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{"app": "simple-server"},
				Ports: []corev1.ServicePort{corev1.ServicePort{
					Name:     "https",
					Protocol: corev1.ProtocolTCP,
					Port:     443,
				}},
			},
		}
		err := sdk.Create(svc)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
		// Generate TLS assets based on the svc and CertConfig.
		se, err := h.ca.GenerateCert(cr, svc, &tlsutil.CertConfig{CertName: "tls"})
		if err != nil {
			return err
		}
		err = sdk.Create(se)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}

		cm, casecret, err := h.ca.CACert(cr)
		if err != nil {
			return err
		}
		err = sdk.Create(cm)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
		err = sdk.Create(casecret)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}

		// Deploy the simple-server using the "quay.io/fanminshi/simple-server:latest" image
		// and TLS assets generated from the above.
		replicas := int32(1)
		de := &v1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "simple-server",
				Namespace: "default",
				Labels:    map[string]string{"app": "simple-server"},
			},
			Spec: v1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "simple-server"}},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "simple-server"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							corev1.Container{
								Name:            "simple-server",
								Image:           "quay.io/fanminshi/simple-server:latest",
								Ports:           []corev1.ContainerPort{corev1.ContainerPort{ContainerPort: 443}},
								ImagePullPolicy: corev1.PullAlways,
								VolumeMounts: []corev1.VolumeMount{corev1.VolumeMount{
									Name:      "tls",
									MountPath: "/etc/tls",
									ReadOnly:  true,
								}},
								Env: []corev1.EnvVar{corev1.EnvVar{
									Name:  "KEY",
									Value: "/etc/tls/server.key",
								}, corev1.EnvVar{
									Name:  "CERT",
									Value: "/etc/tls/server.crt",
								}},
							},
						},
						Volumes: []corev1.Volume{
							corev1.Volume{
								Name: "tls",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: se.Name,
										Items: []corev1.KeyToPath{
											corev1.KeyToPath{
												Key:  "tls.key",
												Path: "server.key",
											},
											corev1.KeyToPath{
												Key:  "tls.crt",
												Path: "server.crt",
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
		err = sdk.Create(de)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}
