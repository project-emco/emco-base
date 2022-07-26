// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package fluxv2

import (
	"context"

	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/internal/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// StartClusterWatcher watches for CR changes in git location
func (p *Fluxv2Provider) StartClusterWatcher(ctx context.Context) error {
	p.gitProvider.StartClusterWatcher(ctx)
	return nil
}

// ApplyStatusCR applies status CR
func (p *Fluxv2Provider) ApplyStatusCR(ctx context.Context, name string, content []byte) error {

	// Add namespace to the status resource, needed by
	// Flux
	//Decode the yaml to create a runtime.Object
	unstruct := &unstructured.Unstructured{}
	//Ignore the returned obj as we expect the data in unstruct
	_, err := utils.DecodeYAMLData(string(content), unstruct)
	if err != nil {
		return err
	}
	// Set Namespace
	unstruct.SetNamespace(p.gitProvider.Namespace)
	b, err := unstruct.MarshalJSON()
	if err != nil {
		return err
	}
	path := p.gitProvider.GetPath("context") + name + ".yaml"
	ref, err := p.gitProvider.Apply(path, nil, b)
	if err != nil {
		return err
	}
	p.gitProvider.Commit(ctx, ref)
	return err

}

// DeleteStatusCR deletes status CR
func (p *Fluxv2Provider) DeleteStatusCR(ctx context.Context, name string, content []byte) error {
	path := p.gitProvider.GetPath("context") + name + ".yaml"
	ref, err := p.gitProvider.Delete(path, nil, content)
	if err != nil {
		return err
	}
	p.gitProvider.Commit(ctx, ref)
	return err
}
