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
    - [ ] Add Image status and created at
    - [ ] Image run, rename, remove commands 
    - [ ] Error on line 78 if log file about ">" character should be resolved
- [ ] Add ability to run containers out of images
- [ ] Ability to filter Containers and Images
- [ ] Gaining at least 50% test coverage
- [ ] Debug container stats as it sounds to have problem in some cases
- [ ] Remove previous messages buttons (when they are not required)
- [ ] Handle callback queries to show messages, e.g. for remove.
    - Helper: &telebot.CallbackResponse{Text: fmt.Sprint(ctx.Data(), "removed!")}
- [ ] Incompatibility between github.com/docker/docker and github.com/moby/moby
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
