# cfar-ca

单元平均 CA-CFAR（Cell-Averaging Constant False Alarm Rate）检测核算工具。输入距离向幅度序列、保护单元数、每侧参考单元数与名义虚警率 Pfa，输出每个被测单元（CUT）的检测阈值、检出标志、相对阈值的裕量，以及整段序列的经验虚警率。这是雷达检测内核：阈值严格取 `α × 参考均值`，放大系数按教科书公式 `α = N(Pfa^{−1/N} − 1)`（N 为参考单元总数）由 Pfa 与参考单元数共同决定，均匀噪声长序列下的检测比例会落在 Pfa 的合理波动带内。不做安防告警、不做阵列栅瓣，也不把 CFAR 阈值当成链路 SNR 余量。

## 用法

```bash
go build ./...
go run . detect example/target-in-noise.json
```

`detect` 打印逐单元表格（索引、幅度、阈值、检出、裕量、状态），末尾给出检出索引、目标（局部峰值）与经验虚警。算例 JSON 形如：

```json
{"amplitudes": [...], "guards": 2, "refs": 8, "pfa": 0.001}
```

`example/` 内置两个算例：

- `target-in-noise.json` —— 中间一个强目标（幅度 28），两侧指数噪声；应只在目标处检出。
- `uniform-noise.json` —— 2000 点均匀指数噪声；经验虚警应接近名义 Pfa。

其他子命令：

```bash
go run . detect -json example/target-in-noise.json   # JSON 导出（无效单元为 null）
go run . sweep -pfa "1e-2 1e-3 1e-4" example/target-in-noise.json  # Pfa 扫描对照
go run . compare -high 1e-3 -low 1e-4 example/target-in-noise.json # 交叉规则对照
go run . alpha -pfa 1e-3 -refs 8        # 打印 α = N(Pfa^{−1/N}−1)
go run . synth -n 4000 -guards 2 -refs 16 -pfa 1e-3 -seed 7  # 均匀噪声经验虚警带校验
go run . version
```

文件名传 `-` 时从 stdin 读取。

### HTTP 服务

`cfar-ca -http :8080` 启动核算 API，`/api/detect` 的请求体为 `{"spec": <算例 JSON>}`，其余接口用 `{"spec": ..., ...}` 结构：

- `GET /health` 健康检查；
- `POST /api/detect` 逐单元阈值/检出/裕量与经验虚警；
- `POST /api/sweep` 对多个 Pfa 扫描 α 与检出数；
- `POST /api/compare` 两个 Pfa 的交叉规则对照；
- `POST /api/alpha` 直接返回放大系数 α。

## 关键约定

- **阈值 = α × 参考均值**。参考单元取 CUT 两侧保护单元外侧各 Refs 个，CUT 自身与保护单元一律不进参考均值——强目标不会抬高自己的阈值。
- **α 由 Pfa 与参考单元数共同决定**：`α = N(Pfa^{−1/N} − 1)`，N = 2×Refs。只把 Pfa 降一个数量级，α 升高、检出变少；只把目标幅度加大，该 CUT 相对阈值的裕量升高。
- **边缘 CUT 参考不足 → 标记无效**，不给阈值、不参与检出，不借用半窗、不补零。
- **非法输入一律报错**（stderr + 非零退出）：Pfa∉(0,1)、负幅度、参考窗伸出序列、保护单元≥参考单元、参考窗长 0、序列为空、文件不存在、JSON 不合法。

## 构建与测试

```bash
go build ./...
go test ./...
```

纯标准库实现，无第三方依赖。

## 能力边界

- 针对单段距离向序列做单元平均检测；不处理多普勒/脉冲维，不做恒虚警的分布自适应（如 OS-CFAR）。
- 幅度序列按平方律检波后的指数背景处理（α 公式即该背景的教科书形式）；瑞利包络场景的换算辅助函数在 `internal/stats` 中提供，但检测主路径按指数背景求解。
- 不收敛类迭代不存在于本内核；经验虚警带基于二项分布近似（±kσ）给出。
