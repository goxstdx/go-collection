# 默认值
TAG_NAME ?=
TAG_MESSAGE ?= ""

.DEFAULT_GOAL := help

.PHONY: add_tag
add_tag:
	@if [ -z "$(TAG_NAME)" ]; then \
		echo "错误: TAG_NAME为空，请提供有效的标签名称"; \
		exit 1; \
	fi
	@echo "使用标签名称: $(TAG_NAME)"
	@TAG_MSG="$(TAG_MESSAGE)"; \
	if [ -z "$$TAG_MSG" ]; then \
		TAG_MSG="$(TAG_NAME)"; \
	fi; \
	git tag $(TAG_NAME) -m "$$TAG_MSG" && \
	git push origin $(TAG_NAME)

.PHONY: push_pkg
push_pkg:
	@if [ -z "$(TAG_NAME)" ]; then \
		echo "错误: TAG_NAME为空，请提供有效的标签名称"; \
		exit 1; \
	fi

	@echo "执行 go test..."
	go test -v ./... || exit 1

	@echo "列出包..."
	go list -m all

	@echo "推送至 https://pkg.go.dev/..."
	go list -m github.com/ZHOUXING1997/collection@$(TAG_NAME)

.PHONY: help
help:
	@echo "可用命令列表:"
	@echo "  add_tag: 添加并推送 Git 标签. 使用: make add_tag TAG_NAME=v1.0.0 [TAG_MESSAGE=\"发布说明\"]"
	@echo "  push_pkg: 构建并列出包. 使用: make push_pkg TAG_NAME=v1.0.0"
