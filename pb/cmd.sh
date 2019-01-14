#!/usr/bin/bash
protoc service.proto --go_out=plugins=grpc:.