# go-facedetect
## v1.0版本（基于虹软）原仓库地址
[go-arcsoft项目地址](https://github.com/PonyWilliam/go-arcsoft)
## 版本
v2.0 通过与python进行socket通信获取人脸数据进行开发（2022年3月18日）  
v1.0 基于虹软人脸识别算法进行开发（old_version目录中）（2021年4月）  
## 关联项目
本项目人脸识别基于百度[PaddlePaddle](https://github.com/PaddlePaddle/Paddle)人工智能框架进行开发，详细请戳[python_facedetect](https://github.com/PonyWilliam/python_facedetect)  
人员信息展示基于electron进行开发，详细请戳[electron_facedetect](https://github.com/PonyWilliam/electron_facedetect)  
管理后端界面基于vue进行开发，详细请见[admin_facedetect](https://github.com/PonyWilliam/admin_facedetect)
管理后端人脸上传前校验基于paddle框架训练模型进行识别以及python Django框架对外暴露校验接口，详细请戳[django_facedetect](https://github.com/PonyWilliam/django_facedetect)
后端部署基于go-micro微服务框架进行开发，详细请见[micro_facedetect](https://github.com/PonyWilliam/micro_facedetect)
## 介绍
go-arcsoft通过移植opencv，虹软sdk，通过cgo链到c++库从而实现在golang上的人脸检测及追踪，同时封装[rfid串口](./RfidUtils/Rfid.go)在rfidutils内实现人脸追踪及上传云端
## 优点
相对于c++，golang拥有优秀的垃圾回收机制，对于内存泄漏有不可比拟的优势。在web方面有强大的库工具可以摆脱curl依赖快速进行http访问及构建http服务器，在串口方面可通过调用Windows等其它平台Api快速实现串口写入  
