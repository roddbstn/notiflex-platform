# Architecture Decision Records

---

## ADR-001: GitOps 도구 — ArgoCD (3장)
**시점**: 2026-09 / **결정**: GitOps 도구로 ArgoCD를 채택하고 Flux, Spinnaker는 사용하지 않는다.
**이유**:
- GKE 친화적이고 UI를 제공하여 배포 상태를 직관적으로 확인 가능
- 선언적 Application CRD로 GitOps 원칙에 맞는 선언적 관리
- selfHeal + prune 활성화로 Git 상태와 클러스터 상태를 자동으로 동기화
- 활발한 생태계, Argo Rollouts(5장)으로 자연스럽게 확장 가능

---

## ADR-002: CI 도구 — GitHub Actions (3장)
**시점**: 2026-09 / **결정**: CI 도구로 GitHub Actions를 채택하고 Cloud Build, Jenkins는 사용하지 않는다.
**이유**:
- GitHub 저장소 네이티브 통합 — 별도 서버 설치 불필요
- Secrets 관리(GCP_SA_KEY, GCP_PROJECT_ID)가 저장소 설정에서 바로 처리
- push 이벤트로 자동 트리거, 이미지 태그 업데이트까지 파이프라인 완결
- 무료 플랜으로 학습 환경에서 충분

---

## ADR-003: 메트릭 모니터링 — Prometheus + Grafana (4장)
**시점**: 2026-09 / **결정**: 메트릭 모니터링으로 kube-prometheus-stack(Prometheus + Grafana)을 채택하고 Datadog, Google Cloud Monitoring은 사용하지 않는다.
**이유**:
- Kubernetes 모니터링의 사실상 표준 (CNCF Graduated), 학습 가치가 높음
- SaaS 구독료 없이 자체 호스팅 — Datadog은 호스트당 $15+/월
- Helm 번들로 Prometheus, Grafana, Alertmanager, kube-state-metrics, node-exporter 6개 컴포넌트를 검증된 버전으로 한 번에 설치
- 이후 Loki(로그), Tempo(트레이스)를 같은 Grafana에 통합하여 도구 파편화 없음

---

## ADR-004: 로그 수집 — Loki + Fluent Bit (4장)
**시점**: 2026-09 / **결정**: 로그 수집으로 Loki + Fluent Bit를 채택하고 ELK Stack, Google Cloud Logging은 사용하지 않는다.
**이유**:
- e2-medium(4GB) 노드에서 ELK(Elasticsearch 최소 2GB+)는 리소스 부족 — Loki는 128Mi로 구동
- Fluent Bit DaemonSet으로 모든 노드의 컨테이너 로그 자동 수집
- ADR-003에서 설치한 Grafana에 데이터소스만 추가하면 메트릭과 같은 UI에서 로그 조회
- 라벨 기반 인덱싱으로 풀텍스트 인덱싱 대비 저장 비용 낮음

---

## ADR-005: 알림 — PrometheusRule + Alertmanager (4장)
**시점**: 2026-09 / **결정**: 알림 시스템으로 PrometheusRule + Alertmanager를 채택하고 Grafana Alerting, PagerDuty는 사용하지 않는다.
**이유**:
- GitOps 호환 — PrometheusRule CRD를 YAML로 Git에 관리, ArgoCD가 자동 적용
- ADR-003에서 설치한 kube-prometheus-stack에 Alertmanager가 이미 포함 — 추가 설치 불필요
- `git blame`으로 "이 알림이 언제 추가됐고 왜 임계값이 이렇지?"를 추적 가능
- 그루핑/억제/라우팅 트리로 다단계 알림 표현 가능

---

## ADR-006: 외부 트래픽 — Gateway API (5장)
**시점**: 2026-09 / **결정**: 외부 트래픽 관리로 Kubernetes Gateway API(GKE 네이티브)를 채택하고 Ingress NGINX, Istio는 사용하지 않는다.
**이유**:
- K8s 1.27 GA, Ingress를 대체하는 공식 차세대 표준
- GKE 네이티브 — `gke-l7-regional-external-managed` GatewayClass가 자동으로 로드밸런서 생성, 별도 Controller 설치 불필요
- Gateway(인프라)/HTTPRoute(앱)로 역할 분리 — Ingress는 단일 리소스에 모든 설정 혼재
- Istio는 리소스 1GB+ 과도, NGINX는 별도 설치 필요

---

## ADR-007: 배포 전략 — Argo Rollouts Blue/Green (5장)
**시점**: 2026-09 / **결정**: 배포 전략으로 Argo Rollouts Blue/Green을 채택하고 K8s 기본 Rolling Update, Flagger는 사용하지 않는다.
**이유**:
- Rolling Update의 한계 해소 — 새 버전 문제 발생 시 롤백 중에도 사용자가 에러를 받는 구간 존재, Blue/Green은 트래픽 전환이 찰나
- ADR-001에서 선택한 ArgoCD와 같은 Argo 생태계 — ArgoCD UI에서 Rollout 상태 확인 가능
- CRD 기반 YAML 선언 — GitOps 흐름과 일관성 유지
- 6장에서 Canary 전략으로 전환 시 strategy 필드만 수정하면 됨 (점진적 진화)
