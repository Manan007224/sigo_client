USE_VENDOR ?= no

ifeq ($(USE_VENDOR),no)
else
  MOD_FLAG := -mod=vendor
endif

test:
	@go test -v ${MOD_FLAG} \
	`go list ./pkg/...` 2>&1 -count=1