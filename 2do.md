- [x] Add logging
- [x] Dockerize, Makefile, Basic CI/CD
- [x] Telegram is at the presentation layer, it doesn't need to be abstracted behind an interface.
- [x] What is tools/get_token.go? there is no need for that. (UPDATE: Removed)
- [x] Use Cobra as explained in Readme.
- [x] ChatID whitelisting (has a TODO, middleware whitelisting) for security concerns
- [x] Define semi-hardcoded buttons! (button unique identifiers all totally hard coded in two places which lead to confusion, one when defining and one when putting in the handler, it's better to define button somehow else)
- [x] Remove session which is used as a one size fits all!
- [x] Debug container logs, also what to do in case of high amount of logs?
- [x] Zero stats problem
- [x] Container Start Stop handlers and functionality
- [x] Welcome message should have button and better message
- [x] Logs streaming
- [x] Add image handlers, next, prev, back
    - [x] Image size should be in human readable units
    - [x] Image ID should be shorted (trimmed)
    - [ ] Image name and tag should be separated???
    - [x] Add Image status and created at
    - [ ] Image run, rename, remove commands 
    - [x] Error on line 78 if log file about ">" character should be resolved (every value is escaped now)
- [ ] Add ability to run containers out of images
- [ ] Ability to filter Containers and Images
- [x] Gaining at least 50% test coverage (~80% with the Docker integration tests)
- [x] Debug container stats as it sounds to have problem in some cases (precpu_stats was read from the wrong field)
- [x] Remove previous messages buttons (when they are not required) (one panel message, edited in place)
- [x] Handle callback queries to show messages, e.g. for remove.
    - Helper: &telebot.CallbackResponse{Text: fmt.Sprint(ctx.Data(), "removed!")}
- [x] Incompatibility between github.com/docker/docker and github.com/moby/moby (only docker/docker v28 now)
- [ ] If the image is in use, include a button to show the containers using it, back button should get back to the list of Images
- [ ] Containers list using filters filters.Args should use a better alternative as All is not a filter and should be an argument for containers list

- [ ] Add Slack
- [ ] Add discord

- [ ] Add github actions security check with gosec (should not exit with status 1 in case of non-critical issues)
    ```yaml
    - name: Gosec Security Scanner
      run: |
        export PATH=$PATH:$(go env GOPATH)/bin
        go install github.com/securego/gosec/v2/cmd/gosec@latest
        echo "[gosec]\n  severity = \"medium\"\n" > .gosec.toml
        gosec ./...
    ```

- [ ] Load yaml config
- [ ] Implement web hook and let the user decide for it.
- [ ] Make it ready to be installed using go install, apt install, etc.
- [ ] Should we replace docker by contract and use it directly as we will never replace it? Not actually, as we may want to plug a mock instead or wrap its functions
- [ ] `dockeroller/docker/docker.go:29` We should get Docker ENVs from user in case of not default
- [ ] Add linter (golangci-linter to actions)
- [x] Pass ID to session like the way dongi works, or maybe better! (tbb sessions)

