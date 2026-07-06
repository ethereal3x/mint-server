# Jenkins CI 配置说明

Jenkins 地址：https://jenkins.l3xx.cc  
镜像仓库：`register.l3xx.cc/mint-server:<git-short-sha>`

流水线定义见仓库根目录 [Jenkinsfile](../../Jenkinsfile)，与 GitHub Actions `.github/workflows/build.yml` 行为对齐。

## 流水线阶段

| 阶段 | 触发条件 | 说明 |
|------|----------|------|
| Checkout | 所有分支 | 拉代码，计算 `GIT_SHORT_SHA` |
| Build & Push | 仅 `master` | 构建并推送镜像到 `register.l3xx.cc`（编译在 Dockerfile 内完成） |
| Update Manifest | 仅 `master` | 更新 `deploy/k8s/mint-server.yaml` 镜像 tag 并 push |

提交信息含 `[skip ci]` 时跳过 Build & Push / Update Manifest，避免 manifest 回写触发循环构建。

## Jenkins Agent 前置依赖

Agent 节点需安装：

- Docker，且 **jenkins 用户有权限访问** `/var/run/docker.sock`（见下方权限说明）
- Git

Go 编译在 `docker build` 阶段由 `Dockerfile` 内的 `golang:1.25-alpine` 完成，agent 无需单独安装 Go。

## 1. 创建 Credentials

路径：**Manage Jenkins → Credentials → System → Global credentials**

| ID | 类型 | 用途 |
|----|------|------|
| `REGISTRY_CREDENTIALS` | Username with password | 登录 `register.l3xx.cc` |
| `GIT_PUSH_CREDENTIALS` | Username with password | push manifest 到 GitHub（Username 填 GitHub 用户名，Password 填 PAT） |

PAT 需具备 `repo` 写权限。

## 2. 创建 Multibranch Pipeline

1. **New Item** → **Multibranch Pipeline**
2. **Branch Sources** → **GitHub**
   - 选择已连接的 GitHub 账号 / App
   - Repository：`ethereal3x/mint-server`
   - Behaviours：Discover branches（PR 可选）
3. **Build Configuration** → Mode: **by Jenkinsfile**
4. **Script Path**：`Jenkinsfile`
5. 保存并 **Scan Repository Now**

## 3. Webhook（GitHub 侧）

GitHub Repo → **Settings → Webhooks**，确认 Jenkins 已注册 push 事件；Multibranch 扫描策略建议：

- 推送时自动扫描分支
- 或定时 scan（如每 5 分钟）

## 4. 与 GitHub Actions 的关系

两套 CI 可并存：

- **GitHub Actions**：依赖 self-hosted runner（`runs-on: self-hosted`）
- **Jenkins**：依赖 jenkins.l3xx.cc agent

若只保留 Jenkins，可在 `.github/workflows/build.yml` 中禁用 push 触发或删除 workflow。

## 5. 手动触发

Jenkins 任务页 → **Build Now**（master 分支会完整走 verify + 镜像推送 + manifest 更新）。

## 6. 常见问题

**卡在 Waiting for agent**  
检查 Jenkins agent 是否在线，Pipeline 使用 `agent any`。

**docker login 失败**  
检查 `REGISTRY_CREDENTIALS` 用户名密码及 registry 网络可达性。

**git push 403**  
检查 PAT 权限，或 Jenkins 使用的 GitHub 账号是否有 master 写权限。

**docker: permission denied**  
将 jenkins 用户加入 docker 组后重启 Jenkins：

```bash
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
```

验证：`sudo -u jenkins docker ps`
