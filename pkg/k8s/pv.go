package k8s

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func MigrateNFSPV(ctx context.Context, client kubernetes.Interface, isDryRun bool) (err error) {

	log.Debugf("Dry run: %t", isDryRun)
	log.Debug("Browse PV to search then of type NFS")
	var dryRun []string

	if isDryRun {
		dryRun = []string{"All"}
	}

	pvList, err := client.CoreV1().PersistentVolumes().List(ctx, v1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "Error when list all PV")
	}

	for _, pv := range pvList.Items {
		if pv.Spec.NFS != nil {
			if pv.Spec.ClaimRef != nil {
				log.Infof("Force NFSv3 on PV %s, linked to PVC %s/%s", pv.Name, pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
			} else {
				log.Infof("Force NFSv3 on PV %s, PV is released", pv.Name)
			}

			pv.Spec.MountOptions = []string{
				"nfsvers=3",
				"nolock",
			}

			if _, err = client.CoreV1().PersistentVolumes().Update(ctx, &pv, v1.UpdateOptions{DryRun: dryRun}); err != nil {
				return errors.Wrapf(err, "Error when force NFSv3 on PV %s", pv.Name)
			}
		}
	}

	return nil
}

func FixNFSPV(ctx context.Context, client kubernetes.Interface, isDryRun bool) (err error) {

	log.Debugf("Dry run: %t", isDryRun)
	log.Debug("Browse PV to search then of type NFS")
	var dryRun []string

	if isDryRun {
		dryRun = []string{"All"}
	}

	pvList, err := client.CoreV1().PersistentVolumes().List(ctx, v1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "Error when list all PV")
	}

	for _, pv := range pvList.Items {
		if pv.Spec.NFS != nil {
			if pv.Spec.ClaimRef != nil {
				log.Infof("Fix nfsv3 on PV %s, linked to PVC %s/%s", pv.Name, pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
			} else {
				log.Infof("Fix nfsv3 on PV %s, PV is released", pv.Name)
			}

			pv.Spec.MountOptions = []string{
				"nfsvers=3",
				"nolock",
			}

			if _, err = client.CoreV1().PersistentVolumes().Update(ctx, &pv, v1.UpdateOptions{DryRun: dryRun}); err != nil {
				return errors.Wrapf(err, "Error when force NFSv3 on PV %s", pv.Name)
			}
		}
	}

	return nil
}
