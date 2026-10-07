COMPLETION_DIR="$(HOME)/.bash_completion.d"

build:
	mkdir -p bin
	go build -o ./bin ./cmd/worktree

install: build
	mkdir -p ~/bin
	cp ./bin/worktree ~/bin/worktree

uninstall:
	rm ~/bin/worktree

completion: build
	mkdir -p $(COMPLETION_DIR)
	./bin/worktree completion bash > "$(COMPLETION_DIR)/worktree"
	echo "completion saved to $(COMPLETION_DIR)/worktree"
