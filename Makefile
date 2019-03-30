# Tool For preparing resources

# Path
ASSETPATH=./asset
ASSETFILE=data.go
ASSETPACKAGE=asset

# Go commands
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean

GOASSETBUILDER=go-assets-builder

# For compatibility
ifeq ($(OS),Windows_NT)
	RM = cmd.exe /C del /Q
else
	RM = rm -f
endif

.PHONY: all assets

all: assets

assets:
	@echo "Ready assets..."
	@cd $(ASSETPATH); \
		$(GOASSETBUILDER) --package=$(ASSETPACKAGE) ./static/ > $(ASSETFILE)

clean:
	$(GOCLEAN)
	@$(RM) $(ASSETPATH)/$(ASSETFILE)

deps:
	dep ensure
	@echo "Installing go-assets-builder..."
	@cd ./vendor/github.com/jessevdk/go-assets-builder && $(GOINSTALL) .