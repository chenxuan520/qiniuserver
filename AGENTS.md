# AGENTS.md — qiniuserver 架构速览

这个仓库是一个自托管的 **Web 前端 + HTTP API**，用于把文件上传到 **七牛云 Kodo** 并进行简单管理。
项目以单个 Go 二进制方式发布（Gin），同时提供：

- 静态前端：位于 `assert/`（UI 主页面是 `assert/index.html`）。
- 后端 API：位于 `/api/*`，把上传/列举/删除等操作代理到七牛云。

未发现 Cursor/Copilot/Trae 的规则文件（未找到 `.cursor/rules/`、`.github/copilot-instructions.md`、`.trae/rules/`）。

## 1) 项目概览

### 运行形态

- **单进程**：`src/main.go` 启动服务，加载配置，初始化七牛 SDK 客户端，然后启动 Gin。
- **前端托管**：Gin 将 `GET /` 指向 `./assert/index.html`，并将 `GET /static/*` 指向 `./assert/`。
- **API**：`/api` 路由组会挂载一个可选的密码认证中间件，然后暴露文件相关接口。

### 主要代码分层（按职责）

- `src/main.go`：入口；组装配置、七牛初始化、路由、静态资源托管。
- `src/config/config.go`：读取 JSON 配置，写入 `config.GlobalConfig`。
- `src/middlerware/auth.go`：API 密码校验（读取 `Authorization` 请求头）。
- `src/controller/file.go`：HTTP handler（上传、URL 上传、列表、删除、info）。
- `src/utils/qiniu.go`：七牛 SDK 封装（init/upload/list/delete）。
- `src/controller/response/res.go`：统一响应结构封装。

### 典型请求链路（上传文件）

1. 浏览器访问 `/`（加载静态 HTML + JS）。
2. 前端 JS 请求 `GET /api/info` 获取 `domain` 和 `upload_path`。
3. 选择/拖拽/粘贴文件后，前端向 `POST /api/upload` 发送 multipart。
4. 后端先把文件保存到本地临时文件，再通过七牛 SDK 上传，随后删除本地临时文件，并返回最终 URL。

### 关键依赖（见 `src/go.mod`）

- Web 框架：`github.com/gin-gonic/gin`
- JSON：`github.com/json-iterator/go`（用于解析配置）
- 七牛云：`github.com/qiniu/go-sdk/v7`

## 2) 构建与命令

### 本地构建/运行

- 构建（推荐；与发布流程一致）：`./build.sh`
- 构建（直接）：`cd src && go build .`
- 从仓库根目录运行（要求当前目录下存在 `./assert/` 和 `./config/`）：`./qiniuserver`

启动后会打印绑定地址（见 `src/main.go`）。

### 测试

- 运行全部测试：`cd src && go test ./...`

注意：`src/utils/qiniu_test.go` 中包含偏“集成测试”的用例（例如 `TestUploadPath`），它会尝试上传一个硬编码的本地文件路径；在干净环境下通常会失败，除非你调整路径并提供可用的七牛配置。

### 发布/打包

- CI：`.github/workflows/build_and_release.yml`
  - Linux/macOS：调用 `./build.sh` 构建。
  - Windows：`cd ./src && go build .` 构建。
  - 打包产物包含二进制 + `./assert` + `./config` + `LICENSE` 的 tar.gz。

## 3) 代码风格

### Go 代码习惯（本仓库现状）

- 使用 `gofmt` 格式化（未发现其它格式化/静态检查配置）。
- 包名遵循小写风格；其中 `middlerware` 是仓库内既有命名（拼写与常见的 `middleware` 不同），属于现有 public surface，改动需谨慎。

### API / JSON 约定

- 配置文件 JSON 使用 `snake_case` 字段（见 `src/config/config.go` 的 struct tag）。
- API 统一响应结构为 `{ "code": number, "data": any }`（见 `src/controller/response/res.go`）。

### 本仓库建议遵循的写法

- handler 逻辑主要放 `controller/`；七牛相关操作通过 `utils/` 统一封装。
- 上传 key/文件名生成应通过 `utils.CreateFileName(...)`，以保持与现有配置语义一致。

## 4) 测试说明

### 现有测试

- `src/utils/string_test.go`：文件名模板逻辑相关的简单测试。
- `src/utils/qiniu_test.go`：需要读取 `../../config/config.json`，依赖网络与七牛有效凭证，并且可能依赖特定本地文件路径。

### 推荐执行方式

- 仅验证纯逻辑时，优先跑不依赖外部环境的用例，例如：`cd src && go test ./utils -run TestCreateFileName`
- 七牛相关测试更像集成测试：建议在专用账号/桶、非生产凭证、可控网络环境下执行，并先调整硬编码的文件路径。

## 5) 安全注意事项

### 认证模型

- 若 `host.password` 为空，则 **`/api/*` 全部不需要认证**（见 `src/middlerware/auth.go`）。
- 若设置了 `host.password`，客户端需要在请求头 `Authorization` 中直接携带该密码值。
- 代码里只有一个非常简单的“失败后延迟”机制（全进程级别的 3 秒节流），不是按 IP/用户维度。

### 密钥与敏感信息

- 七牛 AK/SK 等凭证来自 JSON 配置（见 `src/config/config.go`）。
- `config/demo.json` 是示例模板。
- `config/config.json` 可能包含真实凭证，应按敏感文件处理：避免把内容贴到日志/issue 中，也尽量避免进入 git 历史。

### 与功能直接相关的风险点

- `POST /api/upload_url` 会在服务端下载任意 URL 的内容并上传到七牛（见 `src/controller/file.go`）。对不可信用户开放该接口时需谨慎。

## 6) 配置

### 配置文件路径

启动时按如下顺序读取（见 `src/config/config.go`）：

1. `./config/config.json`（相对当前工作目录）
2. `/config/config.json`（绝对路径，适合容器挂载）

### 配置字段说明

- `host.ip` / `host.port`：Gin 服务监听地址。
- `host.password`：为空则关闭 API 密码校验。
- `qiniu.domain`：用于拼接返回给前端的文件 URL 的域名。
- `qiniu.access_key` / `qiniu.secret_key` / `qiniu.bucket`：七牛凭证与 bucket。
- `qiniu.upload_path`：对象 key 的前缀路径。
- `qiniu.zone`：只能是 `Huanan` / `Huabei` / `Huadong` / `Xingjiapo`（在 `src/utils/qiniu.go` 校验）。
- `qiniu.file_name`：文件名模板（由 `utils.CreateFileName` 解析）：
  - `%d`：时间戳（格式 `20060112150405`）
  - `%f`：原始文件名（空格会替换为 `_`）
  - `%r`：6 位随机字符串
