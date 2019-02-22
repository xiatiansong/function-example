# Readme

这个是miot function的runtime。用户只需要写func-example.go里的代码就可以了。这个工程主要帮助用户写、编译go。

一个基本的流程是：
1. 在func-example.go里修改代码。依赖的调用方法请查看miot.go
2. 对整体工程 go build . 检查编译是否通过
3. 粘贴func-example.go里的代码到miot function页面上传，验证

Note:
1. 请不要更新vendor里的依赖。（这个只是方便写代码用的，runtime不会上传这些依赖）
2. 如果有需要新的依赖，请提需求给我们。