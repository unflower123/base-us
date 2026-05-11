## 端口使用说明

> 后面api服务端口用44xxx，rpc服务端口用55xxx 每个的端口按5递增
> api服务示例：api1 44001，api2 44005 依次递增
> rpc服务示例：rpc1 55001，rpc2 55005 依次增加

## 已经使用端口

### api服务

| 服务名称            | 服务                | 占用端口  |
|-----------------|-------------------|-------|
| manage admin    | admin-api         | 44000 |
| payin           | payin-api         | 44001 |
| payout          | payout-api        | 44005 |
| manage admin    | admin-api         | 44010 |
| merchant admin  | merchant-api      | 44015 |
| abnormal        | abnormal-api      | 44020 |
| seeelement      | seeelement-api    | 44025 |
| notify          | notify-api        | 44030 |
| callback        | callback-api      | 44035 |
| Scheduled Tasks | cron-api          | 44040 |
| bank callback   | bank callback-api | 44060 |

### rpc服务

| 服务名称        | 服务                | 占用端口  |
|-------------|-------------------|-------|
| Payin       | Payin-rpc         | 55001 |
| payout      | payout-rpc        | 55005 |
| manage admin | admin-rpc         | 55010 |
| merchant admin | merchant-rpc      | 55015 |
| abnormal    | abnormal-rpc      | 55020 |
| seeelement  | seeelement-rpc    | 55025 |
| notify      | notify-rpc        | 55030 |
| callback    | callback-rpc      | 55035 |
| Scheduled Tasks  | cron-rpc          | 55040 |
| bank callback   | bank callback-rpc | 55060 |