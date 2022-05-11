package api

type ClusterRegistration struct {
	// pull or push
	SyncMode          string
	KarmadaKubeConfig string
	KarmadaContext    string
	ClusterKubeConfig string
	ClusterContext    string
	MemberClusterName string
}
