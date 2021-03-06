# 模拟 Chain fountain 现象

## 介绍
当烧杯中的链条被拉到烧杯的侧面时，链条会像预期的那样开始从烧杯中抽出。
但是如果条件合适，链条会跃出烧杯，在开始落地之前达到一个更高的高度，形成一个喷泉形状。
你可以在[Mould的视频](https://www.youtube.com/watch?v=_dQJBBklpQQ)中清楚地看到这一点。
令人惊讶的是，在Mould的2013年视频之前，这种效果相对不为人知，因此`Mould Effect`是一个亲切的、替代性的名字。
从那时起，有关`Mould Effect`的相当多的论文已经发表。

链条喷泉很有趣，因为它的解释很棘手。
[Biggins和Warner](http://rspa.royalsocietypublishing.org/content/470/2163/20130689)表明，在一个直观的模型中，静止的链条被简单地拉起来做虹吸运动，没有形成链条喷泉。
你可以按照Biggins和Warner的分析自己去复现一下，[Isaac Physics](https://isaacphysics.org/questions/chain_fountain)提供了一个教程并转载了这一分析。
在这个理想化的模型中，链式拉起的相互作用被假定发生在一个封闭的系统中，链和烧杯之间没有相互作用，因此它是一个[完全非弹性的碰撞](https://en.wikipedia.org/w/index.php?title=Inelastic_collision&oldid=806063088)，会耗散大量的能量。
为了产生链条喷泉，Biggins和Warner预测烧杯必须提供一个额外的非正常力。他们通过假设链条不能超过某种最大可能的曲率，为这种力提供了一种可能的机制。
当我最初看到这个解释时，它似乎是不可理喻的。

一个链条喷泉模拟器将让我们能够运行许多数值实验，并进一步探索这种动态物理学。
模拟器的第一个目标是让我们找到一个能产生连锁喷泉的最小模型。
最小模型中的烧杯-链条相互作用的性质将帮助我们理解这一个反常的理论。
实际上，Biggins & Warner提供了他们的数字实验的简要描述，该实验确实使用椭圆的烧杯通过链条相互作用产生了一个喷泉。
本模拟器使用的方法是基于这一描述的（见方法），并且可以重现他们的结果。

次要的目标是重现对链条喷泉的预测和观察，也许还可以做出一些我们自己的预测。
例如，预测和观察到的链条喷泉的高度通常为初始高度的1.2倍。
此外，链条喷泉将以∝t<sup>2</sup>的速度增长。如果数值误差可以降到足够低，我希望能观察到这两点。
在模拟器中，似乎烧杯的宽度对链状喷泉的高度有很大影响。试着阅读一些已发表研究，希望能够找到更多可以测试的特性。

考虑到`Mould Effect`的可及性，模拟器的另一个目标是使其本身具有可及性。
作为一个Go写的模拟器，这将使任何人都能对链条喷泉进行实验，不管他们是否熟悉数字或是否有能力获得许多长链和烧杯。
我们很幸运，现在生活在一个可执行文件可以互动、执行并在几乎所有东西上运行的世界里，所以我想在视觉、科学、数值实验方面进行探索。
这种平台的选择确实给科学计算带来了一些挑战，但至少在这种情况下，它还没有成为一个限制性因素。

## 方法
Biggins & Warner 对他们的数字实验的描述对下面描述的方法的设计有很大帮助。
### 模型
链条的模型是使用一连串的点质量，其位置为 *x&#818;<sub>i</sub>* ，质量为 *m* ，由弹簧常数 *k* 和长度 *&delta;l* 的刚性弹簧连接。
然后可以用牛顿力学来模拟力学，其中 *F<sub>i</sub>=x&#776;&#818;<sub>i</sub>* 。
作用在每个珠子上的力通常包括联合的弹簧力 *F<sub>i</sub><sup>(L)</sup>=k(x&#818;<sub>i+1</sub>-x&#818;<sub>i</sub>)+k(x&#818;<sub>i</sub>-x&#818;<sub>i-1</sub>)* 和重力 *F<sub>i</sub><sup>(G)</sup>=g/m* ，其中 *g* 是重力加速度。
如果弹簧足够坚硬，这将近似于一个不可伸展的链条。在保持质量密度 *&lambda;=m/&delta;l* 不变的情况下，在 *&delta;l&rightarrow;0* 的极限下建立绳子的模型。链条有时也会与烧杯和地板相互作用。

烧杯与链条的相互作用被模拟为一个弹性恢复力 *F<sup>(B)</sup>=k<sub>B</sub>x<sub>&perp;</sub>* ，其中 *k<sub>B</sub>* 是烧杯的刚度， *x<sub>&perp;</sub>* 是变形距离，即 *x&#818;<sub>i</sub>* 点到烧杯外任何一点的最小距离。
这个力是保守的，所以一旦离开烧杯，链条的所有能量都会被返回。

由于链条喷泉的高度通常与烧杯的高度成正比，地板必须是重要的。地板的模型是将点状质量带到一个立即停止的地方，这样，这种相互作用是完全无弹性的。
地板可以被看作是一种斯托克定律无限粘性的液体。
为了让链条在地板上变平，只有Y方向的运动被这样处理，而X方向的运动被允许为正常运动。这仍然具有预期的效果。

为了使链条进入运动状态，初始位置的选择要使链条悬挂在烧杯的边缘。
没有给链条提供初始速度，但它将通过虹吸运动开始下降。
这种选择很简单，有助于积分的稳定性，也有助于保持积分方法的简单。

这个模型面临的挑战是如何获得高的弹簧刚度*k*。随着*k*的增加，积分时间步长必须变小，以考虑到小位移产生的大加速度。
在拉格朗日力学中可以找到一些避免或帮助解决这一困难的可能的替代方案，它可以明确地考虑到联结约束。
就目前而言，使用目前的积分方法可以实现足够高的*k*。
### 积分
为了使点质量的位置随时间变化而积分，我们使用了基本的韦尔莱积分法。
韦尔莱积分法专门对牛顿运动方程 *x&#776;=F(x)/m*  进行积分，并保持系统总能量守恒。
选择一个恒定的积分时间步长为 *&Delta;t* 使得能够在时间为 *t<sub>i</sub>=&Delta;t&times;i* 时计算位置。
韦尔莱积分法的计算成本非常小，因此 *&Delta;t* 可以选择得非常小，即使是 *&approx;10&micro;s* 也可以轻易运行。
对于初始积分步骤，我们假设没有初始速度计算 *x&#818;<sub>1</sub>=x&#818;<sub>0</sub>+A(x&#818;<sub>0</sub>)&Delta;t<sup>2</sup>/2* 。
然后对于下面的积分步骤，我们计算 *x&#818;<sub>i+1</sub>=2x&#818;<sub>i</sub>-x&#818;<sub>i-1</sub>+A(x&#818;<sub>i</sub>)&Delta;t<sup>2</sup>* 并带入上两次的计算。
函数 *A(x&#818;<sub>i</sub>)* 是第 *i* 个点质量的加速度，使得 *A(x&#818;<sub>i</sub>) = [F<sub>i</sub><sup>(L)</sup>(x&#818;<sub>i</sub>)+F<sub>i</sub><sup>(G)</sup>(x&#818;<sub>i</sub>)+F<sub>i</sub><sup>(B)</sup>(x&#818;<sub>i</sub>)]/m* 。

但是请注意，地板的相互作用不能由这个积分自然地解释。为了纳入这种相互作用，在每个积分步骤之后，检查点质量是否在地板上 *y<sub>i</sub><=0* 或曾经在地板上 *y<sub>i-1</sub><=0* 。
如果是这样，将点质量的 *y* 位置改为 *y=0* 。韦尔莱积分法不会介意这种变化，它将继续进行，就像点质量的能量从来没有任何速度一样。

### 实现细节
虽然对于现代CPU来说，运行韦尔莱积分法是比较容易的，但是将结果存储到内存中则要昂贵得多。
为了允许10µs（100,000/s）数量级的时间步骤，可以指定一些子时间步骤。
对于每个时间步长，积分方法要运行多次，有效地将该时间步长分成许多更小的时间步长。
只有在所有的子时间步骤完成后，结果才会被写入内存。当模拟器播放积分结果时，它将只显示存储在内存中的结果，而不是由子步骤计算的中间结果。

## 运行配置

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