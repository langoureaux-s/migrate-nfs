package k8s

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNodesWhereNfsIsMounted(ctx context.Context, client kubernetes.Interface) (nodes []string, err error) {

	nodes = make([]string, 0)
	pods := make([]corev1.Pod, 0)

	podList, err := client.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Error when get all pods")
	}

	log.Debug("Browse pods to search them that use nfs-client as storage class")

loopPod:
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvc, err := client.CoreV1().PersistentVolumeClaims(pod.Namespace).Get(ctx, volume.PersistentVolumeClaim.ClaimName, v1.GetOptions{})
				if err != nil {
					return nil, errors.Wrapf(err, "Error when get PVC %s/%s", pod.Namespace, volume.PersistentVolumeClaim.ClaimName)
				}
				if pvc.Spec.StorageClassName == nil || *pvc.Spec.StorageClassName == "nfs-client" {
					pods = append(pods, pod)
					continue loopPod
				}
			} else if volume.NFS != nil {
				pods = append(pods, pod)
				continue loopPod
			}
		}
	}

	for _, pod := range pods {
		nodes = append(nodes, pod.Spec.NodeName)
	}

	nodes = funk.UniqString(nodes)
	sort.Strings(nodes)

	return nodes, nil

}
