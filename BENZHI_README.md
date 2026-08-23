# cfar-ca — 单元平均 CA-CFAR 检测核算命令行工具

输入距离向幅度序列、保护单元数、参考窗长与名义虚警率 Pfa，逐单元输出阈值（α×参考均值）、检出标志与经验虚警率。

## 构建 / 运行 / 测试

```text
go build ./...
go run . detect example/target-in-noise.json
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
