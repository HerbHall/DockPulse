IMAGE ?= herbhall/dockpulse
TAG ?= latest

BUILDER=buildx-multi-arch

.PHONY: build-extension install-extension update-extension clean

build-extension:
	docker build --tag=$(IMAGE):$(TAG) .

install-extension:
	docker extension install $(IMAGE):$(TAG)

update-extension:
	docker extension update $(IMAGE):$(TAG)

debug-extension:
	docker extension dev debug $(IMAGE)

reset-extension:
	docker extension dev reset $(IMAGE)

push-extension:
	docker buildx create --name=$(BUILDER) || true
	docker buildx use $(BUILDER)
	docker buildx build --push --platform=linux/amd64,linux/arm64 --tag=$(IMAGE):$(TAG) .

clean:
	docker extension rm $(IMAGE) || true

# Go targets
.PHONY: go-build go-test go-lint

go-build:
	cd backend && go build ./...

go-test:
	cd backend && go test ./...

go-lint:
	cd backend && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...

# Frontend targets
.PHONY: fe-build fe-test fe-lint fe-typecheck

fe-build:
	cd ui && npm run build

fe-test:
	cd ui && npx vitest run

fe-lint:
	cd ui && npx eslint src/

fe-typecheck:
	cd ui && npx tsc --noEmit

# All checks
.PHONY: validate
validate: go-build go-test go-lint fe-typecheck fe-lint fe-test
