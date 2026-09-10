# Notiflex 아키텍처 스냅샷

> **6장 완료 시점** (2026-09-10) 기준

---

## 클러스터 개요

| 항목 | 값 |
|------|-----|
| 클러스터 | notiflex-cluster |
| GCP 프로젝트 | yunsoo-gitaiops-project |
| 리전/존 | asia-northeast3-a (서울) |
| K8s 버전 | v1.35.7-gke.1150000 |
| Workload Identity | yunsoo-gitaiops-project.svc.id.goog |

## 노드풀

| 노드풀 | 머신 타입 | 노드 수 | 특이사항 |
|--------|----------|---------|---------|
| default-pool | e2-medium (2vCPU, 4GB, Spot) | 2 | GKE_METADATA 활성화 (Workload Identity용) |

---

## 네임스페이스별 컴포넌트

### notiflex (앱)

```
외부 트래픽
    │
    ▼
[Gateway: notiflex-gateway]          gke-l7-regional-external-managed
  IP: 35.216.103.153
    │
    ▼
[HTTPRoute: notiflex-route]          /health, /id, /version → notiflex-api:80
    │
    ▼
[Service: notiflex-api]              ClusterIP 34.118.232.85:80
[Service: notiflex-api-canary]       ClusterIP 34.118.229.157:80  ← Canary용
    │
    ▼
[Rollout: notiflex-api]              Argo Rollouts, Canary 전략
  replicas: 4
  stable:  sha-1f27526 (v0.8.0)
  steps:   25% → 120s → 50% → 120s → 75% → 120s → 100%
    │
    ├── [Pod x4] notiflex-api-69df88bd8f-*
    │     serviceAccountName: notiflex-api
    │     env:
    │       POD_NAME         (fieldRef)
    │       VALKEY_PASSWORD_FILE=/mnt/secrets/valkey-password
    │     volumeMounts:
    │       /mnt/secrets  ← CSI 볼륨 마운트
    │         valkey-addr      (GCP Secret Manager에서 주입)
    │         valkey-password  (GCP Secret Manager에서 주입)
    │
    ├── [StatefulSet: valkey-primary]  Bitnami Valkey 9.1.2, standalone
    │     auth: enabled (비밀번호 인증)
    │     PVC: valkey-data-valkey-primary-0 (1Gi, standard-rwo)
    │     Service: valkey-primary:6379
    │
    ├── [ServiceAccount: notiflex-api]
    │     annotation: iam.gke.io/gcp-service-account
    │                 = notiflex-api-sa@yunsoo-gitaiops-project.iam.gserviceaccount.com
    │
    └── [SecretProviderClass: notiflex-secrets]
          provider: gcp
          secrets:
            notiflex-valkey-addr     → /mnt/secrets/valkey-addr
            notiflex-valkey-password → /mnt/secrets/valkey-password
          syncSecret → K8s Secret: notiflex-secrets
```

### argocd (GitOps)

```
[GitHub: roddbstn/notiflex-platform]
    │  push → CI (GitHub Actions)
    │          └── 이미지 빌드 → Artifact Registry
    │              → rollout.yaml 이미지 태그 자동 커밋
    ▼
[ArgoCD: notiflex-smb]
  감시 경로: k8s/smb/
  syncPolicy: automated (selfHeal + prune)
    │
    └── 클러스터 상태를 Git과 자동 동기화
```

### monitoring

```
[Prometheus]  ← kube-prometheus-stack-90.0.0
  수집: K8s 메트릭, Pod 메트릭
  알림 규칙: PrometheusRule (PodRestartTooMany 등)
    │
    ├── [Alertmanager]  알림 라우팅
    ├── [Grafana]       대시보드 (메트릭 + 로그 통합)
    │
[Loki]         로그 저장 (SingleBinary)
[Fluent Bit]   로그 수집 DaemonSet → Loki 전달
```

---

## 시크릿 관리 흐름

```
GCP Secret Manager
  ├── notiflex-valkey-addr      (Valkey 접속 주소)
  └── notiflex-valkey-password  (Valkey 인증 비밀번호)
          │
          │  IAM: roles/secretmanager.secretAccessor
          ▼
  GCP SA: notiflex-api-sa@yunsoo-gitaiops-project.iam.gserviceaccount.com
          │
          │  Workload Identity (roles/iam.workloadIdentityUser)
          ▼
  K8s SA: notiflex-api (notiflex namespace)
          │
          │  SecretProviderClass: notiflex-secrets
          ▼
  Secret Store CSI Driver (secrets-store-csi-driver-1.6.1)
  + GCP Provider DaemonSet
          │
          │  volumeMount: /mnt/secrets/
          ▼
  Pod 내 파일로 마운트
    /mnt/secrets/valkey-addr      → os.ReadFile() → Valkey 주소
    /mnt/secrets/valkey-password  → VALKEY_PASSWORD_FILE env → os.ReadFile()
```

---

## CI/CD 흐름

```
git push (main)
    │
    ▼
GitHub Actions (CI)
  1. Docker 빌드
  2. Artifact Registry 푸시 (sha-XXXXXXX 태그)
  3. k8s/smb/rollout.yaml 이미지 태그 자동 커밋
    │
    ▼
ArgoCD (자동 감지)
    │
    ▼
Argo Rollouts (Canary)
  25% → 120s → 50% → 120s → 75% → 120s → 100%
```

---

## 앱 코드 동작 (v0.8.0)

| 엔드포인트 | 동작 |
|-----------|------|
| `GET /health` | `{"status":"ok"}` 반환 |
| `GET /version` | `{"version":"v0.8.0"}` 반환 |
| `GET /id` | `/mnt/secrets/valkey-addr` 파일 읽기 → Valkey INCR (RESP 직접 구현) → 전역 단조 증가 ID 반환 |

Valkey 연결 시 `VALKEY_PASSWORD_FILE` 환경변수가 있으면 해당 파일에서 비밀번호를 읽어 AUTH 처리, 없으면 `VALKEY_PASSWORD` env 폴백.

---

## 주요 버전

| 컴포넌트 | 버전 |
|---------|------|
| Notiflex API | v0.8.0 (sha-1f27526) |
| Go | 1.25 |
| ArgoCD | v3.5.2 |
| Argo Rollouts | (클러스터 내 설치) |
| Valkey | 9.1.2 (Bitnami Chart 6.2.19) |
| kube-prometheus-stack | 90.0.0 (Prometheus v0.93.1) |
| Loki | 3.6.12 (Chart 7.3.0) |
| Fluent Bit | v2.1.0 (Chart 2.6.0) |
| Secret Store CSI Driver | 1.6.1 |
| GKE | v1.35.7-gke.1150000 |

---

## Artifact Registry

```
asia-northeast3-docker.pkg.dev/yunsoo-gitaiops-project/notiflex/api:<sha>
```
