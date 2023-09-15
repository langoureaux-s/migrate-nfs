package k8s

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

func DisableWorkLoadNFS(ctx context.Context, client kubernetes.Interface, excludeNamespaces []string, includeNamespaces []string, isDryRun bool) (err error) {

	log.Debugf("Dry run: %t", isDryRun)
	log.Debug("Browse Workloads to search them that use nfs-client as storage class")
	var (
		dryRun []string
	)

	if isDryRun {
		dryRun = []string{"All"}
	}

	// Process deployment
	deployments, err := getDeploymentsThatUseNfsPvc(ctx, client, excludeNamespaces, includeNamespaces)
	if err != nil {
		return err
	}
	for _, deployment := range deployments {
		log.Infof("Found deployment %s/%s that use NFS PVC, force replica to 0", deployment.Namespace, deployment.Name)

		deployment.Spec.Replicas = pointer.Int32(0)
		if _, err = client.AppsV1().Deployments(deployment.Namespace).Update(ctx, &deployment, v1.UpdateOptions{DryRun: dryRun}); err != nil {
			return errors.Wrapf(err, "Error when update deployment %s/%s", deployment.Namespace, deployment.Name)
		}
	}

	// Process statefulset
	statefulsets, err := getStatefulsetsThatUseNfsPvc(ctx, client, excludeNamespaces, includeNamespaces)
	if err != nil {
		return err
	}
	for _, statefulset := range statefulsets {
		log.Infof("Found statefulset %s/%s that use NFS PVC, force replica to 0", statefulset.Namespace, statefulset.Name)

		statefulset.Spec.Replicas = pointer.Int32(0)
		if _, err = client.AppsV1().StatefulSets(statefulset.Namespace).Update(ctx, &statefulset, v1.UpdateOptions{DryRun: dryRun}); err != nil {
			return errors.Wrapf(err, "Error when update deployment %s/%s", statefulset.Namespace, statefulset.Name)
		}
	}

	return nil

}

func EnableWorkLoadNFS(ctx context.Context, client kubernetes.Interface, excludeNamespaces []string, includeNamespaces []string, isDryRun bool) (err error) {

	log.Debugf("Dry run: %t", isDryRun)
	log.Debug("Browse Workloads to search them that use nfs-client as storage class")
	var (
		dryRun []string
	)

	if isDryRun {
		dryRun = []string{"All"}
	}

	// Process deployment
	deployments, err := getDeploymentsThatUseNfsPvc(ctx, client, excludeNamespaces, includeNamespaces)
	if err != nil {
		return err
	}
	for _, deployment := range deployments {
		log.Infof("Found deployment %s/%s that use NFS PVC, force replica to 1", deployment.Namespace, deployment.Name)

		deployment.Spec.Replicas = pointer.Int32(1)
		if _, err = client.AppsV1().Deployments(deployment.Namespace).Update(ctx, &deployment, v1.UpdateOptions{DryRun: dryRun}); err != nil {
			return errors.Wrapf(err, "Error when update deployment %s/%s", deployment.Namespace, deployment.Name)
		}
	}

	// Process statefulset
	statefulsets, err := getStatefulsetsThatUseNfsPvc(ctx, client, excludeNamespaces, includeNamespaces)
	if err != nil {
		return err
	}
	for _, statefulset := range statefulsets {

		if statefulset.Spec.Replicas != nil && *statefulset.Spec.Replicas == 0 {
			log.Infof("Found statefulset %s/%s that use NFS PVC, force replica to 1", statefulset.Namespace, statefulset.Name)

			statefulset.Spec.Replicas = pointer.Int32(1)
			if _, err = client.AppsV1().StatefulSets(statefulset.Namespace).Update(ctx, &statefulset, v1.UpdateOptions{DryRun: dryRun}); err != nil {
				return errors.Wrapf(err, "Error when update statefulset %s/%s", statefulset.Namespace, statefulset.Name)
			}
		}
	}

	return nil

}

func getDeploymentsThatUseNfsPvc(ctx context.Context, client kubernetes.Interface, excludeNamespaces []string, includeNamespaces []string) (deployments []appv1.Deployment, err error) {
	deployments = make([]appv1.Deployment, 0)
	var isInclude bool

	log.Debugf("Exclude namespaces: %s", spew.Sdump(excludeNamespaces))
	log.Debugf("Include namespaces: %s", spew.Sdump(includeNamespaces))

	deploymentList, err := client.AppsV1().Deployments("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Error when list all deployment on all namespace")
	}

loopDeployment:
	for _, deployment := range deploymentList.Items {

		// Process exclude / include namespace
		for _, excludeNS := range excludeNamespaces {
			if excludeNS == deployment.Namespace {
				log.Infof("Not process deployment %s/%s because of it is on excluded namespace", deployment.Namespace, deployment.Name)
				continue loopDeployment
			}
		}
		if len(includeNamespaces) > 0 {
			for _, includeNS := range includeNamespaces {
				if includeNS == deployment.Namespace {
					isInclude = true
					break
				}
			}

			if !isInclude {
				log.Infof("Not process deployment %s/%s because of it is not on included namespace", deployment.Namespace, deployment.Name)
				continue loopDeployment
			}
		}

		for _, volume := range deployment.Spec.Template.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvc, err := client.CoreV1().PersistentVolumeClaims(deployment.Namespace).Get(ctx, volume.PersistentVolumeClaim.ClaimName, v1.GetOptions{})
				if err != nil {
					return nil, errors.Wrapf(err, "Error when get PVC %s/%s", deployment.Namespace, volume.PersistentVolumeClaim.ClaimName)
				}
				if pvc.Spec.StorageClassName == nil || *pvc.Spec.StorageClassName == "nfs-client" {
					deployments = append(deployments, deployment)
					break
				}
			} else if volume.NFS != nil {
				deployments = append(deployments, deployment)
				break
			}
		}
	}

	return deployments, nil
}

func getStatefulsetsThatUseNfsPvc(ctx context.Context, client kubernetes.Interface, excludeNamespaces []string, includeNamespaces []string) (statefulsets []appv1.StatefulSet, err error) {
	statefulsets = make([]appv1.StatefulSet, 0)
	var isInclude bool

	log.Debugf("Exclude namespaces: %s", spew.Sdump(excludeNamespaces))
	log.Debugf("Include namespaces: %s", spew.Sdump(includeNamespaces))

	statefulsetList, err := client.AppsV1().StatefulSets("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Error when list all statefulsets")
	}

loopStatefulset:
	for _, statefulset := range statefulsetList.Items {

		// Process exclude / include namespace
		for _, excludeNS := range excludeNamespaces {
			if excludeNS == statefulset.Namespace {
				log.Infof("Not process statefulset %s/%s because of it is on excluded namespace", statefulset.Namespace, statefulset.Name)
				continue loopStatefulset
			}
		}
		if len(includeNamespaces) > 0 {
			for _, includeNS := range includeNamespaces {
				if includeNS == statefulset.Namespace {
					isInclude = true
					break
				}
			}

			if !isInclude {
				log.Infof("Not process deployment %s/%s because of it is not on included namespace", statefulset.Namespace, statefulset.Name)
				continue loopStatefulset
			}
		}

		for _, pvcT := range statefulset.Spec.VolumeClaimTemplates {
			if pvcT.Spec.StorageClassName == nil || *pvcT.Spec.StorageClassName == "nfs-client" {
				statefulsets = append(statefulsets, statefulset)
				continue loopStatefulset
			}
		}

		for _, volume := range statefulset.Spec.Template.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvc, err := client.CoreV1().PersistentVolumeClaims(statefulset.Namespace).Get(ctx, volume.PersistentVolumeClaim.ClaimName, v1.GetOptions{})
				if err != nil {
					return nil, errors.Wrapf(err, "Error when get PVC %s/%s", statefulset.Namespace, volume.PersistentVolumeClaim.ClaimName)
				}
				if pvc.Spec.StorageClassName == nil || *pvc.Spec.StorageClassName == "nfs-client" {
					statefulsets = append(statefulsets, statefulset)
					continue loopStatefulset
				}
			} else if volume.NFS != nil {
				statefulsets = append(statefulsets, statefulset)
				continue loopStatefulset
			}
		}
	}

	return statefulsets, nil
}
