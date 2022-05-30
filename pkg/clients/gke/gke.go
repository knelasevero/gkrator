package gke

import (
	"context"
	"fmt"

	"cloud.google.com/go/container/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	v1alpha1 "github.com/knelasevero/gkrator/api/v1alpha1"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	gkratorclient "github.com/knelasevero/gkrator/pkg/clients"
)

const (
	CloudPlatformRole                         = "https://www.googleapis.com/auth/cloud-platform"
	errSpec = "cluster spec is nil"
	errFetchSAKSecret                         = "could not fetch SecretAccessKey secret: %w"
	errMissingSAK                             = "missing SecretAccessKey"
	errUnableProcessJSONCredentials           = "failed to process the provided JSON credentials: %w"
	errUnableCreateGCPClient                = "failed to create GCP client: %w"
	errUnableGetCredentials                   = "unable to get credentials: %w"
)

type GoogleGKEClient interface {
    CreateCluster(ctx context.Context, req *containerpb.CreateClusterRequest, opts ...gax.CallOption) (*containerpb.Operation, error) 
    CreateNodePool(ctx context.Context, req *containerpb.CreateNodePoolRequest, opts ...gax.CallOption) (*containerpb.Operation, error)
    DeleteCluster(ctx context.Context, req *containerpb.DeleteClusterRequest, opts ...gax.CallOption) (*containerpb.Operation, error)
    DeleteNodePool(ctx context.Context, req *containerpb.DeleteNodePoolRequest, opts ...gax.CallOption) (*containerpb.Operation, error)
	Close() error
}

type GKE struct {
	credentials []byte
	GoogleGKEClient GoogleGKEClient
	projectID string
}


// NewClient constructs a GCP Provider.
func (g *GKE) NewClient(ctx context.Context, gke v1alpha1.GoogleKubernetesEngine, kube kclient.Client, namespace string) (gkratorclient.Client, error) {

	defer func() {
		// closes IAMClient to prevent gRPC connection leak in case of an error.
		if g.GoogleGKEClient == nil {
			_ = g.GoogleGKEClient.Close()
		}
	}()

	g.projectID = gke.Spec.ProjectID

	ts, err := serviceAccountTokenSource(ctx, gke, kube, namespace)
	if err != nil {
		return nil, fmt.Errorf(errUnableCreateGCPClient, err)
	}

	// check if we can get credentials
	_, err = ts.Token()
	if err != nil {
		return nil, fmt.Errorf(errUnableGetCredentials, err)
	}

	c, err := container.NewClusterManagerClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}
	g.GoogleGKEClient = c
	return c, nil
}

func serviceAccountTokenSource(ctx context.Context, gke v1alpha1.GoogleKubernetesEngine, kube kclient.Client, namespace string) (oauth2.TokenSource, error) {
	sr := gke.Spec.Auth.SecretRef
	if sr == nil {
		return nil, nil
	}
	credentialsSecret := &v1.Secret{}
	credentialsSecretName := sr.SecretAccessKey.Name
	objectKey := types.NamespacedName{
		Name:      credentialsSecretName,
		Namespace: namespace,
	}

	err := kube.Get(ctx, objectKey, credentialsSecret)
	if err != nil {
		return nil, fmt.Errorf(errFetchSAKSecret, err)
	}
	credentials := credentialsSecret.Data[sr.SecretAccessKey.Key]
	if (credentials == nil) || (len(credentials) == 0) {
		return nil, fmt.Errorf(errMissingSAK)
	}
	config, err := google.JWTConfigFromJSON(credentials, CloudPlatformRole)
	if err != nil {
		return nil, fmt.Errorf(errUnableProcessJSONCredentials, err)
	}
	return config.TokenSource(ctx), nil
}