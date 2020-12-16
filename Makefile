.PHONY: gen

uname_m = $(shell uname -m)
uname_s = $(shell uname -s)
MODEL_FILE = $(abspath model/models_gen.go)
OUT_FILE = model.go

install-dep:
	if [[ $(uname_s) = Darwin || $(uname_s) = Linux ]]; then \
	  go get -u -v github.com/fatih/gomodifytags && \
	  go get -u -v github.com/Yamashou/gqlgenc; \
	fi

gen:
	@gqlgenc
	@$(MAKE) account-fake-structs

account-fake-structs:
	@gomodifytags -file $(MODEL_FILE) -struct AccountsInsertInput \
		-field Email -add-tags fake -template "{mailgen}" > $(OUT_FILE)
	@gomodifytags -file $(OUT_FILE) -struct AccountsInsertInput \
		-field LastLogin -add-tags fake -template "skip" > $(MODEL_FILE)
	@gomodifytags -file $(MODEL_FILE) -struct AccountsInsertInput \
		-field Role -add-tags fake -template "{rolegen}" > $(OUT_FILE)
	@gomodifytags -file $(OUT_FILE) -struct AccountsInsertInput \
		-field User -add-tags fake -template "skip" > $(MODEL_FILE)
	@gomodifytags -file $(MODEL_FILE) -struct AccountsInsertInput \
		-field UserID -add-tags fake -template "{uuid}" > $(OUT_FILE)
	@cp $(OUT_FILE) $(MODEL_FILE)
	@rm -rf $(OUT_FILE)

