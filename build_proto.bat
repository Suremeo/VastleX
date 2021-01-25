@echo off
protoc --go_out=./ ./vastlex/plugin/actions/protobuf/*.proto