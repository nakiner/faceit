# Info

Directory contains all manifests to start service in production cloud

Traffic termination done by istio itself, so first you need to install istio-operator and configure istio-ingressgateway on target cluster. 

Project DOES NOT include deployments for target database, so starting container postgresql:14 in desired namespace is up to you.

If you do not need to terminate TLS traffic simply remove https block from gateway. 