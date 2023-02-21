all: makedir build/openapi.yaml
	redoc-cli bundle build/openapi.yaml -t ./html/index.hbs -o ../public/index.html \
	    --options.theme.typography.links.color=#695de4 \

sp :=

sp += # add space

my-sort = $(shell echo $(subst $(sp),'\n',$2) | sort $1 --key=1,1 -)

define appendFilesWithTabs
	for file in $(call my-sort,,$(1)); do \
		while IFS= read -r line || [ -n "$$line" ]; do \
			echo "$(2)$$line"; \
		done < $$file; \
	done >> $3
endef

build/models.yaml: ./models/*.yaml
	echo "components:" > $@
	echo "  schemas:" >> $@
	$(call appendFilesWithTabs, $^,""    "", $@)

build/paths.yaml: ./paths/*.yaml
	echo "paths:" > $@
	$(call appendFilesWithTabs, $^,""  "", $@)

build/openapi.yaml: index.yaml build/paths.yaml build/models.yaml
	echo "#This file has been generated automatically. DO NOT MODIFY!" > $@
	for file in $^; do \
		while IFS= read -r line; do \
			echo "$$line"; \
		done < $$file; \
	done >> $@

makedir:
	mkdir -p build;

clean:
	rm -rf build/*
