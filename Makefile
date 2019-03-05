test :
	$(MAKE) -C lambda/send test

build :
	$(MAKE) -C lambda/send build

terratest : build
	$(MAKE) -C deploy/test/send_lambda

deploy-% : build
	$(MAKE) -C deploy/$*/send_lambda