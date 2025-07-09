# Cargo commands
CARGO := cargo

# Logger
define log_info
	echo -e "[\033[0;33m*\033[0m] $(1)"
endef

define log_success
	echo -e "[\033[0;32m+\033[0m] Done"
endef

release:
	@ $(call log_info,[release] Compiling...)
	@ $(CARGO) build --release
	@ $(call log_success)

debug:
	@ $(call log_info,[debug] Compiling...)
	@ $(CARGO) build --release --features debug
	@ $(call log_success)

clean:
	@ $(call log_info,Cleaning build artifacts)
	@ rm -rf target
	@ $(call log_success)

.PHONY: release debug clean
