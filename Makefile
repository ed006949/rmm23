include Makefile.local

DATE		:=	$(shell date +%s)
GIT_STATUS	:=	$(shell git status --short)
GIT_COMMIT	:=	$(shell git rev-parse --short HEAD)
GO_VERSION	:=	$(shell go version | awk '{print $$3}' | sed 's/^go//')

LDFLAGS		:=	-s -w -X '$(NAME)/src/l.buildName=$(NAME)' -X '$(NAME)/src/l.buildTime=$(DATE)' -X '$(NAME)/src/l.buildCommit=$(GIT_COMMIT)'
GOBUILD		:=	go build -trimpath

.PHONY: all build clean commit diff execute fix init install lint normalize race release run status test update upgrade vet gitignore init_hook init_localpackage clean-clean-clean-clean init-init-init-init

all: commit race build

build:
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o ./bin/$(NAME) ./src/*.go
	GOOS=freebsd GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o ./bin/$(NAME)-freebsd-amd64 ./src/*.go
	$(GOBUILD) -gcflags=all="-N -l" -ldflags="$(LDFLAGS)" -o ./bin/$(NAME)-delve ./src/*.go

clean:
	-go clean -i -r -x -cache -testcache -modcache -fuzzcache
	-rm -v go.mod go.sum
	-find . -name ".DS_Store" -delete
	-find . -name "._.DS_Store" -delete
	-go mod init $(TARGET)
	-go get -u ./...
	-go mod tidy

commit: status normalize
ifneq ($(GIT_STATUS),)
	git add .
	git commit --no-edit
	git push
endif

diff:
	git diff

execute:
	./bin/$(NAME) $(COMMAND_LINE)

fix:
	go fix ./...

init:
	go mod init $(TARGET)
	go get -u ./...
	go mod tidy

install:
	@echo $(NAME) $(PACKAGE) $(TARGET) $(DATE) $(GIT_STATUS)

lint:
	golangci-lint run ./... --fix

normalize: fix lint vet update

race:
	go run -race ./... $(COMMAND_LINE)

release: commit
	git tag v$(VERSION)
	git push origin v$(VERSION)
	gh release create v$(VERSION) --generate-notes --latest=true

run:
	go run -ldflags="$(LDFLAGS)" -trimpath ./... $(COMMAND_LINE)

status:
	git status

test:
	go test ./...

update:
	go get -u ./...
	go mod tidy

upgrade:
	go mod edit -go=$(GO_VERSION)
	$(MAKE) update

vet:
	go vet -composites=false ./...

#
# possibly destructive actions
#
gitignore:
	curl -o ./.gitignore $(GITIGNORE_URL)
	cat ./.local.gitignore >> ./.gitignore

init_hook:
	@echo "installing hook 'prepare-commit-msg'"
	@echo '#!/bin/sh' > ./.git/hooks/prepare-commit-msg
	@echo '' >> ./.git/hooks/prepare-commit-msg
	@echo 'COMMIT_MSG_FILE=$$1' >> ./.git/hooks/prepare-commit-msg
	@echo 'COMMIT_SOURCE=$$2' >> ./.git/hooks/prepare-commit-msg
	@echo 'SHA1=$$3' >> ./.git/hooks/prepare-commit-msg
	@echo 'OLLAMA_MODEL="mevatron/diffsense:1.5b"' >> ./.git/hooks/prepare-commit-msg
	@echo 'git diff --staged | ollama run "$$OLLAMA_MODEL" | tee -a "$$COMMIT_MSG_FILE"' >> ./.git/hooks/prepare-commit-msg
	chmod -v +x ./.git/hooks/prepare-commit-msg

#
# init local package
# > make init_localpackage localpackage=package_name
#
init_localpackage:
ifneq ($(localpackage),)
	mkdir ./src/$(localpackage)
	echo "package $(localpackage)" > ./src/$(localpackage)/const.go
	echo "package $(localpackage)" > ./src/$(localpackage)/errors.go
	echo "package $(localpackage)" > ./src/$(localpackage)/func.go
	echo "package $(localpackage)" > ./src/$(localpackage)/init.go
	echo "package $(localpackage)" > ./src/$(localpackage)/method.go
	echo "package $(localpackage)" > ./src/$(localpackage)/type.go
	echo "package $(localpackage)" > ./src/$(localpackage)/var.go
endif

clean-clean-clean-clean: clean
	-gh auth logout

#
# new repo init
# > make init-init-init-init
#
init-init-init-init: clean-clean-clean-clean
	-gh auth logout
	gh auth login --with-token < ~/.git_token
	-gh repo delete $(NAME) --yes
	-rm -Rfv ./.git
	git init
	oco hook set
	git config commit.gpgSign false
	gh repo create $(NAME) --private --source=.
	git add .
	git commit -m "Makefile initial commit ($(DATE))"
	git push --set-upstream origin master
	go mod init $(TARGET)
	go get -u ./...
	go mod tidy
	git add .
	git commit -m "Makefile initial update ($(DATE))"
	git push
