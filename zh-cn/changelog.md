# 更新日志

## _2.0.1 开发中_
## 待做项
- [ ] 切换全部golang正则表达式为标准表达式模式
- [ ] 重写前端界面
- [ ] 系统发布组增加市场,官方发布组会自动同步,也可自定义添加和采用
- [ ] 增加搜索关键字数据类型
- [ ] 备份导入
- [ ] 识别不到语言之后 使用媒体发布国家的语言
- [ ] 新手引导 
- [ ] 从Sonarr导入下载器 

## 2.0.0.6
- [x] 兼容站点中非中文标题在 [连载] 前导致的匹配失败

## 2.0.0.5
- [x] 优化版本发布相关

## 2.0.0.4
- [x] 增加自动备份功能
- [x] 增加自动计算连续集数/整季文件大小
- [x] 增加导出配置
- [x] 增加数据源状态开关
- [x] 优化种子解析能力

## 2.0.0.3
- [x] 优化兼容标题匹配结果
- [x] 修复定时刷新bug

## 2.0.0.2
- [x] 增加rss自定义标题输出格式
  - 动漫默认格式:
    - `[{releaseGroup}][{chineseTitle}] {title} - {season}{episode} ({abEpisode}) [{language}][{quality}][{video}][{audio}][{mediaType}]`
  - TV电视剧默认格式:
    - `[{releaseGroup}][{chineseTitle}] {title} - {season}{episode} [{language}][{quality}][{video}][{audio}][{mediaType}]`
  - 使用方式查看 [系统配置](/zh-cn/system.md)
- [ ] 增加加入自定义发布组(必须全词匹配)


## 2.0.0.1 
- [x] 增加若干种子识别
- [x] 修复若干bug
 
## 2.0.0 (2022-08-19)
- [x] 优化页面搜索页码错误问题
- [x] 修复页面数据源搜索后总数不正确问题
- [x] 优化自动分析标题时机为录入数据源时进行分析
- [x] 修改全自动匹配下的 包含标题 为 正则表达式
- [x] 优化测试匹配的功能速度
- [x] 增加代理 数据源 类型 支持 prowlarr jacket
  - 代理jacket模式
    - 分组设置
       - 分组中设置代理jacket模式
       - 分钟中绑定jakcet数据源 暂定支持一个 多个不知道会不会卡
    - 分组媒体配置
      - 不可设置单独数据源 直接使用分组中绑定的jacket数据源
    - 搜索流程
      - sonarr 调用 toznab api 进行搜索
      - 获取搜索词
      - 转发搜索jacket 
      - 获取jacket返回结果
      - 解析返回数据内容
      - 进行匹配`分组媒体`下配置的媒体数据
      - 匹配季集信息
      - 格式化输出sonarr认识的数据项
- [x] 修复torznab rss 同步问题
- [x] 增加优化识别率

## 1.5.6 (2022-08-09) 
- 增加 系统设置-全局代理设置
- 修改 数据源代理设置为开关模式
- 修改 `tmdb` 代理选择方式
- 优化 `Sonarr数据源` 展示效果
- 新增 `全新匹配组件` 识别精度更高 更准确
- 优化无pubDate设置默认为添加时
- 增加内置站点
  - 蜜柑计划
  - 萌番组
  - ACG.RIP
  - 動漫花園
  - 爱恋动漫
  - 漫猫动漫
  - 末日动漫
  - eztv
  - RARBG
- 常规修改
  - 增加nyaa网站支持
  - 优化测试匹配不使用可查看的列表匹配内容 改为直接使用数据源
  - 修复手动正则模式bug
  - 修复手动指定季导致第几集为0的问题
  - 修复标题中涵盖第几季 但是绝对集数 输出错误问题

## 1.5.5 (2022-07-26)
- 增加数据源缓存天数功能
- 兼容Sonarr集数断连的解析能力
  - 但是Sonarr必须包含标题中的集数 如 Sonarr中的第三季里面的12集 必须存在 否则标题 xx S03E12 将会匹配不成功
- 修改最低别名匹配字数5个字以上
- 常规修改
  - 优化缓存列表为倒序排序
  - 修复缓存天数bug
  - 解决重复下载重命名失败问题


## 1.5.4 (2022-07-20)
- 修复生成匹配规则问题


## 1.5.3 (2022-07-08)
- 修改数据库内核优化兼容性
- 增加自定义端口号
- 增加TheMoviedb 自定义代理/HOST
- 优化分组媒体列表页面中的展示文字
- 优化服务配置中的必填项
- 优化http请求内核
- 增加筛选发布组
- 修复全自动模式下填入包含规则bug问题
- 小版本更新
  - 增加Nc-raws baha 为繁体

## 1.5.2 (2022-06-30)
- 优化公告弹窗
- 优化种子下载监听任务
- 增加Transmission
- 增加group.xml?refresh=1 即时刷新分组数据
- 优化测试匹配规则时处理语言和质量
- 优化下载器监听日志输出

## 1.5.1 (2022-06-15)
- 增加集数 [05v2] 此类型的兼容
- 增加group.xml文件中的标题
- 增加group.xml增加使用了那个分组使用了那个数据源
- 修复首页今日下载量统计
- 优化保存监听任务成功过多次数
- 已删除种子任务自动设置为已完成
- 优化繁体big5语种为Chinese Traditional
- 优化qb下载进度100% 设置监控任务bug问题
- 优化数据源全部问题
- 修改新的授权接口

## v1.4 (2022-05-17) 

- **增加种子v2 哈希兼容**
- Torznab 免费开放
- 一键添加Rss/Torznab 免费开放
- 增加一键Torznab
- 优化将数据源名称从4个限制为3个
- 优化同步删除媒体数据
- 优化qBittorrent保存配置检测主机地址
- 优化Torznab搜索集搜索不到问题
- 优化首页永久会员显示


## v1.2-v1.3(2022-05-10)

- 增加自动化Sonarr媒体导入
- 增加分组媒体规则自动识别
- 增加Torznab数据格式提供
- 优化rss xml 里面的数据结构
- 增加服务设置-Sonarr-刷新时间间隔自定义
- 增加Tmdb查看其它媒体标题
- 优化Tmdb数据媒体同步速度
- 修改请求支持gzip 压缩xml大小
- 修改数据源超时时间为60s


- 当前版本维护更新记录
  - 2022-05-13 
    - 修复1.1升级到1.3的配置兼容
    - 修复测试菜单还原
    - 修复Season为0的问题
  - 2022-05-11 
    - 优化前端分组展示
    - 增加界面右上角增加教程按钮
    - 修复含有SP季度 自动计算季信息错误问题
    - 修复Torznab搜索功能

## v1.1(2022-04-29)

- 修复数据长度为0的bug