# Notiflex Platform

B2B 알림 SaaS 플랫폼 **Notiflex**의 인프라 및 배포 저장소입니다.

## 개요

이 저장소는 Notiflex 서비스의 애플리케이션 코드, 쿠버네티스 매니페스트, CI/CD 파이프라인을 관리합니다.

## 기술 스택

| 영역 | 기술 |
|------|------|
| 언어 | Go (표준 라이브러리) |
| 컨테이너 | scratch 베이스 이미지 |
| 오케스트레이션 | Kubernetes (GKE Standard, Zonal) |
| CI/CD | GitHub Actions + ArgoCD |
| 이미지 저장소 | Artifact Registry (asia-northeast3) |
| 모니터링 | Prometheus + Grafana |
| 로깅 | Loki + Fluent Bit |

## 디렉터리 구조

```
notiflex-platform/
├── app/                  # Go 애플리케이션 소스
├── k8s/
│   └── smb/              # Kubernetes 매니페스트
└── .github/
    └── workflows/        # GitHub Actions CI/CD 파이프라인
```

## 인프라

- **GCP 프로젝트**: `yunsoo-gitaiops-project`
- **리전**: `asia-northeast3` (서울)
- **클러스터**: GKE Standard, Zonal (`asia-northeast3-a`)
- **Artifact Registry**: `asia-northeast3-docker.pkg.dev/yunsoo-gitaiops-project/notiflex`
