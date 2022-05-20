# 模拟 Chain fountain 现象

运行后将生成一个 `config.json` 文件，解释以及默认值如下：

```
{
    "beadMass": 5, 铁球的重量 (g)
    "linkLength": 0.01, 链接的长度 (m)
    "initialHeight": 1, 开始的高度 (m)
    "timeStepSize": 0.001, 积分时间跨度 (s)
    "subSteps": 10, 积分子步骤数量
    "linkStiffness": 10000000, 链接刚度 (g/s^2)
    "gravity": 9.8, 重力 (m/s^2)
    "beakerWidth": 0.08, 烧杯宽度 (m)
    "beakerHeight": 0.02, 烧杯深度 (m)
    "beakerThickness": 0.05, 烧杯厚度 (m)
    "beakerStiffness": 10000000, 烧杯刚度 (g/s^2)
    "totalBeads": 300, 铁球总数
    "XOffset": 400, 初始X偏移
    "YOffset": 50, 初始Y偏移
    "zoom": 600, 缩放倍数
    "playSpeed": 2, 播放速度
}
```

运行时可用鼠标拖动界面。