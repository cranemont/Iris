# Iris

> Online Judge 시스템을 위한 채점 서버
> >개발이 진행중인 프로젝트로, 현재는 RabbitMQ에만 연결 가능합니다.

Online Judge 시스템을 위한 채점 서버입니다. 사용자가 제출한 코드를 컴파일하고 주어진 테스트케이스에 대한 입/출력과 리소스 사용량 등을 확인하여 정답 여부를 반환합니다. 서버의 출력 결과는 connector(WIKI 참조 - [WIP](#)) 및 sandbox 바이너리에 따라 달라질 수 있습니다. 

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
# Tests
- WIP
---
# Roadmap
- WIP
---
# Authors and acknowledgment

### Authors
@cranemont, @mixxeo

### Sandbox 정보
- [Judger](https://github.com/QingdaoU/Judger)

---
# License

MIT License
