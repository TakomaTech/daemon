GOCMD=go
BINARY=daemon

.PHONY: all build plugin cpp run clean

all: build

build:
	$(GOCMD) build -o $(BINARY) ./cmd/daemon

plugin:
	$(GOCMD) build -buildmode=plugin -o plugins/example.so ./plugins/example

cpp:
	g++ -std=c++17 -O2 -c internal/cpp/engine.cpp -o internal/cpp/engine.o
	ar rcs internal/cpp/libengine.a internal/cpp/engine.o

install: build plugin
	install -d $(DESTDIR)/usr/local/bin
	install -m 755 $(BINARY) $(DESTDIR)/usr/local/bin/$(BINARY)
	install -d $(DESTDIR)/usr/local/lib/daemon/plugins
	install -m 644 plugins/example.so $(DESTDIR)/usr/local/lib/daemon/plugins/

run: build
	./$(BINARY)

run-headless: build
	DAEMON_HEADLESS=1 ./$(BINARY)

clean:
	rm -f $(BINARY) plugins/example.so internal/cpp/engine.o internal/cpp/libengine.a
