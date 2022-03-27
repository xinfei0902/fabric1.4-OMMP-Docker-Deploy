#!/bin/bash
curPath=$PWD

## 构建go mod
go mod init

## 检测
go mod tidy 

## vendor
go bulid 