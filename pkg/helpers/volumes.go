package helpers

import (
	corev1 "k8s.io/api/core/v1"
)

// GetPodVolumeType is a function that retrieves the type of a pod volume
func GetPodVolumeType(source corev1.VolumeSource) string {
	switch {
		case source.HostPath != nil:
			return "hostPath"
		case source.EmptyDir != nil:
			return "emptyDir"
		case source.GCEPersistentDisk != nil:
			return "gcePersistentDisk"
		case source.AWSElasticBlockStore != nil:
			return "awsElasticBlockStore"
		case source.Secret != nil:
			return "secret"
		case source.NFS != nil:
			return "nfs"
		case source.ISCSI != nil:
			return "iscsi"
		case source.Glusterfs != nil:
			return "glusterfs"
		case source.PersistentVolumeClaim != nil:
			return "persistentVolumeClaim"
		case source.RBD != nil:
			return "rbd"
		case source.FlexVolume != nil:
			return "flexVolume"
		case source.Cinder != nil:
			return "cinder"
		case source.CephFS != nil:
			return "cephfs"
		case source.Flocker != nil:
			return "flocker"
		case source.DownwardAPI != nil:
			return "downwardAPI"
		case source.FC != nil:
			return "fc"
		case source.AzureFile != nil:
			return "azureFile"
		case source.AzureDisk != nil:
			return "azureDisk"
		case source.ConfigMap != nil:
			return "configMap"
		case source.VsphereVolume != nil:
			return "vsphereVolume"
		case source.Quobyte != nil:
			return "quobyte"
		case source.PhotonPersistentDisk != nil:
			return "photonPersistentDisk"
		case source.Projected != nil:
			return "projected"
		case source.PortworxVolume != nil:
			return "portworxVolume"
		case source.ScaleIO != nil:
			return "scaleIO"
		case source.StorageOS != nil:
			return "storageos"
		case source.CSI != nil:
			return "csi"
	}
	return ""

}

// GetPersistentVolumeType is a function that retrieves the type of a persistent volume
func GetPersistentVolumeType(source corev1.PersistentVolumeSource) string {
	switch {
		case source.HostPath != nil:
			return "hostPath"
		case source.GCEPersistentDisk != nil:
			return "gcePersistentDisk"
		case source.AWSElasticBlockStore != nil:
			return "awsElasticBlockStore"
		case source.NFS != nil:
			return "nfs"
		case source.ISCSI != nil:
			return "iscsi"
		case source.Glusterfs != nil:
			return "glusterfs"
		case source.RBD != nil:
			return "rbd"
		case source.FlexVolume != nil:
			return "flexVolume"
		case source.Cinder != nil:
			return "cinder"
		case source.CephFS != nil:
			return "cephfs"
		case source.Flocker != nil:
			return "flocker"
		case source.FC != nil:
			return "fc"
		case source.AzureFile != nil:
			return "azureFile"
		case source.AzureDisk != nil:
			return "azureDisk"
		case source.VsphereVolume != nil:
			return "vsphereVolume"
		case source.Quobyte != nil:
			return "quobyte"
		case source.PhotonPersistentDisk != nil:
			return "photonPersistentDisk"
		case source.PortworxVolume != nil:
			return "portworxVolume"
		case source.ScaleIO != nil:
			return "scaleIO"
		case source.StorageOS != nil:
			return "storageos"
		case source.CSI != nil:
			return "csi"
	}
	return ""

}