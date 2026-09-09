# Notiflex Platform

Notiflex — B2B 알림 SaaS 플랫폼의 인프라 및 배포 저장소이다.

## 프로젝트 개요

- **서비스**: Notiflex — 기업 대상 알림 발송 SaaS (이메일, SMS, 푸시 등)
- **이 저장소의 역할**: 애플리케이션 코드 + 쿠버네티스 매니페스트 + CI/CD 파이프라인

## 기술 스택

- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (최소 크기)
- **오케스트레이션**: Kubernetes (GKE Standard, Zonal)
- **CI/CD**: GitHub Actions + ArgoCD (예정)

## GCP 설정

- **프로젝트 ID**: `yunsoo-gitaiops-project`
- **리전**: `asia-northeast3` (서울)
- **존**: `asia-northeast3-a`
- **Artifact Registry**: `asia-northeast3-docker.pkg.dev/yunsoo-gitaiops-project/notiflex`

## 디렉터리 구조

```
notiflex-platform/
├── CLAUDE.md
├── app/           # Go 애플리케이션 소스
├── k8s/
│   └── smb/       # Kubernetes 매니페스트
└── .github/
    └── workflows/ # GitHub Actions CI/CD 파이프라인
```

## 행동 규칙

1. 명령 실행 전 현재 상태를 확인한다 (`kubectl get`, `gcloud config list` 등).
2. 파일 변경 전 기존 내용을 먼저 읽는다.
3. 에러 발생 시 원인을 분석하고 해결 방안을 제시한 뒤 진행한다.
4. 매니페스트 작성 시 네임스페이스 (`notiflex`)를 명시한다.
5. **모든 `kubectl` 명령에 `--context gke-sysnet4admin_book_gitaiops`를 반드시 지정한다** (잘못된 클러스터 대상 실행 방지).
6. 리소스 생성/삭제 전에는 영향 범위를 먼저 설명한다.
7. 이미지 태그는 `latest`를 쓰지 않고 명시적 버전 (`v0.1.0` 등)을 사용한다.
8. 토큰, 키, 비밀번호는 코드/매니페스트에 하드코딩하지 않는다 (환경변수, GitHub Secrets, Secret Manager 사용).

> **행동 규칙 5번이 특히 중요합니다.** 이 책은 `claude --dangerously-skip-permissions` 모드로 진행하므로 kubectl 명령이 승인 없이 바로 실행됩니다. 만약 로컬에 minikube, kind 같은 다른 클러스터가 남아있고 현재 컨텍스트가 거기를 가리키고 있다면, `kubectl delete` 하나로 엉뚱한 클러스터를 지우는 사고가 날 수 있습니다. `--context`를 항상 지정하면 이런 실수를 원천 차단합니다.
