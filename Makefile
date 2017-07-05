################################################################
# Makefile for register
# 2015-12-30 M.Horigome
#
# 注意
# Windowsでは、MSYS2とGNU Makeが別途必要です
################################################################

VERSION=1.0.0.1
TARGET=flist

ifeq ($(OS),Windows_NT)
	# For Windows
	THIS_TARGET=$(TARGET).exe
else
	# For Linux, MacOSX
	THIS_TARGET=$(TARGET)
endif

# Packages Lists
# ビルドターゲットが利用するVendorパッケージを指定します
VENDERS= golang.org/x/text/encoding/japanese \
		 golang.org/x/text/transform

################################################################
GOPATH=$(CURDIR)/_vendor
GO=GOPATH=$(GOPATH) go

#GO_COMMON_OPTS=
GO_BUILD=$(GO) build -ldflags "-X main.version=$(VERSION)" $(GO_COMMON_OPTS)
GO_TEST=$(GO) test -v $(GO_COMMON_OPTS)
GO_GET=$(GO) get -d -v -u $(GO_COMMON_OPTS)

# _shared の設定
SRC = $(wildcard $(CURDIR)/*.go)
OUT = ./_build

# Cross compile
OS_WIN =GOOS=windows
OS_LINUX =GOOS=linux
OS_MAC =GOOS=darwin

AR_AMD64 =GOARCH=amd64
AR_386 =GOARCH=386

################################################################

# このアーキテクチャでのビルド
this:
	@echo "==> make this target"
	@$(GO_BUILD) -o $(OUT)/$(THIS_TARGET) $(SRC)

# ビルド、テスト
all:win_x64 win_x86 linux_x86 linux_x64 mac_x64

vendor_clean:
	rm -rf ./_vendor/src

vendor_get: vendor_clean
	$(GO_GET) $(VENDERS)

# 各パッケージの内部の.gitを削除すること
# * Windows では find コマンドが動作しないかもしれない。
vendor_update: vendor_get
	rm -rf `find ./_vendor/src -type d -name .git` \
	&& rm -rf `find ./_vendor/src -type d -name .hg` \
	&& rm -rf `find ./_vendor/src -type d -name .bzr` \
	&& rm -rf `find ./_vendor/src -type d -name .svn`

# main のテスト
test_main:
	@$(GO_TEST)

# 全てのテスト
test:test_main

clean:
	@echo "==> clean output"
	rm -rf $(OUT)


# 各アーキテクチャごとのビルド
win_x64:
	@echo "==> make Windows x64 target"
	$(OS_WIN) $(AR_AMD64) $(GO_BUILD) -o $(OUT)/win_x64/$(TARGET).exe $(SRC)

win_x86:
	@echo "==> make Windows x86 target"
	$(OS_WIN) $(AR_386) $(GO_BUILD) -o $(OUT)/win_x86/$(TARGET).exe $(SRC)

linux_x64:
	@echo "==> make linux x64 target"
	$(OS_LINUX) $(AR_AMD64) $(GO_BUILD) -o $(OUT)/linux_x64/$(TARGET) $(SRC)

linux_x86:
	@echo "==> make linux x86 target"
	$(OS_LINUX) $(AR_386) $(GO_BUILD) -o $(OUT)/linux_x86/$(TARGET) $(SRC)

mac_x64:
	@echo "==> make MacOS X target"
	$(OS_MAC) $(AR_AMD64) $(GO_BUILD) -o $(OUT)/mac_x64/$(TARGET) $(SRC)
