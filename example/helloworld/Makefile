PLUGIN_DIR := com.github.morikuni.goStreamDeckSDK.sdPlugin

generate-manifest:
	go run gen-manifest/main.go > $(PLUGIN_DIR)/manifest.json
	
build:
	go build -o $(PLUGIN_DIR)/helloworld main.go
	
install-mac: generate-manifest build
	rm -r ~/Library/Application\ Support/com.elgato.StreamDeck/Plugins/$(PLUGIN_DIR) || true
	cp -R $(PLUGIN_DIR) ~/Library/Application\ Support/com.elgato.StreamDeck/Plugins/$(PLUGIN_DIR)