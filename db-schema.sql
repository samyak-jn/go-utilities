CREATE TABLE ClusterCredentials (
	
	-- Primary key for the credentials (UID)
	clustercredentials_cred_id VARCHAR (48) UNIQUE PRIMARY KEY,

	-- API URL for the cluster 
	-- Example: https://api.ci-ln-dlfw0qk-f76d1.origin-ci-int-gce.dev.openshift.com:6443
	host VARCHAR (512),

	-- State 1) kube_config containing a token to a service account that has the permissions we need.
	kube_config VARCHAR (65000),

	-- State 1) The name of a context within the kube_config 
	kube_config_context VARCHAR (64),

	-- State 2) ServiceAccount bearer token from the target manager cluster
	serviceaccount_bearer_token VARCHAR (128),

	-- State 2) The namespace of the ServiceAccount
	serviceaccount_ns VARCHAR (128),

	seq_id serial
);