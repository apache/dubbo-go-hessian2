# UUID-info

1. explain dubbo-go-hessian2 strut UUID
- JavaServer -> create UUID -> GO Client -> UUID struct (PASS)
- dubbo-go-hessian2 cannot create UUID strut
- see jdk source code of class:[java.util.UUID] learning how to create UUID struct
- see https://github.com/satori/go.uuid
