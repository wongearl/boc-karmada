# scheduler.image.pullPolicyBOC集成Karmada方案设计

## 1. 初步目标

1. 6月3号前KONK替换现有BOC多集群应用发布方式。简化流程，同时考虑老版本升级。
2. 6月3号之前KONK,支持多集群应用发布，支持集群间应用的差异化配置。支持常见的应用多活模式，例如同城双活，两地三中心等。

## 2. karmada部署

证书有效时间设定为100年。

### 2.1 helm方式部署

#### 2.1.1 要求

* kubernetes 1.16+
* helmV3+

#### 2.1.2 部署步骤

##### 2.1.2.1 部署karmad控制平面

根据实际需求配置以下内容：

| Name                                     | Description                                                  | Value                                                        | 需要修改的值                      |
| ---------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | --------------------------------- |
| `installMode`                            | InstallMode "host", "agent" and "component" are provided, "host" means install karmada in the control-cluster, "agent" means install agent client in the member cluster, "component" means install selected components in the control-cluster | `"host"`                                                     |                                   |
| `clusterDomain`                          | Default cluster domain for karmada                           | `"cluster.local"`                                            |                                   |
| `certs.mode`                             | Mode "auto" and "custom" are provided, "auto" means auto generate certificate, "custom" means use user certificate | `"auto"`                                                     |                                   |
| `certs.auto.expiry`                      | Expiry of the certificate                                    | `"876000h"`                                                  |                                   |
| `certs.auto.hosts`                       | Hosts of the certificate                                     | `["kubernetes.default.svc","*.etcd.karmada-system.svc.cluster.local","*.karmada-system.svc.cluster.local","*.karmada-system.svc","localhost","127.0.0.1"]` |                                   |
| `certs.custom.caCrt`                     | CA CRT of the certificate                                    | `""`                                                         |                                   |
| `certs.custom.crt`                       | CRT of the certificate                                       | `""`                                                         |                                   |
| `certs.custom.key`                       | KEY of the certificate                                       | `""`                                                         |                                   |
| `certs.custom.frontProxyCaCrt`           | CA CRT of the front proxy certificate                        | `""`                                                         |                                   |
| `certs.custom.frontProxyCrt`             | CRT of the front proxy certificate                           | `""`                                                         |                                   |
| `certs.custom.frontProxyKey`             | KEY of the front proxy certificate                           | `""`                                                         |                                   |
| `etcd.mode`                              | Mode "external" and "internal" are provided, "external" means use external ectd, "internal" means install a etcd in the cluster | `"internal"`                                                 |                                   |
| `etcd.external.servers`                  | Servers of etcd                                              | `""`                                                         |                                   |
| `etcd.external.registryPrefix`           | Use to registry prefix of etcd                               | `"/registry/karmada"`                                        |                                   |
| `etcd.external.certs.caCrt`              | CA CRT of the etcd certificate                               | `""`                                                         |                                   |
| `etcd.external.certs.crt`                | CRT of the etcd certificate                                  | `""`                                                         |                                   |
| `etcd.external.certs.key`                | KEY of the etcd certificate                                  | `""`                                                         |                                   |
| `etcd.internal.replicaCount`             | Target replicas of the etcd                                  | `1`                                                          |                                   |
| `etcd.internal.image.repository`         | Image of the etcd                                            | `"k8s.gcr.io/etcd"`                                          | 修改etcd镜像仓库地                |
| `etcd.internal.image.pullPolicy`         | Image pull policy of the etcd                                | `"IfNotPresent"`                                             |                                   |
| `etcd.internal.image.tag`                | Image tag of the etcd                                        | `"3.4.13-0"`                                                 | 修改etcd镜像tag                   |
| `etcd.internal.storageType`              | StorageType of the etcd, accepts "hostPath", "pvc"           | `"hostPath"`                                                 |                                   |
| `etcd.internal.pvc.storageClass`         | StorageClass of the etcd, takes effect when `etcd.internal.storageType` is "pvc" | `""`                                                         |                                   |
| `etcd.internal.pvc.size`                 | Storage size of the etcd, takes effect when `etcd.internal.storageType` is "pvc" | `""`                                                         |                                   |
| `etcd.internal.resources`                | Resource quota of the etcd                                   | `{}`                                                         |                                   |
| `scheduler.labels`                       | Labels of the scheduler deployment                           | `{"app": "karmada-scheduler"}`                               |                                   |
| `scheduler.replicaCount`                 | Target replicas of the scheduler                             | `1`                                                          |                                   |
| `scheduler.podLabels`                    | Labels of the scheduler pods                                 | `{}`                                                         |                                   |
| `scheduler.podAnnotations`               | Annotaions of the scheduler pods                             | `{}`                                                         |                                   |
| `scheduler.imagePullSecrets`             | Image pull secret of the scheduler                           | `[]`                                                         |                                   |
| `scheduler.image.repository`             | Image of the scheduler                                       | `"swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-scheduler"` | 修改karmada-scheduler镜像仓库地址 |
| `scheduler.image.tag`                    | Image tag of the scheduler                                   | `"latest"`                                                   | 修改镜像tag                       |
| `scheduler.image.pullPolicy`             | Image pull policy of the scheduler                           | `"IfNotPresent"`                                             |                                   |
| `scheduler.resources`                    | Resource quota of the scheduler                              | `{}`                                                         |                                   |
| `scheduler.nodeSelector`                 | Node selector of the scheduler                               | `{}`                                                         |                                   |
| `scheduler.affinity`                     | Affinity of the scheduler                                    | `{}`                                                         |                                   |
| `scheduler.tolerations`                  | Tolerations of the scheduler                                 | `[]`                                                         |                                   |
| `webhook.labels`                         | Labels of the webhook deployment                             | `{"app": "karmada-webhook"}`                                 |                                   |
| `webhook.replicaCount`                   | Target replicas of the webhook                               | `1`                                                          |                                   |
| `webhook.podLabels`                      | Labels of the webhook pods                                   | `{}`                                                         |                                   |
| `webhook.podAnnotations`                 | Annotaions of the webhook pods                               | `{}`                                                         |                                   |
| `webhook.imagePullSecrets`               | Image pull secret of the webhook                             | `[]`                                                         |                                   |
| `webhook.image.repository`               | Image of the webhook                                         | `"swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-webhook"` | 修改镜像仓库地址                  |
| `webhook.image.tag`                      | Image tag of the webhook                                     | `"latest"`                                                   | 修改镜像tag                       |
| `webhook.image.pullPolicy`               | Image pull policy of the webhook                             | `"IfNotPresent"`                                             |                                   |
| `webhook.resources`                      | Resource quota of the webhook                                | `{}`                                                         |                                   |
| `webhook.nodeSelector`                   | Node selector of the webhook                                 | `{}`                                                         |                                   |
| `webhook.affinity`                       | Affinity of the webhook                                      | `{}`                                                         |                                   |
| `webhook.tolerations`                    | Tolerations of the webhook                                   | `[]`                                                         |                                   |
| `controllerManager.labels`               | Labels of the karmada-controller-manager deployment          | `{"app": "karmada-controller-manager"}`                      |                                   |
| `controllerManager.replicaCount`         | Target replicas of the karmada-controller-manager            | `1`                                                          |                                   |
| `controllerManager.podLabels`            | Labels of the karmada-controller-manager pods                | `{}`                                                         |                                   |
| `controllerManager.podAnnotations`       | Annotaions of the karmada-controller-manager pods            | `{}`                                                         |                                   |
| `controllerManager.imagePullSecrets`     | Image pull secret of the karmada-controller-manager          | `[]`                                                         |                                   |
| `controllerManager.image.repository`     | Image of the karmada-controller-manager                      | `"swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-controller-manager"` | 修改镜像仓库地址                  |
| `controllerManager.image.tag`            | Image tag of the karmada-controller-manager                  | `"latest"`                                                   | 修改镜像tag                       |
| `controllerManager.image.pullPolicy`     | Image pull policy of the karmada-controller-manager          | `"IfNotPresent"`                                             |                                   |
| `controllerManager.resources`            | Resource quota of the karmada-controller-manager             | `{}`                                                         |                                   |
| `controllerManager.nodeSelector`         | Node selector of the karmada-controller-manager              | `{}`                                                         |                                   |
| `controllerManager.affinity`             | Affinity of the karmada-controller-manager                   | `{}`                                                         |                                   |
| `controllerManager.tolerations`          | Tolerations of the karmada-controller-manager                | `[]`                                                         |                                   |
| `apiServer.labels`                       | Labels of the karmada-apiserver deployment                   | `{"app": "karmada-apiserver"}`                               |                                   |
| `apiServer.replicaCount`                 | Target replicas of the karmada-apiserver                     | `1`                                                          |                                   |
| `apiServer.podLabels`                    | Labels of the karmada-apiserver pods                         | `{}`                                                         |                                   |
| `apiServer.podAnnotations`               | Annotaions of the karmada-apiserver pods                     | `{}`                                                         |                                   |
| `apiServer.imagePullSecrets`             | Image pull secret of the karmada-apiserver                   | `[]`                                                         |                                   |
| `apiServer.image.repository`             | Image of the karmada-apiserver                               | `"k8s.gcr.io/kube-apiserver"`                                | 修改镜像仓库地址                  |
| `apiServer.image.tag`                    | Image tag of the karmada-apiserver                           | `"v1.19.1"`                                                  | 修改镜像tag                       |
| `apiServer.image.pullPolicy`             | Image pull policy of the karmada-apiserver                   | `"IfNotPresent"`                                             |                                   |
| `apiServer.resources`                    | Resource quota of the karmada-apiserver                      | `{}`                                                         |                                   |
| `apiServer.hostNetwork`                  | Deploy karmada-apiserver with hostNetwork. If there are multiple karmadas in one cluster, you'd better set it to "false" | `"true"`                                                     |                                   |
| `apiServer.nodeSelector`                 | Node selector of the karmada-apiserver                       | `{}`                                                         |                                   |
| `apiServer.affinity`                     | Affinity of the karmada-apiserver                            | `{}`                                                         |                                   |
| `apiServer.tolerations`                  | Tolerations of the karmada-apiserver                         | `[]`                                                         |                                   |
| `apiServer.serviceType`                  | Service type of apiserver, accepts "ClusterIP", "NodePort", "LoadBalancer" | `"ClusterIP"`                                                |                                   |
| `apiServer.nodePort`                     | Node port for apiserver, takes effect when `apiServer.serviceType` is "NodePort". If no port is specified, the nodePort will be automatically assigned. | `0`                                                          |                                   |
| `aggregatedApiServer.labels`             | Labels of the karmada-aggregated-apiserver deployment        | `{"app": "karmada-aggregated-apiserver"}`                    |                                   |
| `aggregatedApiServer.replicaCount`       | Target replicas of the karmada-aggregated-apiserver          | `1`                                                          |                                   |
| `aggregatedApiServer.podLabels`          | Labels of the karmada-aggregated-apiserver pods              | `{}`                                                         |                                   |
| `aggregatedApiServer.podAnnotations`     | Annotaions of the karmada-aggregated-apiserver pods          | `{}`                                                         |                                   |
| `aggregatedApiServer.imagePullSecrets`   | Image pull secret of the karmada-aggregated-apiserver        | `[]`                                                         |                                   |
| `aggregatedApiServer.image.repository`   | Image of the karmada-aggregated-apiserver                    | `"swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-aggregated-apiserver"` | 修改镜像仓库地址                  |
| `aggregatedApiServer.image.tag`          | Image tag of the karmada-aggregated-apiserver                | `"latest"`                                                   | 修改镜像tag                       |
| `aggregatedApiServer.image.pullPolicy`   | Image pull policy of the karmada-aggregated-apiserver        | `"IfNotPresent"`                                             |                                   |
| `aggregatedApiServer.resources`          | Resource quota of the karmada-aggregated-apiserver           | `{requests: {cpu: 100m}}`                                    |                                   |
| `aggregatedApiServer.nodeSelector`       | Node selector of the karmada-aggregated-apiserver            | `{}`                                                         |                                   |
| `aggregatedApiServer.affinity`           | Affinity of the karmada-aggregated-apiserver                 | `{}`                                                         |                                   |
| `aggregatedApiServer.tolerations`        | Tolerations of the karmada-aggregated-apiserver              | `[]`                                                         |                                   |
| `kubeControllerManager.labels`           | Labels of the kube-controller-manager deployment             | `{"app": "kube-controller-manager"}`                         |                                   |
| `kubeControllerManager.replicaCount`     | Target replicas of the kube-controller-manager               | `1`                                                          |                                   |
| `kubeControllerManager.podLabels`        | Labels of the kube-controller-manager pods                   | `{}`                                                         |                                   |
| `kubeControllerManager.podAnnotations`   | Annotaions of the kube-controller-manager pods               | `{}`                                                         |                                   |
| `kubeControllerManager.imagePullSecrets` | Image pull secret of the kube-controller-manager             | `[]`                                                         |                                   |
| `kubeControllerManager.image.repository` | Image of the kube-controller-manager                         | `"k8s.gcr.io/kube-controller-manager"`                       | 修改镜像仓库地址                  |
| `kubeControllerManager.image.tag`        | Image tag of the kube-controller-manager                     | `"v1.19.1"`                                                  | 修改镜像tag                       |
| `kubeControllerManager.image.pullPolicy` | Image pull policy of the kube-controller-manager             | `"IfNotPresent"`                                             |                                   |
| `kubeControllerManager.resources`        | Resource quota of the kube-controller-manager                | `{}`                                                         |                                   |
| `kubeControllerManager.nodeSelector`     | Node selector of the kube-controller-manager                 | `{}`                                                         |                                   |
| `kubeControllerManager.affinity`         | Affinity of the kube-controller-manager                      | `{}`                                                         |                                   |
| `kubeControllerManager.tolerations`      | Tolerations of the kube-controller-manager                   | `[]`                                                         |                                   |

执行部署命令：

```console
$ helm install karmada -n karmada-system --create-namespace ./charts
```

部署完成后生成的karmada-system命名空间中名为karmada-kubeconfig的secret即karmada控制平面kubeconfig。

```
apiVersion: v1
kind: Config
clusters:
  - cluster:
      certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURtekNDQW9PZ0F3SUJBZ0lVU25VNzVVQlJOSEp6MlVjazM3aHhnZlcyYUJJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1hERUxNQWtHQTFVRUJoTUNlSGd4Q2pBSUJnTlZCQWdNQVhneENqQUlCZ05WQkFjTUFYZ3hDakFJQmdOVgpCQW9NQVhneENqQUlCZ05WQkFzTUFYZ3hDekFKQmdOVkJBTU1BbU5oTVJBd0RnWUpLb1pJaHZjTkFRa0JGZ0Y0Ck1DQVhEVEl5TURVeE1UQTROVEF6TVZvWUR6SXhNakl3TkRFM01EZzFNRE14V2pCY01Rc3dDUVlEVlFRR0V3SjQKZURFS01BZ0dBMVVFQ0F3QmVERUtNQWdHQTFVRUJ3d0JlREVLTUFnR0ExVUVDZ3dCZURFS01BZ0dBMVVFQ3d3QgplREVMTUFrR0ExVUVBd3dDWTJFeEVEQU9CZ2txaGtpRzl3MEJDUUVXQVhnd2dnRWlNQTBHQ1NxR1NJYjNEUUVCCkFRVUFBNElCRHdBd2dnRUtBb0lCQVFDMUZ1OTZZMmUwdU1ZMzBOKzFBdmF4Z3llQTJwUUR4dGwyc1JQdkI4SjQKYlV1ZnQrNERrNXpocFI2di9mdHRvbkVWbzBXaENyQVRIb2NaUGthYzJFSHRaeXgyTzB3ekhkbDMzeFFDZ3F2dgozelZQSHREd3ZIbnZxV056NjAwMHNnTjRBQWlTR2FIQlRaeTQyd1FvOTRaWktGcXN3NGtBTTl3dERucWh6dWVUCi94R2JVWjVydkcrMFJHVXRaTjBLb2wzNFFoZHJIMUtOSHRidllqT0VRVE9waDBsWVVYcnlnWklEeXd1QTNtcmMKdEVnUlZwcGRPWEtCOHNwaEJDT3Vady9EYjUwTXM3bk9nWXdka2JSaStVcVBOZzFTUkRCV1V4OEowK2xKSThSbQp3TW9oSGZxazF5b0JWaDhBTDlRcDRlS3l0bWR5dFZyblNXbFFyaU9JeGxCbEFnTUJBQUdqVXpCUk1CMEdBMVVkCkRnUVdCQlNNeGFkRFc3dTRlZUlmanNub3pFS3hBREp2aERBZkJnTlZIU01FR0RBV2dCU014YWREVzd1NGVlSWYKanNub3pFS3hBREp2aERBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBYwpDV0pVRzBqbk15Uk9MaGx4TUZHaWRjbUVlTFp1OXZybzc3d1NEQjRwd1BRUElMemdqYUNLNVNoN2FNWEdWUUtZCm5VRVN3dHE5OG9oWkUrVzlDb1dUQkNsc0Z6UDRnN0VOSTF2MkFvSEppOXdaY2NNY1BtM1liOWU5anArbkxpM3UKVXhmS1ZRZk1kUHpmSlo3dGdtcjEvcXY5bkxIdTR0NjBrUm9BLzZqWUk4WklaYXRNYzVaY3JKZXd2eG96bm82SAp6Y3ZSTVF3RGFsUlNtREVTTER1WEhZNXc1RTlaUFMyU0E4ZWdvTDJ4M2VQcmc2eWNsQ3RDMnlSbVA3MElSeTN5CjJGamNqVEY1b3NqMlZIQml2ZzdTMjRtK2UrM3NhbWw0dlhZZTFGUXZrVE13RDNhb01oeDRxWUVvcW9xOHJYejAKNG9rY1pYUndGWlQrVytoeEpMZlkKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
      insecure-skip-tls-verify: false
      server: https://karmada-apiserver.karmada-system.svc.cluster.local:5443
    name: karmada-apiserver
users:
  - user:
      client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVPVENDQXlHZ0F3SUJBZ0lVV3E2RXVIOGIxNDFySXpxc3gwQ3o2dGRSeFZNd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1hERUxNQWtHQTFVRUJoTUNlSGd4Q2pBSUJnTlZCQWdNQVhneENqQUlCZ05WQkFjTUFYZ3hDakFJQmdOVgpCQW9NQVhneENqQUlCZ05WQkFzTUFYZ3hDekFKQmdOVkJBTU1BbU5oTVJBd0RnWUpLb1pJaHZjTkFRa0JGZ0Y0Ck1DQVhEVEl5TURVeE1UQTRORFl3TUZvWUR6SXhNakl3TkRFM01EZzBOakF3V2pBd01SY3dGUVlEVlFRS0V3NXoKZVhOMFpXMDZiV0Z6ZEdWeWN6RVZNQk1HQTFVRUF4TU1jM2x6ZEdWdE9tRmtiV2x1TUlJQklqQU5CZ2txaGtpRwo5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBMVhaT1B0anY1V3VGTE94aUduMnRibnRZZitUS25MTE9vWm5ZCkFnWWNONmQxeTMvbDJiaUV0QWhEZ09JTDBMUFRPeFJTa1pCbjZraGVuWVMweUowVDJlcGdITnFJbXRzWHQrRjMKdFIvWjlVOG44SFNEbDhjSjBDQzRmU3NneEZGNnNLRzdIYThucUdxTmdGdkllYVRHektpbE9LRnVTM3h3ZGg2QgpLcHN5SGxHVkwvc3dVeWpWc0owclVraSt2WG5XVEt4T1UvMDFvbzJORVpvTlJOTklid3VKSjJ0OElveGZkcERPCnhQVUpVYnhULzZOYkhRbkY2ZGswWEhncHNteGNsTzdySk9KOEQ1TitZVXVOTEFYWnVuVlVTU0NBVjdaVkZvYnkKSG91K3JVTC9rSFo3OW9LYVdEYU5kRDViVVNEL0NlYVBMNzBZREU2R2k0amV5dW1vT3dJREFRQUJvNElCR3pDQwpBUmN3RGdZRFZSMFBBUUgvQkFRREFnV2dNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01DQmdnckJnRUZCUWNECkFUQU1CZ05WSFJNQkFmOEVBakFBTUIwR0ExVWREZ1FXQkJUYVRFU3JTblpTVS9MK3ZpeTQ3OXp1Sit3cWhEQWYKQmdOVkhTTUVHREFXZ0JTTXhhZERXN3U0ZWVJZmpzbm96RUt4QURKdmhEQ0Jsd1lEVlIwUkJJR1BNSUdNZ2hacgpkV0psY201bGRHVnpMbVJsWm1GMWJIUXVjM1pqZ2ljcUxtVjBZMlF1YTJGeWJXRmtZUzF6ZVhOMFpXMHVjM1pqCkxtTnNkWE4wWlhJdWJHOWpZV3lDSWlvdWEyRnliV0ZrWVMxemVYTjBaVzB1YzNaakxtTnNkWE4wWlhJdWJHOWoKWVd5Q0ZDb3VhMkZ5YldGa1lTMXplWE4wWlcwdWMzWmpnZ2xzYjJOaGJHaHZjM1NIQkg4QUFBRXdEUVlKS29aSQpodmNOQVFFTEJRQURnZ0VCQUQxQ3dpeVhWSktNQkJHcW9WWldpZGlGeGplN2pjRlhRbFRWd3MvZFY0RmM5MzA4ClhJM05lN3ViaURCeXJVS25UZzFUSFRnRTh3Vi9FT3VEeWFDTGZwRnBKZFhHSmdsQ1orWVlIemRSdDFMdGpCc2MKelJZZHFVdmFFKyt6ems0VVNqZEdMWU1TSnBkcEgvRG1sZ0Y0RXRQTFNvN0pvWmIrcnJmdENkWDB1OG9LUFZxQQo2dTgyN3A2RGJlSlRkaEdqa2dScThERmZYODZIbGpIaklFMlhFMnpYcEdYTDNlVk8wRm12Rk9icm41YUFJUkN6CjJxOGNBYm4zRkNxWmVRdjVlamZrU0pQZ1QzN3c5c3JIMDdTc1prL3ZPUTVMZk5nQVRISzdlN2IwM2JlNkRINkcKRWNHYm5ubU5EUU9WL2FXcmUwQm9qMlRUWFl4OXcybU9uNy84cG5zPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
      client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBMVhaT1B0anY1V3VGTE94aUduMnRibnRZZitUS25MTE9vWm5ZQWdZY042ZDF5My9sCjJiaUV0QWhEZ09JTDBMUFRPeFJTa1pCbjZraGVuWVMweUowVDJlcGdITnFJbXRzWHQrRjN0Ui9aOVU4bjhIU0QKbDhjSjBDQzRmU3NneEZGNnNLRzdIYThucUdxTmdGdkllYVRHektpbE9LRnVTM3h3ZGg2Qktwc3lIbEdWTC9zdwpVeWpWc0owclVraSt2WG5XVEt4T1UvMDFvbzJORVpvTlJOTklid3VKSjJ0OElveGZkcERPeFBVSlVieFQvNk5iCkhRbkY2ZGswWEhncHNteGNsTzdySk9KOEQ1TitZVXVOTEFYWnVuVlVTU0NBVjdaVkZvYnlIb3UrclVML2tIWjcKOW9LYVdEYU5kRDViVVNEL0NlYVBMNzBZREU2R2k0amV5dW1vT3dJREFRQUJBb0lCQUF3ci85QXpuSkpQYkR1Zwppd09Kc1E4QXQ0NHJaS1pFeCtXTkdUVWNWaFdTVmZReHFkQ2RaZXZDSU45RGhIcjFGaEZqV2tYMG53aEw4aUUzClJQdS9nVGRHMXc0dUkrWDRva1NZOWJOOVNuZGplUnFMK2tqNFQ0WHZwN1Z6ZFIzY3E0dFEzWk9XdmtNck9FQmUKZTNMOFExMitMQldybHRkMUZCQ1lNck5VUUNwY1Q4disyaXNIVk9zdVpwQVk4bEo2VTJLVTJHYkEwTnBhWEYzOQp5Q1BpNDhrb1ZvaHN5OGxNTTA5aUhmaUxac1YzcE96aHBjRzFEWUcxclkyUHFtKzJKdWVkMXhCNFZOZW1EdWF6CjdQNGxKN3dYZFhGbDNKZWRKRTVSWk05c3FDalR3TmZWZElQV1d6QWtOaHFGSjMvazRGK29DZ1RUaTVMVnF6bkwKaVJuMXBNRUNnWUVBL3VzbmZtZy90dWFNMWFlTHEvOVN2Zno5VzFIM2hWT2g3aE9CY2V0QmtkT3lmNmlFdHJiQQpOekJtdnpvbXFTTVAxT3NKOHhtU0QvRUdtTW5aTHdoT0JsaVg1TWM1bitXVkJKL3VmampINUtrczQrZVAxcUErCjNDVHV4ejZNKytLWXRjaDcwZXloUitEQ0IrRml1YnkrTWhpZDZCWEs0N0tKejJXSUhQQi96cHNDZ1lFQTFsNGgKQm1ReWZlZ0w5OGc2bk5GdEx5RW5QbVBVemNEbGl4b3NNS1NsenJlUUlyaXlTRk42dTVZbDI1SWNUOUdkNE16RQpoU3dJL2xRZ1JmNmpGdTVPN3IzbUZWaldVMHRMT3lpdCs0OENwQlFsbWkydnc3YzEvT3R6R1JPTCtma01hS0xsCjRqRWtVVVNaQkllK0dzS1p1UGluYW8ySzMzS3Y5RndBSmRYWFZ1RUNnWUFEcmYxSjg5TkhucGhWMVM5TThraU4KZVlObXVBNHNuSUp2MzFRMUFzSlZpb3EvRVAycGJZWGt4Z3dWb285QVRjTkN2WW5OT29kRE42Vnl1eWNwYUtOSApzQ3V4SDdjVE9jc20zL0FmWGs1MFhJVExYV0pVSk1nRGdYejQzMHhGcm9XcHloVVBlS3p6VHFrK1YrQ0c4ZFNGCkFKbjI2YW1lRTh1dGRMZThRelIxRlFLQmdBV205YmYxYkY5bGZ1dEpuRUlHUTVxNmhRNWdFM3haV2JRUlBKa0wKdmxOMUE1Zm14c3loWnRzTFduUXZwZzkwdDNUMThUaVJzS3NFRE5YTG9RRTV4MXNFSnN6cFNyMW5mdFJRZmtYagp2MjRVR2VtMnlxdWVhUTRDSjBiQi9TN2FJam1nRWUrazNCQklmc0JmMCtOZ3ZpemlZWWV0czd1d2luTW0rZG9GCmlXZ0JBb0dBSWEvZUhEaTB0dWZmZDFQMTZ2Z0N2QlpRTGFXY0pkVHhoWWVhdDN3a3BwUHk4bCthd0lITFd1ZE8KRE9udnhjWnJiVTJLS3hSWHIweUJ1WmQ5NnFLL0ZLODZXQ01BOVpEdk1GZ0ljaER6aFhla3RGdFQyUGNMQy9nQgpzdlphdDFjbGROc1JsUGlkRUVpeWNzelg5RXY5alg5VlVFYSs0bFVMdWxIY09IVTZHUFE9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
    name: karmada-apiserver
contexts:
  - context:
      cluster: karmada-apiserver
      user: karmada-apiserver
    name: karmada-apiserver
current-context: karmada-apiserver
```

卸载：

```
$ helm uninstall karmada -n karmada-system
```



## 3. 添加成员集群

在boc新建集群时，选择是否需要将集群加入karmada作为成员集群，如果需要加入karmada作为成员集群则选择是pull模式还是push模式。

### 3.1 push模式

通俗讲此种模式需要karmada控制平面需要能够连接到各个成员集群的API server。

在push模式下karmada控制平面将直接访问成员集群的kube-apiserver获取成员集群的状态和部署清单，并且由karmada控制面的execution controller将应用推送到成员集群。



```go
type ClusterRegistration struct{
	// pull or push 
	SyncMode string   
    // 保存karmada控制平面kubeconfig的secret
	KarmadaKubeConfig string
	KarmadaContext string
    // 保存成员集群kubeconfig的secret
	ClusterKubeConfig string
	ClusterContext string
    // 设置成员集群名
	MemberClusterName string
}
```



### 3.2 pull模式

在pull模式下karmada 控制平面不会访问成员集群，而是将其委托给一个名为`karmada-agent`.

每个`karmada-agent`服务于一个集群并负责：

- 将集群注册到 Karmada（创建`Cluster`对象）
- 维护集群状态并向 Karmada 报告（更新`Cluster`对象状态）
- 从 Karmada 执行空间（命名空间，`karmada-es-<cluster name>`）观察清单并部署到为其服务的集群。

举个例子：

1. 在成员集群创建

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: karmada-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: karmada-agent
rules:
  - apiGroups: ['*']
    resources: ['*']
    verbs: ['*']
  - nonResourceURLs: ['*']
    verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: karmada-agent
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: karmada-agent
subjects:
  - kind: ServiceAccount
    name: karmada-agent-sa
    namespace: karmada-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: karmada-agent-sa
  namespace: karmada-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: karmada-agent
  namespace: karmada-system
  labels:
    app: karmada-agent
spec:
  replicas: 1
  selector:
    matchLabels:
      app: karmada-agent
  template:
    metadata:
      labels:
        app: karmada-agent
    spec:
      serviceAccountName: karmada-agent-sa
      tolerations:
        - key: node-role.kubernetes.io/master
          operator: Exists
      containers:
        - name: karmada-agent
          image: registry.cn-hangzhou.aliyuncs.com/earl-k8s/karmada-agent:v1.1.1-18-gdcd27196
          command:
            - /bin/karmada-agent
            - --karmada-kubeconfig=/etc/kubeconfig/karmada-kubeconfig
            - --karmada-context=karmada
            - --cluster-name=member2
            - --cluster-status-update-frequency=10s
            - --v=4
          volumeMounts:
            - name: kubeconfig
              mountPath: /etc/kubeconfig
      volumes:
        - name: kubeconfig
          secret:
            secretName: karmada-kubeconfig # 需要先在该成员集群的该命名空间创建保存了karmada控制平面kubeconfig的secret
```

