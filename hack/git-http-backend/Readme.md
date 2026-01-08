## Start server

```bash
docker compose up -d 
```

## Access Web UI

http://git.local.gardener.cloud:5123/test/repo

> User: `test`  
> Password: `testtest`

## Clone

```bash
git clone http://test:testtest@git.local.gardener.cloud:5123/test/repo.git
```

## Configure Git Remote in Landscape Repo

`git-sync-secret.yaml`:
```yaml
stringData:
  password: testtest
  username: test
```

`gotk-sync.yaml`:
```yaml
  url: http://git.local.gardener.cloud:3000/test/repo
```

## Gardener Local Configurations

`helm-release.yaml`:
```yaml
apiVersion: helm.toolkit.fluxcd.io/v2
kind: HelmRelease
metadata:
  name: gardener-operator
  namespace: garden
spec:
  values:
    env:
      - name: GARDENER_OPERATOR_LOCAL
        value: "true"
    hostAliases:
      - hostnames:
          - api.virtual-garden.local.gardener.cloud
        ip: 10.2.10.2
        # config:
        # featureGates:
        # IstioTLSTermination: true
        # VPAInPlaceUpdates: true
```

`kustomization.yaml`:
```yaml
patches:
  - path: helm-release.yaml
```

### Local Provider Extension

```yaml
apiVersion: operator.gardener.cloud/v1alpha1
kind: Extension
metadata:
  annotations:
    security.gardener.cloud/pod-security-enforce: privileged
  name: provider-local
spec:
  deployment:
    admission:
      runtimeCluster:
        helm:
          ociRepository:
            ref: registry.local.gardener.cloud:5001/local-skaffold_gardener-extension-admission-local_charts_runtime:v1.135.0-dev-7f22b34644-dirty@sha256:20f0ff847cab9f95a58495831144fc045bd45ca8f3bc5250fa06e3d194289e42
      values:
        logLevel: debug
      virtualCluster:
        helm:
          ociRepository:
            ref: registry.local.gardener.cloud:5001/local-skaffold_gardener-extension-admission-local_charts_application:v1.135.0-dev-7f22b34644-dirty@sha256:2332a4e503770aa826983ccaf5c06c6a7034e14622a2613f90e5dd4861b94c93
    extension:
      helm:
        ociRepository:
          ref: registry.local.gardener.cloud:5001/local-skaffold_gardener-extension-provider-local_charts_extension:v1.135.0-dev-7f22b34644-dirty@sha256:082d82d8f851d62a0b25948adb824e43a548140b48b353174c5f082f72a2ed8a
      injectGardenKubeconfig: true
      runtimeClusterValues:
        logLevel: debug
      values:
        imageVectorOverwrite:
          images:
            - name: machine-controller-manager-provider-local
              ref: registry.local.gardener.cloud:5001/local-skaffold_machine-controller-manager-provider-local:v1.135.0-dev-7f22b34644-dirty@sha256:62d8769bc9497082b7df2c82caccefc5ba11a7461331e6f5886be892c694ed9a
        logLevel: debug
  resources:
    - kind: BackupBucket
      primary: true
      type: local
    - kind: BackupEntry
      primary: true
      type: local
    - kind: ControlPlane
      primary: true
      type: local
    - kind: DNSRecord
      primary: true
      type: local
    - kind: Infrastructure
      primary: true
      type: local
    - kind: Worker
      primary: true
      type: local
    - kind: OperatingSystemConfig
      primary: true
      type: local
    - clusterCompatibility:
        - shoot
        - seed
      kind: Extension
      lifecycle:
        delete: AfterKubeAPIServer
        migrate: AfterKubeAPIServer
        reconcile: BeforeKubeAPIServer
      primary: true
      type: local-ext-seed
      workerlessSupported: true
    - clusterCompatibility:
        - shoot
        - garden
      kind: Extension
      primary: true
      type: local-ext-shoot
      workerlessSupported: true
    - kind: Extension
      lifecycle:
        reconcile: AfterWorker
      primary: true
      type: local-ext-shoot-after-worker
```
