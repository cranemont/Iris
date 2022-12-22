# Iris

> Online Judge 시스템을 위한 채점 서버

---

# Installation

## Requirements

- Go 1.19 or newer
- Docker

### Optional

- RabbitMQ Server
- API Server for Testcase
  - Testcase는 환경 변수 `TESTCASE_SERVER_URL`로 지정된 URL에서 HTTP GET 요청을 통해 가져옵니다.
    - Request URL은 `$TESTCASE_SERVER_URL/testcase_id`입니다. [RabbitMQ Connector](#)를 사용할 경우 `testcase_id`는 [MessageID](https://www.rabbitmq.com/consumers.html#message-properties)로 전달해야 합니다
  - `TESTCASE_SERVER_URL`을 설정하지 않을 경우 `testcase/preset-source.go`의 데이터를 사용합니다.
  - 예시 서버 추가 예정(WIP)
- Redis Server
  - 캐싱된 테스트케이스를 여러 채점 서버에서 공유하기 위해 Redis를 사용합니다
  - 캐시 서버가 없는 경우 테스트케이스는 메모리에 캐싱됩니다. 자세한 사용법은 [WIKI](#)를 참고하세요

##

---

# Usage

## With RabbitMQ

- WIP

## With Others

- [WIKI](#)

---

# Authors and acknowledgment

@cranemont, @mixxeo

## [Judger](https://github.com/QingdaoU/Judger)

# License

MIT License
