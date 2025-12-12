# 彩票中奖查询接口

## 接口说明

这是一个**无需登录验证**的公开接口，用于批量检查双色球或大乐透彩票是否中奖。

## 接口信息

- **接口地址**: `/lottery/check`
- **请求方式**: POST
- **是否需要登录**: ❌ 否（公开接口）
- **Content-Type**: application/json

## 请求参数

### 双色球示例

```json
{
  "lotteryType": "双色球",
  "winningSSQ": {
    "redBalls": [1, 5, 12, 18, 23, 30],
    "blueBall": 8
  },
  "purchasedSSQ": [
    {
      "redBalls": [1, 5, 12, 18, 23, 30],
      "blueBall": 8
    },
    {
      "redBalls": [2, 6, 13, 19, 24, 31],
      "blueBall": 9
    },
    {
      "redBalls": [1, 5, 12, 18, 23, 31],
      "blueBall": 8
    }
  ]
}
```

### 大乐透示例

```json
{
  "lotteryType": "大乐透",
  "winningDLT": {
    "frontZone": [5, 12, 18, 23, 30],
    "backZone": [3, 8]
  },
  "purchasedDLT": [
    {
      "frontZone": [5, 12, 18, 23, 30],
      "backZone": [3, 8]
    },
    {
      "frontZone": [1, 12, 18, 23, 30],
      "backZone": [3, 9]
    },
    {
      "frontZone": [5, 12, 18, 23, 31],
      "backZone": [3, 8]
    }
  ]
}
```

### 参数说明

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| lotteryType | string | 是 | 彩票类型，可选值：`双色球`、`大乐透` |
| winningSSQ | object | 条件 | 双色球中奖号码，当 lotteryType 为双色球时必填 |
| winningDLT | object | 条件 | 大乐透中奖号码，当 lotteryType 为大乐透时必填 |
| purchasedSSQ | array | 条件 | 购买的双色球号码（批量），当 lotteryType 为双色球时必填 |
| purchasedDLT | array | 条件 | 购买的大乐透号码（批量），当 lotteryType 为大乐透时必填 |

### 双色球号码格式

```json
{
  "redBalls": [1, 5, 12, 18, 23, 30],  // 红球：6个号码，范围 1-33，不能重复
  "blueBall": 8                         // 蓝球：1个号码，范围 1-16
}
```

### 大乐透号码格式

```json
{
  "frontZone": [5, 12, 18, 23, 30],  // 前区：5个号码，范围 1-35，不能重复
  "backZone": [3, 8]                  // 后区：2个号码，范围 1-12，不能重复
}
```

## 返回结果

### 成功响应示例（双色球）

```json
{
  "code": 0,
  "data": {
    "lotteryType": "双色球",
    "results": [
      {
        "numbers": {
          "redBalls": [1, 5, 12, 18, 23, 30],
          "blueBall": 8
        },
        "isWinning": true,
        "prize": "一等奖",
        "prizeLevel": 1,
        "matchDetail": "红球匹配6个，蓝球匹配"
      },
      {
        "numbers": {
          "redBalls": [2, 6, 13, 19, 24, 31],
          "blueBall": 9
        },
        "isWinning": false,
        "prize": "未中奖",
        "prizeLevel": 0,
        "matchDetail": "红球匹配0个，蓝球未匹配"
      },
      {
        "numbers": {
          "redBalls": [1, 5, 12, 18, 23, 31],
          "blueBall": 8
        },
        "isWinning": true,
        "prize": "三等奖",
        "prizeLevel": 3,
        "matchDetail": "红球匹配5个，蓝球匹配"
      }
    ]
  },
  "msg": "检查成功"
}
```

### 返回字段说明

| 字段名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| data | object | 返回数据 |
| data.lotteryType | string | 彩票类型 |
| data.results | array | 所有彩票的中奖结果 |
| data.results[].numbers | object | 该彩票的号码 |
| data.results[].isWinning | bool | 是否中奖 |
| data.results[].prize | string | 奖项名称 |
| data.results[].prizeLevel | int | 奖项等级（0-未中奖，1-一等奖，2-二等奖...） |
| data.results[].matchDetail | string | 匹配详情 |
| msg | string | 返回消息 |

## 中奖规则

### 双色球中奖规则

| 奖级 | 中奖条件 |
|------|---------|
| 一等奖 | 6个红球 + 1个蓝球 |
| 二等奖 | 6个红球 |
| 三等奖 | 5个红球 + 1个蓝球 |
| 四等奖 | 5个红球 或 4个红球 + 1个蓝球 |
| 五等奖 | 4个红球 或 3个红球 + 1个蓝球 |
| 六等奖 | 1个蓝球 |

### 大乐透中奖规则

| 奖级 | 中奖条件 |
|------|---------|
| 一等奖 | 前区5个 + 后区2个 |
| 二等奖 | 前区5个 + 后区1个 |
| 三等奖 | 前区5个 + 后区0个 |
| 四等奖 | 前区4个 + 后区2个 |
| 五等奖 | 前区4个 + 后区1个 |
| 六等奖 | 前区3个 + 后区2个 |
| 七等奖 | 前区4个 + 后区0个 |
| 八等奖 | 前区3个 + 后区1个 或 前区2个 + 后区2个 |
| 九等奖 | 前区3个 + 后区0个 或 前区1个 + 后区2个 或 前区2个 + 后区1个 或 前区0个 + 后区2个 |

## 使用示例（curl）

### 双色球查询

```bash
curl -X POST http://localhost:8888/lottery/check \
  -H "Content-Type: application/json" \
  -d '{
    "lotteryType": "双色球",
    "winningSSQ": {
      "redBalls": [1, 5, 12, 18, 23, 30],
      "blueBall": 8
    },
    "purchasedSSQ": [
      {
        "redBalls": [1, 5, 12, 18, 23, 30],
        "blueBall": 8
      },
      {
        "redBalls": [1, 5, 12, 18, 23, 31],
        "blueBall": 8
      }
    ]
  }'
```

### 大乐透查询

```bash
curl -X POST http://localhost:8888/lottery/check \
  -H "Content-Type: application/json" \
  -d '{
    "lotteryType": "大乐透",
    "winningDLT": {
      "frontZone": [5, 12, 18, 23, 30],
      "backZone": [3, 8]
    },
    "purchasedDLT": [
      {
        "frontZone": [5, 12, 18, 23, 30],
        "backZone": [3, 8]
      },
      {
        "frontZone": [1, 12, 18, 23, 30],
        "backZone": [3, 8]
      }
    ]
  }'
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 7 | 参数错误或验证失败 |

## 注意事项

1. ✅ 本接口**无需登录验证**，可直接调用
2. ✅ 支持批量查询多个彩票号码
3. ✅ 一次请求只能查询一种彩票类型（双色球或大乐透）
4. ✅ 号码范围和数量会自动验证
5. ✅ 重复号码会被检测并返回错误

