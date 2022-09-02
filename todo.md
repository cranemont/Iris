
- [x] handler에서 함수 reflection 따서 string map으로 register 할수 있는 함수 구현
  - reflection성능 문제때문에 사용 안함. 직접 mapping
- [x] judger.Judge에서 run/grader 고루틴 관리
- [ ] compile/run config 패키지 완성
- [ ] testcase Manager 구현
  - [x] cache manager
  - [x] task에 testcase struct 생성해서 넣기
   - [ ] submissionDTO로 넘어왔으면 그거 사용
- [x] grader
  - [x] whitespace trim
- [x] sandbox 연결
  - testcase I/O 관리
- [ ] Result code 정의
- [ ] logger
- [x] error handler
- [ ] test code 작성
- [ ] MaxOutputSize Testcase 크기 따라서 설정하도록 하기
- [ ] libjudger 실행 개수는 CPU 수에 맞출것. 되도록이면 Docker compose의 환경변수로 지정
- [ ] run 실패시 grade 하지 않도록 수정
  - 중간단계 끊겼을시 마무리 정리 제대로 되도록
- [ ] task 결과 정리하고 디렉토리 지우는 OnExec 구현
필요없는 부분에서 struct 사용하지 말 것
- package의 함수 형태로 사용
일관되게 error 처리하기

----
Error code를 어떻게 다룰 것인가
어떻게 user에게 일관된 return을 보내줄 것인가
중단점이 발생했을때 해당 Task를 어떻게 관리할 것인가