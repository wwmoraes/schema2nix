.PHONY: hooks/%
hooks/pre-bump:: check
	$(if ${VERSION},nix-update schema2nix --flake --version ${VERSION})

hooks/post-bump::
	git push
	$(if ${VERSION},git push $(shell git remote) ${VERSION})
