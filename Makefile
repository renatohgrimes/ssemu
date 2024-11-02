release:
	sh scripts/make_env.sh
	sh scripts/make_images.sh
	sh scripts/make_binaries_linux.sh
	sh scripts/make_binaries_windows.sh
	sh scripts/make_release.sh

tests:
	sh scripts/make_env.sh
	sh scripts/make_images.sh
	sh scripts/make_binaries_linux.sh
	sh scripts/run_tests.sh

compose:
	sh scripts/make_env.sh
	sh scripts/make_images.sh
	sh scripts/make_binaries_linux.sh
	sh scripts/run_compose.sh

devenv:
	sh scripts/make_env_dev.sh

lint:
	sh scripts/run_lint.sh