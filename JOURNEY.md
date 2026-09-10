# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-09-08 | Claude Code 2.1.90, statusline 설정 완료 |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-09-08 | gcloud 583.0.0, 프로젝트: yunsoo-gitaiops-project, 리전: asia-northeast3 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-09-08 | roddbstn/notiflex-platform (public) |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-09-09 | notiflex-cluster, e2-medium x2, Spot VM, Gateway API 활성화 |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-09-09 | api:v0.1.0 빌드 1분 10초, Pod 2개 Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-09-09 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-09-08 | ArgoCD v3.5.2 설치, notiflex-smb Application 생성, selfHeal+prune 활성화 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-09-08 | /version 엔드포인트 추가(v0.1.1→v0.1.2→revert→v0.1.3), Rolling Update 확인 |
| ch3 | 3.4 CI | ✅ | 2026-09-08 | GitHub Actions CI, GCP_SA_KEY/GCP_PROJECT_ID Secrets 등록, sha 태그 자동 빌드 |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-09-08 | CI가 deployment.yaml 이미지 태그 자동 커밋+push → ArgoCD 자동 감지 배포 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-09-10 | kube-prometheus-stack 설치, Prometheus/Grafana/Alertmanager Running |
| ch4 | 4.3 로그 수집 | ✅ | 2026-09-10 | Loki SingleBinary + Fluent Bit DaemonSet, Grafana Explore에서 로그 확인 |
| ch4 | 4.4 알림 | ✅ | 2026-09-10 | PrometheusRule(PodRestartTooMany 5분/2회), Alertmanager 연동 |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-09-10 | Gateway API(gke-l7-regional-external-managed), 외부 IP 35.216.103.153 |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-09-10 | Argo Rollouts Blue/Green, autoPromotionSeconds:30, v0.3.0 배포 확인 |
| ch6 | 6.1 캐시 | ⬜ | | |
| ch6 | 6.2 시크릿 관리 | ⬜ | | |
| ch6 | 6.3 Canary 전환 | ⬜ | | |
| ch7 | 7.2 멀티 노드풀 | ⬜ | | |
| ch7 | 7.3 App of Apps | ⬜ | | |
| ch7 | 7.4 멀티테넌시 | ⬜ | | |
| ch8 | 8.1 메시징 | ⬜ | | |
| ch8 | 8.2 트레이싱 | ⬜ | | |
| ch8 | 8.3 CronJob | ⬜ | | |
| ch9 | 9.1 저장소 분석 | ⬜ | | |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ⬜ | | |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| 컨테이너 베이스 이미지 | scratch | Alpine, Distroless | 최소 크기, 보안 표면 최소화 |
| GitOps 도구 | ArgoCD | Flux, Spinnaker | GKE 친화적, UI 제공, 선언적 Application CRD, 활발한 생태계 |
| CI 도구 | GitHub Actions | Cloud Build, Jenkins | 저장소 네이티브 통합, Secrets 관리 간편, 별도 서버 불필요 |
| 메트릭 모니터링 | Prometheus + Grafana (kube-prometheus-stack) | Datadog, Google Cloud Monitoring | K8s 표준, 무료, Helm 번들, Loki·Tempo와 Grafana 통합 |
| 로그 수집 | Loki + Fluent Bit | ELK Stack, Google Cloud Logging | 경량(128Mi), e2-medium 적합, Grafana 네이티브 통합 |
| 알림 | PrometheusRule + Alertmanager | Grafana Alerting, PagerDuty | GitOps 호환 CRD, kube-prometheus-stack에 포함, git blame으로 변경 추적 |
| 외부 트래픽 | Gateway API (gke-l7-regional-external-managed) | Ingress NGINX, Istio | K8s 차세대 표준, GKE 네이티브(추가 설치 불필요), 역할 분리 |
| 배포 전략 | Argo Rollouts Blue/Green | Rolling Update, Flagger | 트래픽 전환 즉시, ArgoCD 동일 생태계, 6장 Canary로 진화 가능 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 초기 설정 |
| Notiflex 이미지 | sha-4280b31 (v0.3.0) | v0.1.0→v0.1.3→v0.2.0→v0.3.0, Blue/Green 배포 전환 |
| ArgoCD | v3.5.2 | ch3.2 설치 |
| Kafka | | |
| OTel SDK | | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium (Spot) | 2 | notiflex-api |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2 | Python 3.9로 gcloud 설치 실패 | brew install python@3.12 후 CLOUDSDK_PYTHON 설정 |
| ch2 | Cloud Build API 비활성화 | gcloud services enable cloudbuild.googleapis.com |
| ch3 | Dockerfile WORKDIR 오타(/appㅊ) | Edit 도구로 /app으로 수정 |
| ch3 | GitHub Actions 쓰기 권한 실패(403) | gh api로 default_workflow_permissions=write 설정 |
| ch3 | CI 빌드 후 로컬 push 충돌 | git pull origin main --no-rebase 후 push |
