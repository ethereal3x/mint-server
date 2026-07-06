# Jenkins CI 配置说明

Jenkins 地址：https://jenkins.l3xx.cc  
镜像仓库：`register.l3xx.cc/mint-server:<git-short-sha>`

流水线定义见仓库根目录 [Jenkinsfile](../../Jenkinsfile)。

## 流水线阶段

| 阶段 | 触发条件 | 说明 |
|------|----------|------|
| Checkout | 所有分支 | 拉代码，计算 `GIT_SHORT_SHA` |
| Build & Push | 仅 `master` | 构建并推送镜像到 `register.l3xx.cc`（编译在 Dockerfile 内完成） |

提交信息含 `[skip ci]` 时跳过 Build & Push。

## 部署（手动）

Jenkins **不会** push GitHub 修改 manifest。镜像推送成功后，手动更新 K8s 部署：

1. 在 Jenkins 构建日志中确认镜像 tag（如 `e4faac1`）
2. 修改 `deploy/k8s/mint-server.yaml` 中的 image：

   ```yaml
   image: register.l3xx.cc/mint-server:e4faac1
   ```

3. commit 并 push，由 ArgoCD 同步到集群

后续可接入 ArgoCD Image Updater 自动跟踪 registry 新 tag。

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

## 2. 创建 Multibranch Pipeline

1. **New Item** → **Multibranch Pipeline**
2. **Branch Sources** → **GitHub**
   - 选择已连接的 GitHub 账号 / App
   - Repository：`ethereal3x/mint-server`
   - Behaviours：Discover branches（PR 可选，建议 Filter 仅 `master`）
3. **Build Configuration** → Mode: **by Jenkinsfile**
4. **Script Path**：`Jenkinsfile`
5. 保存并 **Scan Repository Now**

## 3. Webhook（GitHub 侧）

GitHub Repo → **Settings → Webhooks**，确认 Jenkins 已注册 push 事件。

## 4. 与 GitHub Actions 的关系

两套 CI 可并存。若只保留 Jenkins，可在 `.github/workflows/build.yml` 中禁用 push 触发或删除 workflow。

## 5. 手动触发

Jenkins 任务页 → **Build Now**（master 分支构建并推送镜像）。

## 6. 常见问题

**卡在 Waiting for agent**  
检查 Jenkins agent 是否在线，Pipeline 使用 `agent any`。

**docker login 失败**  
检查 `REGISTRY_CREDENTIALS` 用户名密码及 registry 网络可达性。

**docker: permission denied**  
将 jenkins 用户加入 docker 组后重启 Jenkins：

```bash
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
```

验证：`sudo -u jenkins docker ps`
