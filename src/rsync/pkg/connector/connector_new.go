// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package connector

import (
	"context"
	"fmt"

	"strings"

	pkgerrors "github.com/pkg/errors"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/internal/utils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/plugins/anthos"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/plugins/azurearcv2"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/plugins/fluxv2"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/plugins/k8s"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/plugins/k8sexp"
	. "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/types"
)

// Connection is for a cluster
type Provider struct {
	cid string
}

func NewProvider(id interface{}) Provider {
	return Provider{
		cid: fmt.Sprintf("%v", id),
	}
}

func (p *Provider) GetClientProviders(ctx context.Context, app, cluster, level, namespace string) (ClientProvider, error) {
	// Default Provider type
	var providerType string = "k8s"

	result := strings.SplitN(cluster, "+", 2)
	if len(result) != 2 {
		log.Error("Invalid cluster name format::", log.Fields{"cluster": cluster})
		return nil, pkgerrors.New("Invalid cluster name format")
	}
	kc, err := utils.GetKubeConfig(ctx, cluster, level, namespace)
	if err != nil {
		if !strings.Contains(err.Error(), "Invalid kubeconfig") {
			return nil, err
		}
	}

	if len(kc) > 0 {
		providerType = "k8s"
	} else {
		// GitOps uses level "0" credentials at this time
		c, err := utils.GetGitOpsConfig(ctx, cluster, "0", "default")
		if err != nil {
			return nil, err
		}
		providerType = c.Props.GitOpsType
		if providerType == "" {
			return nil, pkgerrors.New("No provider type specified")
		}
	}

	switch providerType {
	case "k8s":
		cl, err := k8s.NewK8sProvider(ctx, p.cid, app, cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		return cl, nil
		// This case is unused at this time.
		// In the above case K8s plugin each resource is written
		// to the cluster individually. The disadvantage is
		// that if any resource fails then the application is
		// left in bad state on the cluster with some resources
		// already applied on the cluster. In this plugin all
		// application resources are collected in a temporary
		// file, and then applied together to the cluster.
		// All or no resources will be applied to the cluster.
		// More Disk space is required in this approach.
	case "k8sExp":
		cl, err := k8sexp.NewK8sProvider(ctx, p.cid, app, cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		return cl, nil

	case "fluxcd":
		cl, err := fluxv2.NewFluxv2Provider(ctx, p.cid, app, cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		return cl, nil
	case "azureArcV2":
		cl, err := azurearcv2.NewAzureArcProvider(ctx, p.cid, app, cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		return cl, nil
	case "anthos":
		cl, err := anthos.NewAnthosProvider(ctx, p.cid, app, cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		return cl, nil
	}
	return nil, pkgerrors.New("Provider type not supported")
}
